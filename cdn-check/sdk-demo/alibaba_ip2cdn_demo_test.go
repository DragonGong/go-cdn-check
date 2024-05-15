package sdk_demo

import "testing"

func Test__main(t *testing.T) {
	if err := alibabaMain(); err != nil {
		t.Errorf("_main() error = %v", err)
	}
}
