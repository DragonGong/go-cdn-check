package sdk_demo

import (
	"encoding/json"
	"fmt"
	cdnCheck "github.com/dragongong/go-cdn-check/cdn-check"
	"github.com/dragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/api/client"
	"github.com/dragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/common/auth"
	"os"
)

func wangsuDemo() {

	querySpecificIPBelongRequest := client.QuerySpecificIPBelongRequest{}
	var subquerySpecificIPBelongRequest0 = "116.114.98.35"
	querySpecificIPBelongRequest.SetIp([]*string{&subquerySpecificIPBelongRequest0})
	var config auth.AkskConfig
	config.AccessKey = os.Getenv("GO_CDN_CHECK_WANGSU_ACCESS_KEY")
	config.SecretKey = os.Getenv("GO_CDN_CHECK_WANGSU_SECRET_KEY")
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
