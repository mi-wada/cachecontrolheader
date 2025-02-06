# cachecontrolheader

A Go library to parse HTTP Cache-Control headers based on [RFC9111](https://datatracker.ietf.org/doc/html/rfc9111.html).

## Usage

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/mi-wada/cachecontrolheader"
)

func main() {
	res, err := http.Get("https://example.com")
	if err != nil {
		panic(err)
	}
	cacheControl, err := cachecontrolheader.Parse(res.Header.Get("Cache-Control"))
	if err != nil {
		panic(err)
	}
	fmt.Println(cacheControl.MaxAge)
	fmt.Println(cacheControl.Private)
  fmt.Println(cacheControl.MustRevalidate)
}
```

## Install

```shell
go get github.com/mi-wada/cachecontrolheader@latest
```
