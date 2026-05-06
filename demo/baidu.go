package demo

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"os"
)

func baiduDemo() {
	client := GetDefaultClient()
	ipsInfo, _ := client.GetIpListInfo([]string{"116.114.98.35", "59.24.3.174"}, "describeIp")
	fmt.Println(ipsInfo)
}

func GetDefaultClient() *cdn.Client {
	ak := os.Getenv("GO_CDN_CHECK_BAIDU_ACCESS_KEY")
	sk := os.Getenv("GO_CDN_CHECK_BAIDU_SECRET_KEY")
	endpoint := "https://cdn.baidubce.com"
	client, _ := cdn.NewClient(ak, sk, endpoint)
	return client
}
