package sdk_demo

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
)

// qps 无限制   100ip
func baiduDemo() {
	client := GetDefaultClient()
	ipsInfo, _ := client.GetIpListInfo([]string{"116.114.98.35", "59.24.3.174"}, "describeIp")
	fmt.Println(ipsInfo)
}

func GetDefaultClient() *cdn.Client {
	ak := "ALTAKJ0b68ZR2DYTcomAxjo3kp"
	sk := "6a06b1b9ccd04880a35c44a769608c5f"
	endpoint := "https://cdn.baidubce.com"
	client, _ := cdn.NewClient(ak, sk, endpoint)
	return client
}
