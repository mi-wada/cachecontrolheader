package cachecontrolheader_test

import (
	"fmt"

	"github.com/mi-wada/cachecontrolheader"
)

func Example() {
	s := "max-age=3600, private, must-revalidate"
	h, err := cachecontrolheader.Parse(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(h.MaxAge, h.Private, h.MustRevalidate)
	// Output: 1h0m0s true true
}
