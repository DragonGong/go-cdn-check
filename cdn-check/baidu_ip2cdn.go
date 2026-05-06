package cdn_check

import (
	"errors"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/clouddragongong/go-cdn-check/cdn-check/config"
	"time"
)

var _ Checker = (*BaiduChecker)(nil)

type BaiduChecker struct {
	client            *cdn.Client
	cdnCheckerManager *CdnCheckerManager
}

func NewBaiduChecker() Checker {
	ak := config.CdnConfig.Baidu.AccessKey
	sk := config.CdnConfig.Baidu.SecretKey
	if ak == "" || sk == "" {
		return nil
	}
	endpoint := "https://cdn.baidubce.com"
	client, _ := cdn.NewClient(ak, sk, endpoint)
	return &BaiduChecker{
		client:            client,
		cdnCheckerManager: NewCdnCheckerManager("baidu", make(chan struct{}, baiduQps), baiduQps, baiduMaxIpsNums),
	}
}

func (c *BaiduChecker) CheckConfig() bool {
	return len(config.CdnConfig.Baidu.AccessKey) > 0 && len(config.CdnConfig.Baidu.SecretKey) > 0
}

func (c *BaiduChecker) GetType() string {
	return c.cdnCheckerManager.checkerName
}
func (c *BaiduChecker) check(ips []string) (map[string]bool, error) {
	if len(ips) > c.cdnCheckerManager.MaxIpsNum {
		return nil, errors.New("the number of IPs exceeds the limit")
	}
	ipsInfos, err := c.client.GetIpListInfo(ips, "describeIp")
	if err != nil {
		time.Sleep(time.Millisecond * time.Duration(1000/c.cdnCheckerManager.Qps))
		ipsInfos, err = c.client.GetIpListInfo(ips, "describeIp")
		if err != nil {
			return nil, err
		}
	}
	ipCdnMap := make(map[string]bool)
	// 得保证map的key是有原始数据过来的,传接口之后出现丢失的ip为false
	for _, ip := range ips {
		ipCdnMap[ip] = false
	}
	for _, ipsInfo := range ipsInfos {
		if _, ok := ipCdnMap[ipsInfo.IP]; ok {
			ipCdnMap[ipsInfo.IP] = ipsInfo.IsCdnIp
		}
	}
	return ipCdnMap, nil
}

func (c *BaiduChecker) CdnFilter(ips []string) []string {
	return c.cdnCheckerManager.CdnFilter(ips, c.check)
}
