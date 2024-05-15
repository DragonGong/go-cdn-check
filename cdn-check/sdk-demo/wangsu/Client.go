package wangsu

import (
	"fmt"
	"github.com/clouddragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/api/client"
	"github.com/clouddragongong/go-cdn-check/cdn-check/sdk-demo/wangsu/common/auth"
)

func main() {

	querySpecificIPBelongRequest := client.QuerySpecificIPBelongRequest{}

	var config auth.AkskConfig
	config.AccessKey = "{accessKey}"
	config.SecretKey = "{secretKey}"
	config.EndPoint = "open.chinanetcenter.com"
	config.Uri = "/api/si/tools/ipCheck"
	config.Method = "POST"
	response := auth.Invoke(config, querySpecificIPBelongRequest.String())
	fmt.Printf("response body is %#v\n", response)
}
