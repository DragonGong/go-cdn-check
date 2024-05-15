package cdn_check

import (
	"fmt"
	"github.com/clouddragongong/go-cdn-check/cdn-check/config"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

func Test_CdnCheck(t *testing.T) {
	config.CdnAkSkConfigPath = "./config/fig.yaml"
	config.ParseCdnConfig()
	cdnChecker := NewCdnChecker()
	//res := cdnChecker.CdnCheckFilter([]string{"116.114.98.35", "59.24.3.174",})
	p := time.Now()
	containCdn := []string{"116.114.98.35", "59.24.3.174"}
	cdnChecker.CdnCheckFilter(append(
		generateUniqueIPs(300), containCdn...))
	cost := time.Now().Sub(p)
	fmt.Println(cost)
}

func Test_CdnCheckOversea(t *testing.T) {
	cdnChecker := NewCdnChecker()
	p := time.Now()
	containCdn := []string{"120.52.22.96", "59.24.3.174", "180.163.57.128"}
	res := cdnChecker.CdnCheckOversea(append(
		generateUniqueIPs(300), containCdn...))
	fmt.Println(time.Now().Sub(p))
	//fmt.Println(res)
	count := 0
	var cdn []string
	for k, v := range res {
		if v {
			cdn = append(cdn, k)
			count++
		}
	}
	fmt.Printf("the total ip len is %d, the cdn ip len is %d , the cdn ip is  %v", len(res), count, cdn)
}

// generateRandomIP 生成一个随机的IP地址，并检查它是否已存在于提供的map中
func generateRandomIP(usedIPs map[string]struct{}, mu *sync.Mutex) (string, bool) {
	mu.Lock()
	defer mu.Unlock()

	for {
		// 随机生成四个数字作为IP地址的各个部分
		var ipParts [4]byte
		for i := range ipParts {
			ipParts[i] = byte(rand.Intn(256))
		}

		// 将四个数字组合成一个IP地址的字符串表示
		ip := net.IPv4(ipParts[0], ipParts[1], ipParts[2], ipParts[3]).String()

		// 检查这个IP是否已经被使用过
		if _, exists := usedIPs[ip]; !exists {
			// 如果没有被使用过，则添加到已使用列表中，并返回该IP
			usedIPs[ip] = struct{}{}
			return ip, true
		}
	}
}

// generateUniqueIPs 生成指定数量的唯一随机IP地址
func generateUniqueIPs(n int) []string {
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机数种子
	usedIPs := make(map[string]struct{})
	ips := make([]string, 0, n)

	var mu sync.Mutex // 用于同步访问usedIPs map

	// 不断尝试生成新的IP地址，直到我们有了n个不重复的IP
	for len(ips) < n {
		ip, ok := generateRandomIP(usedIPs, &mu)
		if ok {
			ips = append(ips, ip)
		}
	}

	return ips
}
