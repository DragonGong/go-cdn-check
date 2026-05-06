package gocdncheck

import (
	cdncheck "github.com/clouddragongong/go-cdn-check/cdn-check"
	"github.com/clouddragongong/go-cdn-check/cdn-check/config"
)

// ProviderCredentials contains the access key pair for a single CDN provider.
type ProviderCredentials struct {
	AccessKey string
	SecretKey string
}

// Credentials groups provider credentials used by provider-backed checks.
type Credentials struct {
	Baidu   ProviderCredentials
	Alibaba ProviderCredentials
	Wangsu  ProviderCredentials
}

// Result maps an IP to whether it is recognized as a CDN IP.
type Result map[string]bool

// CDNIPs returns the IPs recognized as CDN endpoints.
func (r Result) CDNIPs() []string {
	var ips []string
	for ip, isCDN := range r {
		if isCDN {
			ips = append(ips, ip)
		}
	}
	return ips
}

// OriginIPs returns the IPs not recognized as CDN endpoints.
func (r Result) OriginIPs() []string {
	var ips []string
	for ip, isCDN := range r {
		if !isCDN {
			ips = append(ips, ip)
		}
	}
	return ips
}

// Detector is the public detector type exposed by the root package.
type Detector struct {
	checker *cdncheck.CdnChecker
}

// New creates a detector using the credentials currently loaded in process memory.
func New() *Detector {
	return &Detector{checker: cdncheck.NewCdnChecker()}
}

// NewWithConfig loads explicit credentials and creates a detector from them.
func NewWithConfig(cfg Credentials) *Detector {
	SetCredentials(cfg)
	return New()
}

// NewFromEnv loads provider credentials from environment variables and creates a detector.
func NewFromEnv() *Detector {
	LoadConfigFromEnv()
	return New()
}

// NewFromConfigFile loads provider credentials from a local config file and creates a detector.
func NewFromConfigFile(path string) (*Detector, error) {
	if err := LoadConfig(path); err != nil {
		return nil, err
	}
	return New(), nil
}

// LoadConfig loads provider credentials from the given YAML file.
func LoadConfig(path string) error {
	config.CdnAkSkConfigPath = path
	return config.ParseCdnConfig()
}

// LoadConfigFromEnv loads provider credentials from environment variables.
func LoadConfigFromEnv() {
	config.SetConfig(config.LoadFromEnv())
}

// SetCredentials replaces the in-process provider credentials used by new detectors.
func SetCredentials(cfg Credentials) {
	config.SetConfig(config.CdnProviderAuth{
		Baidu: config.BaiduCdn{
			AccessKey: cfg.Baidu.AccessKey,
			SecretKey: cfg.Baidu.SecretKey,
		},
		Alibaba: config.AlibabaCdn{
			AccessKey: cfg.Alibaba.AccessKey,
			SecretKey: cfg.Alibaba.SecretKey,
		},
		Wangsu: config.WangsuCdn{
			AccessKey: cfg.Wangsu.AccessKey,
			SecretKey: cfg.Wangsu.SecretKey,
		},
	})
}

// CurrentCredentials returns the credentials currently loaded in process memory.
func CurrentCredentials() Credentials {
	cfg := config.GetConfig()
	return Credentials{
		Baidu: ProviderCredentials{
			AccessKey: cfg.Baidu.AccessKey,
			SecretKey: cfg.Baidu.SecretKey,
		},
		Alibaba: ProviderCredentials{
			AccessKey: cfg.Alibaba.AccessKey,
			SecretKey: cfg.Alibaba.SecretKey,
		},
		Wangsu: ProviderCredentials{
			AccessKey: cfg.Wangsu.AccessKey,
			SecretKey: cfg.Wangsu.SecretKey,
		},
	}
}

// Check reports whether each IP is recognized as a CDN IP.
func (d *Detector) Check(ips []string) Result {
	return Result(d.checker.CdnCheckFilter(ips))
}

// CheckOversea performs detection using the projectdiscovery/cdncheck dataset.
func (d *Detector) CheckOversea(ips []string) Result {
	return Result(d.checker.CdnCheckOversea(ips))
}

// CdnCheckFilter is kept for backward compatibility. Prefer Check.
func (d *Detector) CdnCheckFilter(ips []string) Result {
	return d.Check(ips)
}

// CdnCheckOversea is kept for backward compatibility. Prefer CheckOversea.
func (d *Detector) CdnCheckOversea(ips []string) Result {
	return d.CheckOversea(ips)
}

// Check is a convenience helper that creates a default detector and runs Check.
func Check(ips []string) Result {
	return New().Check(ips)
}

// CheckOversea is a convenience helper that creates a default detector and runs CheckOversea.
func CheckOversea(ips []string) Result {
	return New().CheckOversea(ips)
}
