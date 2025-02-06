package cachecontrolheader_test

import (
	"fmt"

	"github.com/mi-wada/cachecontrolheader"
)

func Example() {
	s := "max-age=3600, must-revalidate, private"
	h, err := cachecontrolheader.Parse(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(h.MaxAge, h.MustRevalidate, h.Private)
	// Output: 1h0m0s true true
}
