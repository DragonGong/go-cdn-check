# go-cdn-check

`go-cdn-check` is a Go package for checking whether an IP belongs to a CDN provider.

It supports two usage modes:

- Zero-config offline detection via `projectdiscovery/cdncheck`
- Provider API based detection for Baidu CDN, Alibaba CDN, and Wangsu when credentials are supplied

## Install

```bash
go get github.com/clouddragongong/go-cdn-check
```

## Quick Start

```go
package main

import (
	"fmt"

	gocdncheck "github.com/clouddragongong/go-cdn-check"
)

func main() {
	result := gocdncheck.Check([]string{
		"1.1.1.1",
		"8.8.8.8",
		"120.52.22.96",
	})
	fmt.Println(result)
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
gocdncheck.LoadConfigFromEnv()
checker := gocdncheck.New()
result := checker.CdnCheckFilter([]string{"116.114.98.35"})
```

Or load a local config file:

```go
err := gocdncheck.LoadConfig("./cdn-check/config/config.local.yaml")
if err != nil {
	panic(err)
}
checker := gocdncheck.New()
```

A sample config file is available at [cdn-check/config/config.example.yaml](./cdn-check/config/config.example.yaml).

## Notes

- Real credentials are no longer stored in the repository.
- Log files and local credential files are ignored by git.
- By default the package only logs to stdout. Set `GO_CDN_CHECK_LOG_DIR` if you want file logs.
