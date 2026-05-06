package gocdncheck

import "fmt"

func ExampleCheck() {
	result := Check([]string{"120.52.22.96", "8.8.8.8"})
	fmt.Println(len(result) > 0)
	// Output: true
}
