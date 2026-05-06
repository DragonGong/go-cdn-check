package main

import (
	"fmt"
	"log"

	gocdncheck "github.com/dragongong/go-cdn-check"
)

func main() {
	detector := gocdncheck.NewFromEnv()
	result := detector.Check([]string{"116.114.98.35"})

	if len(result) == 0 {
		log.Fatal("no results returned")
	}

	fmt.Println(result)
}
