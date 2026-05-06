package gocdncheck

import "testing"

func TestSetCredentials(t *testing.T) {
	previous := CurrentCredentials()
	t.Cleanup(func() {
		SetCredentials(previous)
	})

	cfg := Credentials{
		Baidu: ProviderCredentials{
			AccessKey: "baidu-ak",
			SecretKey: "baidu-sk",
		},
		Alibaba: ProviderCredentials{
			AccessKey: "ali-ak",
			SecretKey: "ali-sk",
		},
		Wangsu: ProviderCredentials{
			AccessKey: "ws-ak",
			SecretKey: "ws-sk",
		},
	}

	SetCredentials(cfg)
	current := CurrentCredentials()

	if current != cfg {
		t.Fatalf("credentials mismatch: got %#v want %#v", current, cfg)
	}
}

func TestCheckReturnsResults(t *testing.T) {
	previous := CurrentCredentials()
	SetCredentials(Credentials{})
	t.Cleanup(func() {
		SetCredentials(previous)
	})

	result := Check([]string{"120.52.22.96", "8.8.8.8"})
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestResultHelpers(t *testing.T) {
	result := Result{
		"1.1.1.1": true,
		"8.8.8.8": false,
		"9.9.9.9": true,
	}

	if len(result.CDNIPs()) != 2 {
		t.Fatalf("expected 2 cdn ips, got %d", len(result.CDNIPs()))
	}
	if len(result.OriginIPs()) != 1 {
		t.Fatalf("expected 1 origin ip, got %d", len(result.OriginIPs()))
	}
}
