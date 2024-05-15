package config

import (
	"github.com/clouddragongong/go-cdn-check/logger"
	"gopkg.in/yaml.v2"
	"os"
)

var CdnConfig = &CdnProviderAuth{}

type CdnProviderAuth struct {
	Baidu   BaiduCdn   `yaml:"Baidu"`
	Alibaba AlibabaCdn `yaml:"Alibaba"`
	Wangsu  WangsuCdn  `yaml:"Wangsu"`
}
type BaiduCdn struct {
	AccessKey string `yaml:"AccessKey"`
	SecretKey string `yaml:"SecretKey"`
}

type AlibabaCdn struct {
	AccessKey string `yaml:"AccessKey"`
	SecretKey string `yaml:"SecretKey"`
}

type WangsuCdn struct {
	AccessKey string `yaml:"AccessKey"`
	SecretKey string `yaml:"SecretKey"`
}

var (
	CdnAkSkConfigPath    string
	DefaultCdnConfigPath = "./cdn-check/config/config.yaml"
)

func init() {
	ParseCdnConfig()
}

func ParseCdnConfig() {
	if CdnAkSkConfigPath == "" {
		CdnAkSkConfigPath = DefaultCdnConfigPath
		logger.Log.Warn("read cdn config warning:", "reason", "cdn auth config.yaml path is nil")
	}
	dataBytes, err := os.ReadFile(CdnAkSkConfigPath)
	if err != nil {
		logger.Log.Error("read cdn config error:", "reason:", err)
	}
	err = yaml.Unmarshal(dataBytes, CdnConfig)
	if err != nil {
		logger.Log.Error("read cdn config error , unmarshal", "reason", err)
	}
}
