# cachecontrolheader

[![Godoc Reference](https://godoc.org/github.com/mi-wada/cachecontrolheader?status.svg)](http://godoc.org/github.com/mi-wada/cachecontrolheader)

A Go package to parse HTTP Cache-Control headers based on [RFC9111](https://datatracker.ietf.org/doc/html/rfc9111.html).

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
	cacheControl := cachecontrolheader.Parse(res.Header.Get("Cache-Control"))
	fmt.Println(cacheControl.MaxAge)
	fmt.Println(cacheControl.MustRevalidate)
	fmt.Println(cacheControl.Private)
}
```

## Install

```shell
go get github.com/mi-wada/cachecontrolheader@latest
```
