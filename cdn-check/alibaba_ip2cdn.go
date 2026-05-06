package cdn_check

import (
	"encoding/json"
	"errors"
	"fmt"
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/dragongong/go-cdn-check/cdn-check/config"
	"github.com/dragongong/go-cdn-check/logger"
	"strings"
	"sync"
	"time"
)

// 单用户限制50次/秒
// Description:
// 使用AK&SK初始化账号Client
// 必须开启cdn的权限

type AlibabaChecker struct {
	runtime           *util.RuntimeOptions
	client            *cdn20180510.Client
	cdnCheckerManager *CdnCheckerManager
	maxConcurrency    int // 保险起见的最大阈值，如果ip数量大于它，就报错
}

type AlibabaResponse struct {
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
	Body       AlibabaBody       `json:"body"`
}

type AlibabaBody struct {
	CdnIp       string `json:"CdnIp"`
	ISP         string `json:"ISP"`
	IspEname    string `json:"IspEname"`
	Region      string `json:"Region"`
	RegionEname string `json:"RegionEname"`
	RequestId   string `json:"RequestId"`
}

func NewAlibabaChecker() Checker {
	if len(config.CdnConfig.Alibaba.AccessKey) == 0 || len(config.CdnConfig.Alibaba.SecretKey) == 0 {
		return nil
	}
	client, err := CreateAlibabaCdnClient()
	if err != nil {
		logger.Log.Error("NewAlibabaChecker error ")
		return nil
	}
	return &AlibabaChecker{
		maxConcurrency:    20,
		runtime:           &util.RuntimeOptions{},
		client:            client,
		cdnCheckerManager: NewCdnCheckerManager("alibaba", alibabaQps, alibabaMaxIpsNums),
	}
}

func CreateAlibabaCdnClient() (_result *cdn20180510.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(config.CdnConfig.Alibaba.AccessKey),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(config.CdnConfig.Alibaba.SecretKey),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Cdn
	config.Endpoint = tea.String("cdn.aliyuncs.com")
	_result = &cdn20180510.Client{}
	_result, _err = cdn20180510.NewClient(config)
	return _result, _err
}

func (a *AlibabaChecker) CheckConfig() bool {
	return len(config.CdnConfig.Alibaba.AccessKey) > 0 && len(config.CdnConfig.Alibaba.SecretKey) > 0
}
func (a *AlibabaChecker) GetType() string {
	return a.cdnCheckerManager.checkerName
}
func (a *AlibabaChecker) doCheck(ip string) (bool, error) {
	describeIpInfoRequest := &cdn20180510.DescribeIpInfoRequest{
		IP: tea.String(ip),
	}
	res, tryErr := func() (res string, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_res, err := a.client.DescribeIpInfoWithOptions(describeIpInfoRequest, a.runtime)
		if err != nil {
			return "", err
		}
		resByte, _ := json.Marshal(*_res)
		res = string(resByte)
		return
	}()
	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, err := util.AssertAsString(error.Message)
		if err != nil {
			return false, err
		}
	}
	response := &AlibabaResponse{}
	err := json.Unmarshal([]byte(res), response)
	if err != nil {
		return false, err
	}
	if response.StatusCode != 200 {
		return false, fmt.Errorf("alibaba response.StatusCode is %d", response.StatusCode)
	}

	if response.Body.CdnIp == "True" {
		return true, nil
	} else {
		return false, nil
	}
}
func (a *AlibabaChecker) checkItem(ips []string) (map[string]bool, error) {
	if len(ips) > 1 {
		return nil, errors.New("the number of IPs exceeds the limit")
	}
	m := make(map[string]bool)
	m[ips[0]], _ = a.doCheck(ips[0])
	return m, nil
}
func (a *AlibabaChecker) check(ips []string) (map[string]bool, error) {
	if len(ips) > a.cdnCheckerManager.MaxIpsNum {
		return nil, errors.New("the number of IPs exceeds the limit")
	} else if len(ips) == 0 {
		return nil, errors.New("the number of Ips is zero")
	}
	ipCdnMap := make(map[string]bool)
	for _, ip := range ips {
		isCdn, err := a.doCheck(ip)
		if err != nil {
			return nil, err
		}
		ipCdnMap[ip] = isCdn
	}
	return ipCdnMap, nil
}
func (a *AlibabaChecker) checkConcurrency(ips []string) (map[string]bool, error) {
	// todo：阿里巴巴接口并发还没开发完，由于它只能一次传一个参数，导致速度很慢，现在再尝试做并发调接口，还未完成

	if len(ips) > a.maxConcurrency {
		return nil, errors.New("the number of IPs exceeds the limit")
	}
	resultC := make(chan string, a.maxConcurrency)
	count := 0
	ticker := time.NewTicker(time.Millisecond * 1)
	defer ticker.Stop()

	wg := &sync.WaitGroup{}
	for range ticker.C {
		if count >= len(ips) {
			break
		}
		wg.Add(1)
		go func(countG int) {
			checkRes, _ := a.doCheck(ips[countG])
			if checkRes {
				// 如果是cdn
				resultC <- ips[countG]
			}
			wg.Done()
		}(count)
		count++
	}
	go func() {
		wg.Wait()
		close(resultC)
	}()
	ipCdnMap := make(map[string]bool)
	for _, ip := range ips {
		ipCdnMap[ip] = false
	}
	for ip := range resultC {
		if _, ok := ipCdnMap[ip]; ok {
			ipCdnMap[ip] = true
		}
	}
	return ipCdnMap, nil
}

func (a *AlibabaChecker) CdnFilter(ips []string) []string {
	return a.cdnCheckerManager.CdnFilter(ips, a.check)
}
