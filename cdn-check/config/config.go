package config

import (
	"errors"
	"github.com/clouddragongong/go-cdn-check/logger"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"sync"
)

var CdnConfig = &CdnProviderAuth{}
var configMu sync.RWMutex

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
	DefaultCdnConfigPath = "./cdn-check/config/config.local.yaml"
)

func init() {
	if err := ParseCdnConfig(); err != nil {
		logger.Log.Debug("load cdn config skipped", "reason", err)
	}
}

func ParseCdnConfig() error {
	path := strings.TrimSpace(CdnAkSkConfigPath)
	if path == "" {
		path = strings.TrimSpace(os.Getenv("GO_CDN_CHECK_CONFIG"))
	}
	if path == "" {
		path = DefaultCdnConfigPath
	}

	cfg := LoadFromEnv()
	if dataBytes, err := os.ReadFile(path); err == nil {
		if err := yaml.Unmarshal(dataBytes, &cfg); err != nil {
			return err
		}
		cfg = mergeEnvOverrides(cfg)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	CdnAkSkConfigPath = path
	SetConfig(cfg)
	return nil
}

func LoadFromEnv() CdnProviderAuth {
	return mergeEnvOverrides(CdnProviderAuth{})
}

func SetConfig(cfg CdnProviderAuth) {
	configMu.Lock()
	defer configMu.Unlock()
	*CdnConfig = cfg
}

func GetConfig() CdnProviderAuth {
	configMu.RLock()
	defer configMu.RUnlock()
	return *CdnConfig
}

func mergeEnvOverrides(cfg CdnProviderAuth) CdnProviderAuth {
	cfg.Baidu.AccessKey = envOr("GO_CDN_CHECK_BAIDU_ACCESS_KEY", cfg.Baidu.AccessKey)
	cfg.Baidu.SecretKey = envOr("GO_CDN_CHECK_BAIDU_SECRET_KEY", cfg.Baidu.SecretKey)
	cfg.Alibaba.AccessKey = envOr("GO_CDN_CHECK_ALIBABA_ACCESS_KEY", cfg.Alibaba.AccessKey)
	cfg.Alibaba.SecretKey = envOr("GO_CDN_CHECK_ALIBABA_SECRET_KEY", cfg.Alibaba.SecretKey)
	cfg.Wangsu.AccessKey = envOr("GO_CDN_CHECK_WANGSU_ACCESS_KEY", cfg.Wangsu.AccessKey)
	cfg.Wangsu.SecretKey = envOr("GO_CDN_CHECK_WANGSU_SECRET_KEY", cfg.Wangsu.SecretKey)
	return cfg
}

func envOr(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
