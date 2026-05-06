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

	fmt.Println("all results:", result)
	fmt.Println("cdn ips:", result.CDNIPs())
	fmt.Println("origin ips:", result.OriginIPs())
}
