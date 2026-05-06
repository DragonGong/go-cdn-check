package sdk_demo

import (
	"os"
	"testing"
)

func Test_baiduDemo(t *testing.T) {
	if os.Getenv("GO_CDN_CHECK_BAIDU_ACCESS_KEY") == "" || os.Getenv("GO_CDN_CHECK_BAIDU_SECRET_KEY") == "" {
		t.Skip("set Baidu CDN credentials to run this integration test")
	}
	baiduDemo()
}
