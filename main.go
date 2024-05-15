package go_cdn_check

import (
	"fmt"
	cdn_check "github.com/clouddragongong/go-cdn-check/cdn-check"
)

func main() {
	ip := []string{"127.0.0.1"}
	checker := cdn_check.NewCdnChecker()
	m := checker.CdnCheckFilter(ip)
	fmt.Println(m)
}
