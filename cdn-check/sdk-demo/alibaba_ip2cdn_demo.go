package sdk_demo

import (
	"encoding/json"
	"fmt"
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"strings"
)

// 单用户限制50次/秒
// Description:
//
// 使用AK&SK初始化账号Client
//
// @return Client
//
// @throws Exception
// 必须开启cdn的权限
func CreateClient() (_result *cdn20180510.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(os.Getenv("GO_CDN_CHECK_ALIBABA_ACCESS_KEY")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(os.Getenv("GO_CDN_CHECK_ALIBABA_SECRET_KEY")),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Cdn
	config.Endpoint = tea.String("cdn.aliyuncs.com")
	_result = &cdn20180510.Client{}
	_result, _err = cdn20180510.NewClient(config)
	return _result, _err
}

func alibabaMain() (_err error) {
	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	describeIpInfoRequest := &cdn20180510.DescribeIpInfoRequest{
		IP: tea.String("173.245.48.12"),
	}
	runtime := &util.RuntimeOptions{
		ReadTimeout:    tea.Int(1000),
		ConnectTimeout: tea.Int(1000),
		IgnoreSSL:      tea.Bool(false),
		Autoretry:      tea.Bool(true),
		MaxAttempts:    tea.Int(1),
	}
	res, tryErr := func() (res string, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_res, err := client.DescribeIpInfoWithOptions(describeIpInfoRequest, runtime)
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
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	fmt.Println(res)
	response := make(map[string]interface{})
	err := json.Unmarshal([]byte(res), &response)
	if err != nil {
		return err
	}
	return _err
}
