package cdn_check

import (
	"github.com/clouddragongong/go-cdn-check/logger"
	"sync"
	"time"
)

type CdnCheckerManager struct {
	checkerName string
	Signal      chan struct{}
	Qps         int //每秒
	MaxIpsNum   int //最大的ip数
}

func NewCdnCheckerManager(name string, signal chan struct{}, qps int, maxIpsNum int) *CdnCheckerManager {
	manager := &CdnCheckerManager{
		checkerName: name,
		Signal:      signal,
		Qps:         qps,
		MaxIpsNum:   maxIpsNum,
	}
	go manager.addQpsBySecond()
	return manager
}

// 由于非阻塞,只需要限制qps,不需要限制min
func (c *CdnCheckerManager) addQpsBySecond() {
	for {
		time.Sleep(time.Millisecond * time.Duration(1000/c.Qps))
		c.Signal <- struct{}{}
	}
}

// CdnFilter 控制全局的并发
func (c *CdnCheckerManager) CdnFilter(ips []string, fn func(ip []string) (map[string]bool, error)) []string {
	// fn代表着每个企业的check函数
	initialLength := len(ips)
	timeBegin := time.Now()
	ipsInputs := make([][]string, 0)
	for {
		if len(ips) < c.MaxIpsNum && len(ips) > 0 {
			ipsInputs = append(ipsInputs, ips)
			break
		} else if len(ips) <= 0 {
			break
		}
		ipsInputs = append(ipsInputs, ips[:c.MaxIpsNum])
		ips = ips[c.MaxIpsNum:]
	}
	index := 0
	success := 0
	wg := &sync.WaitGroup{}
	resultC := make(chan []string, len(ipsInputs))
	for range c.Signal {
		if index >= len(ipsInputs) {
			break
		}
		wg.Add(1)
		go func(index int) {
			resMap, err := fn(ipsInputs[index])
			if err != nil {
				resultC <- ipsInputs[index]
				logger.Log.Info("gomap CdnFilter 接口出现了问题,识别为不是cdn", "err", err, "checker name", c.checkerName)
				wg.Done()
				return
			}
			success++
			var resList []string
			for ip, isCdn := range resMap {
				if !isCdn {
					resList = append(resList, ip)
				}
			}
			resultC <- resList
			wg.Done()
		}(index)
		index++
	}
	go func() {
		wg.Wait()
		close(resultC)
	}()
	var res []string
	for ipsNotCdn := range resultC {
		if ipsNotCdn == nil {
			continue
		}
		res = append(res, ipsNotCdn...)
	}
	timeCostSecond := time.Now().Sub(timeBegin).Seconds()
	logger.Log.Info("[gomap][CdnFilter]识别了ip:", "原本的总ip数：", initialLength, "识别为cdn的ip数：", initialLength-len(res), "success number", success, "total number", len(ipsInputs), "checker:", c.checkerName, "const time", timeCostSecond)
	return res
	//for {
	//	select {
	//	case <-c.Signal:
	//
	//		if len(ips) < c.MaxIpsNum && len(ips) > 0 {
	//			ipCdnMapRes, err := fn(ips)
	//			if err != nil {
	//				// 接口出现了问题，那就是加入res
	//				res = append(res, ips...)
	//				logger.Log.Info("gomap CdnFilter 接口出现了问题,识别为不是cdn", "err", err)
	//			} else {
	//				res = append(res, c.IpCdnMapToSlice(ipCdnMapRes)...)
	//			}
	//			logger.Log.Info("[gomap][CdnFilter]本组识别了ip:", "原本的总ip数：", initialLength, "识别为cdn的ip数：", initialLength-len(res))
	//			return
	//		} else if len(ips) <= 0 {
	//			logger.Log.Info("[gomap][CdnFilter]本组识别了ip:", "原本的总ip数：", initialLength, "识别为cdn的ip数：", initialLength-len(res))
	//			return
	//		} else {
	//			ipCdnMapRes, err := fn(ips[:c.MaxIpsNum])
	//			if err != nil {
	//				// 接口出现了问题，那就是加入res
	//				res = append(res, ips[:c.MaxIpsNum]...)
	//				logger.Log.Info("gomap CdnFilter 接口出现了问题,识别为不是cdn", "err", err)
	//			} else {
	//				res = append(res, c.IpCdnMapToSlice(ipCdnMapRes)...)
	//			}
	//		}
	//		ips = ips[c.MaxIpsNum:]
	//	}
	//}
}

func (c *CdnCheckerManager) IpCdnMapToSlice(m map[string]bool) []string {
	var result []string
	for key, value := range m {
		if !value { // 如果值为false
			result = append(result, key) // 将键添加到结果切片中
		}
	}
	return result // 返回结果切片
}
