package sdk_demo

import (
	"os"
	"testing"
)

func Test_wangsuDemo(t *testing.T) {
	if os.Getenv("GO_CDN_CHECK_WANGSU_ACCESS_KEY") == "" || os.Getenv("GO_CDN_CHECK_WANGSU_SECRET_KEY") == "" {
		t.Skip("set Wangsu CDN credentials to run this integration test")
	}
	wangsuDemo()
}
