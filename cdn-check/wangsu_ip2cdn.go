package cdn_check

import (
	"encoding/json"
	"errors"
	"github.com/dragongong/go-cdn-check/cdn-check/config"
	"github.com/dragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/api/client"
	"github.com/dragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/common/auth"
)

type WangsuChecker struct {
	config            auth.AkskConfig
	cdnCheckerManager *CdnCheckerManager
}

// CheckListItem 定义了checkList中的单个元素
type CheckListItem struct {
	Response string `json:"response"`
	IP       string `json:"ip"`
}

// Result 定义了包含checkList的result结构体
type Result struct {
	CheckList []CheckListItem `json:"checkList"`
}

// Response 定义了整个响应的结构体，包括Result和code
type WangsuResponse struct {
	Result Result `json:"result"`
	Code   int    `json:"code"`
	Error  string `json:"error,omitempty"` // 可以添加一个可选的error字段来处理可能的错误信息
}

func NewWangsuChecker() Checker {
	if len(config.CdnConfig.Wangsu.AccessKey) == 0 || len(config.CdnConfig.Wangsu.SecretKey) == 0 {
		return nil
	}
	return &WangsuChecker{
		config: struct {
			AccessKey     string
			SecretKey     string
			Uri           string
			Method        string
			EndPoint      string
			SignedHeaders string
		}{
			AccessKey: config.CdnConfig.Wangsu.AccessKey,
			SecretKey: config.CdnConfig.Wangsu.SecretKey,
			Uri:       "/api/si/tools/ipCheck",
			Method:    "POST",
			EndPoint:  "open.chinanetcenter.com",
		},
		cdnCheckerManager: NewCdnCheckerManager("wangsu", wangsuQps, wangsuMaxIpsNums),
	}
}
func (w *WangsuChecker) CheckConfig() bool {
	return len(config.CdnConfig.Wangsu.AccessKey) > 0 && len(config.CdnConfig.Wangsu.SecretKey) > 0
}
func (w *WangsuChecker) GetType() string {
	return w.cdnCheckerManager.checkerName
}
func (w *WangsuChecker) check(ips []string) (map[string]bool, error) {
	if len(ips) > w.cdnCheckerManager.MaxIpsNum {
		return nil, errors.New("the number of IPs exceeds the limit")
	}
	querySpecificIPBelongRequest := client.QuerySpecificIPBelongRequest{}
	//= "116.114.98.35"
	var subquerySpecificIPBelongRequestList []*string
	for _, ip := range ips {
		subquerySpecificIPBelongRequestList = append(subquerySpecificIPBelongRequestList, &ip)
	}
	querySpecificIPBelongRequest.SetIp(subquerySpecificIPBelongRequestList)
	response := auth.Invoke(w.config, querySpecificIPBelongRequest.String())
	responseStruct := &WangsuResponse{}
	err := json.Unmarshal([]byte(response), responseStruct)
	if err != nil {
		return nil, err
	}
	if responseStruct.Code != 200 {
		return nil, errors.New("code isn't 200")
	}

	ipCdnMap := make(map[string]bool)
	// 得保证map的key是有原始数据过来的,传接口之后出现丢失的ip为false
	for _, ip := range ips {
		ipCdnMap[ip] = false
	}
	for _, ipsInfo := range responseStruct.Result.CheckList {
		if _, ok := ipCdnMap[ipsInfo.IP]; ok {
			if ipsInfo.Response == "yes" {
				ipCdnMap[ipsInfo.IP] = true
			}
		}
	}
	return ipCdnMap, nil
}

func (w *WangsuChecker) CdnFilter(ips []string) []string {
	return w.cdnCheckerManager.CdnFilter(ips, w.check)
}
