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

func Example_errorOnUnknown() {
	s := "max-age=3600, must-revalidate, private, ???"
	_, err := cachecontrolheader.Parse(s, cachecontrolheader.ErrorOnUnknown())
	fmt.Println(err)
	// Output: unknown directive: ???
}

func Example_errorOnInvalidValues() {
	s := "max-age=invalid, must-revalidate, private"
	_, err := cachecontrolheader.Parse(s, cachecontrolheader.ErrorOnInvalidValues())
	fmt.Println(err)
	// Output: failed to parse duration for directive max-age=invalid: time: invalid duration "invalids"
}
