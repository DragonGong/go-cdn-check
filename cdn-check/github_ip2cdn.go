package cdn_check

import (
	"github.com/projectdiscovery/cdncheck"
	"net"
)

type GithubChecker struct {
	client            *cdncheck.Client
	cdnCheckerManager *CdnCheckerManager
}

func NewGithubChecker() Checker {
	return &GithubChecker{
		client:            cdncheck.New(),
		cdnCheckerManager: NewCdnCheckerManager("github", githubQps, githubMaxIpsNums),
	}
}
func (g GithubChecker) CheckConfig() bool {
	return true
}

func (g GithubChecker) check(ips []string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, ip := range ips {
		res[ip] = false
	}
	for _, ip := range ips {
		if matched, _, _ := g.client.CheckCDN(net.ParseIP(ip)); matched {
			res[ip] = true
		}
	}
	return res, nil
}

func (g GithubChecker) CdnFilter(ip []string) []string {
	return g.cdnCheckerManager.CdnFilter(ip, g.check)
}

func (g GithubChecker) GetType() string {
	return g.cdnCheckerManager.checkerName
}
