package cdn_check

import (
	"fmt"
	"github.com/dragongong/go-cdn-check/logger"
	"github.com/projectdiscovery/cdncheck"
	"net"
)

type CdnChecker struct {
	checkers []Checker
}
type Checker interface {
	CheckConfig() bool
	check(ips []string) (map[string]bool, error)
	CdnFilter(ip []string) []string
	GetType() string
}

func NewCdnChecker() *CdnChecker {
	cdnChecker := &CdnChecker{}
	alibabaChecker := NewAlibabaChecker()
	baiduChecker := NewBaiduChecker()
	wangsuChecker := NewWangsuChecker()
	githubChecker := NewGithubChecker()
	cdnChecker.Register(alibabaChecker)
	cdnChecker.Register(baiduChecker)
	cdnChecker.Register(wangsuChecker)
	cdnChecker.Register(githubChecker)
	return cdnChecker
}
func (c *CdnChecker) Register(checker Checker) {
	if checker == nil {
		return
	}
	if !checker.CheckConfig() {
		logger.Log.Info("register warning: config is not right", "checker name", checker.GetType())
		return
	}
	if c.checkers == nil {
		c.checkers = make([]Checker, 0)
	}
	logger.Log.Info("cdn register success", "checker name", checker.GetType())
	c.checkers = append(c.checkers, checker)
}

// 判断cdn的err在这里处理
func (c *CdnChecker) CdnCheckFilter(ips []string) map[string]bool {
	if len(ips) == 0 {
		return nil
	}
	ipCdnMap := make(map[string]bool)
	for _, ip := range ips {
		ipCdnMap[ip] = true
	}
	for i := 0; i < len(c.checkers); i++ {
		if ips == nil {
			return ipCdnMap
		}
		ips = c.checkers[i].CdnFilter(ips)
	}
	for _, ip := range ips {
		if _, ok := ipCdnMap[ip]; ok {
			ipCdnMap[ip] = false
		}
	}

	var cdnList []string
	for ip, isCdn := range ipCdnMap {
		if isCdn {
			cdnList = append(cdnList, ip)
		}
	}
	logger.Log.Info("the ip addresses of the cdn the checkers recognized :", "cdnList:", fmt.Sprint(cdnList))
	return ipCdnMap
}

func (c *CdnChecker) CdnCheckOversea(ips []string) map[string]bool {
	client := cdncheck.New()
	res := make(map[string]bool)
	for _, ip := range ips {
		res[ip] = false
	}

	for _, ip := range ips {
		if matched, _, _ := client.CheckCDN(net.ParseIP(ip)); matched {
			res[ip] = true
		}
	}
	return res
}
