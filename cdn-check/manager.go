package cdn_check

import (
	"github.com/clouddragongong/go-cdn-check/logger"
	"sync"
	"sync/atomic"
	"time"
)

type CdnCheckerManager struct {
	checkerName string
	Qps         int //每秒
	MaxIpsNum   int //最大的ip数
}

func NewCdnCheckerManager(name string, qps int, maxIpsNum int) *CdnCheckerManager {
	return &CdnCheckerManager{
		checkerName: name,
		Qps:         qps,
		MaxIpsNum:   maxIpsNum,
	}
}

// CdnFilter 控制全局的并发
func (c *CdnCheckerManager) CdnFilter(ips []string, fn func(ip []string) (map[string]bool, error)) []string {
	if len(ips) == 0 {
		return nil
	}

	timeBegin := time.Now()
	initialLength := len(ips)
	ipsInputs := chunkIPs(ips, c.MaxIpsNum)
	resultC := make(chan []string, len(ipsInputs))
	wg := &sync.WaitGroup{}
	var success int64

	interval := time.Second
	if c.Qps > 0 {
		interval = time.Second / time.Duration(c.Qps)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for index, batch := range ipsInputs {
		if index > 0 {
			<-ticker.C
		}
		wg.Add(1)
		batch := append([]string(nil), batch...)
		go func(currentBatch []string) {
			defer wg.Done()
			resMap, err := fn(currentBatch)
			if err != nil {
				resultC <- currentBatch
				logger.Log.Info("gomap CdnFilter 接口出现了问题,识别为不是cdn", "err", err, "checker name", c.checkerName)
				return
			}
			atomic.AddInt64(&success, 1)
			var resList []string
			for ip, isCdn := range resMap {
				if !isCdn {
					resList = append(resList, ip)
				}
			}
			resultC <- resList
		}(batch)
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
	timeCostSecond := time.Since(timeBegin).Seconds()
	logger.Log.Info("[gomap][CdnFilter]识别了ip:", "原本的总ip数：", initialLength, "识别为cdn的ip数：", initialLength-len(res), "success number", success, "total number", len(ipsInputs), "checker:", c.checkerName, "const time", timeCostSecond)
	return res
}

func chunkIPs(ips []string, maxBatchSize int) [][]string {
	if len(ips) == 0 {
		return nil
	}
	if maxBatchSize <= 0 || maxBatchSize >= len(ips) {
		return [][]string{append([]string(nil), ips...)}
	}

	batches := make([][]string, 0, (len(ips)+maxBatchSize-1)/maxBatchSize)
	for start := 0; start < len(ips); start += maxBatchSize {
		end := start + maxBatchSize
		if end > len(ips) {
			end = len(ips)
		}
		batches = append(batches, append([]string(nil), ips[start:end]...))
	}
	return batches
}
