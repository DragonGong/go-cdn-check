package sdk_demo

import (
	"encoding/json"
	"fmt"
	cdnCheck "github.com/clouddragongong/go-cdn-check/cdn-check"
	"github.com/clouddragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/api/client"
	"github.com/clouddragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/common/auth"
)

func wangsuDemo() {

	querySpecificIPBelongRequest := client.QuerySpecificIPBelongRequest{}
	var subquerySpecificIPBelongRequest0 = "116.114.98.35"
	querySpecificIPBelongRequest.SetIp([]*string{&subquerySpecificIPBelongRequest0})
	var config auth.AkskConfig
	config.AccessKey = "7a7VWPGTL2cWgzJ5SxQqhvKztxQQlHOk4I55"
	config.SecretKey = "x9rG3E3XqP3WgPoCjosB07qxi3UV4eOGgTI5dtYzUaLu9TR9STsYsgkwj9hHVU4J"
	config.EndPoint = "open.chinanetcenter.com"
	config.Uri = "/api/si/tools/ipCheck"
	config.Method = "POST"
	response := auth.Invoke(config, querySpecificIPBelongRequest.String())
	responseMap := &cdnCheck.WangsuResponse{}
	err := json.Unmarshal([]byte(response), responseMap)
	if err != nil {
		return
	}
	fmt.Printf("response body is %#v\n", response)
}
