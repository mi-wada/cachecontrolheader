package cachecontrolheader_test

import (
	"fmt"

	"github.com/mi-wada/cachecontrolheader"
)

func Example() {
	s := "max-age=3600, must-revalidate, private"
	h := cachecontrolheader.Parse(s)
	fmt.Println(h.MaxAge, h.MustRevalidate, h.Private, h.MaxStale)
	// Output: 1h0m0s true true <nil>
}

func Example_parseStrict() {
	s := "max-age=3600, must-revalidate, private, ???"
	_, err := cachecontrolheader.ParseStrict(s)
	fmt.Println(err)

	s = "max-age=invalid, must-revalidate, private"
	_, err = cachecontrolheader.ParseStrict(s)
	fmt.Println(err)

	// Output:
	// unknown directive: ???
	// failed to parse the value of directive(max-age=invalid): time: invalid duration "invalids"
}
