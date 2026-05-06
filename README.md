# go-cdn-check

`go-cdn-check` is a Go package for checking whether an IP belongs to a CDN provider.

It supports two usage modes:

- Zero-config offline detection via `projectdiscovery/cdncheck`
- Provider API based detection for Baidu CDN, Alibaba CDN, and Wangsu when credentials are supplied

The root package is the supported public API. The `cdn-check/...` subpackages are internal implementation details.

## Install

```bash
go get github.com/dragongong/go-cdn-check
```

## Quick Start

```go
package main

import (
	"fmt"

	gocdncheck "github.com/dragongong/go-cdn-check"
)

func main() {
	detector := gocdncheck.New()
	result := detector.Check([]string{
		"1.1.1.1",
		"8.8.8.8",
		"120.52.22.96",
	})
	fmt.Println(result)
	fmt.Println(result.CDNIPs())
}
```

## Use Provider Credentials

You can load credentials from environment variables:

```bash
export GO_CDN_CHECK_BAIDU_ACCESS_KEY="..."
export GO_CDN_CHECK_BAIDU_SECRET_KEY="..."
export GO_CDN_CHECK_ALIBABA_ACCESS_KEY="..."
export GO_CDN_CHECK_ALIBABA_SECRET_KEY="..."
export GO_CDN_CHECK_WANGSU_ACCESS_KEY="..."
export GO_CDN_CHECK_WANGSU_SECRET_KEY="..."
```

Then:

```go
detector := gocdncheck.NewFromEnv()
result := detector.Check([]string{"116.114.98.35"})
```

Or load a local config file:

```go
detector, err := gocdncheck.NewFromConfigFile("./cdn-check/config/config.local.yaml")
if err != nil {
	panic(err)
}
_ = detector
```

A sample config file is available at [cdn-check/config/config.example.yaml](./cdn-check/config/config.example.yaml).

Runnable examples are available in [examples/basic](./examples/basic) and [examples/from-env](./examples/from-env).

## Notes

- Public usage should depend on the root package `github.com/dragongong/go-cdn-check`.
- Prefer `Detector.Check` / `Detector.CheckOversea`; `CdnCheckFilter` and `CdnCheckOversea` are legacy compatibility names.
- Real credentials are no longer stored in the repository.
- Log files and local credential files are ignored by git.
- By default the package only logs to stdout. Set `GO_CDN_CHECK_LOG_DIR` if you want file logs.
