package gocdncheck

import (
	cdncheck "github.com/clouddragongong/go-cdn-check/cdn-check"
	"github.com/clouddragongong/go-cdn-check/cdn-check/config"
)

type Credentials = config.CdnProviderAuth
type Checker = cdncheck.CdnChecker

func New() *Checker {
	return cdncheck.NewCdnChecker()
}

func NewWithConfig(cfg Credentials) *Checker {
	config.SetConfig(cfg)
	return New()
}

func LoadConfig(path string) error {
	config.CdnAkSkConfigPath = path
	return config.ParseCdnConfig()
}

func LoadConfigFromEnv() {
	config.SetConfig(config.LoadFromEnv())
}

func SetCredentials(cfg Credentials) {
	config.SetConfig(cfg)
}

func Check(ips []string) map[string]bool {
	return New().CdnCheckFilter(ips)
}

func CheckOversea(ips []string) map[string]bool {
	return New().CdnCheckOversea(ips)
}
