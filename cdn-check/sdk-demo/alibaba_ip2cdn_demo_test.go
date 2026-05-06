package sdk_demo

import (
	"os"
	"testing"
)

func Test__main(t *testing.T) {
	if os.Getenv("GO_CDN_CHECK_ALIBABA_ACCESS_KEY") == "" || os.Getenv("GO_CDN_CHECK_ALIBABA_SECRET_KEY") == "" {
		t.Skip("set Alibaba CDN credentials to run this integration test")
	}
	if err := alibabaMain(); err != nil {
		t.Errorf("_main() error = %v", err)
	}
}
