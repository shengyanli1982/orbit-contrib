# Compressor

**Compressor** is a simple middleware for compressing response data. It can be used in `Gin` and `orbit`. And designed to be used in improve network transmission efficiency.

`Compressor` use interface `CodecWriter` to compress response data. Now it supports `gzip` and `deflate` algorithm. You can also implement your own `CodecWriter` to support other algorithm.

`Compressor` based golang native package and other powerful packages:

-   [gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)
-   [compress/gzip](https://pkg.go.dev/compress/gzip)
-   [compress/flate](https://pkg.go.dev/compress/flate)

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/compressor
```

## Quick Start

`Compressor` is designed to be used in improve network transmission efficiency. It is recommended to use it with `orbit` or `gin`.

### Config

`Compressor` has a config object, which can be used to configure the batch process behavior. The config object can be used following methods to set.

-   `WithCompressLevel` : set the compress level. The default is `6`.
-   `WithWriterCreateFunc` : set the writer create function. The default is `DefaultWriterCreateFunc`.
-   `WithMatchFunc` : set the match function. The default is `DefaultLimitMatchFunc`.
-   `WithIpWhitelist` : set the whitelist. The default is `DefaultIpWhitelist`.

### Compressor

#### 1. GZip

`gzip` is default algorithm for `Compressor`.

**Methods**

-   `HandlerFunc` : return a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop` : stop the compressor. It is empty function, no need to call it.

**Example**

```go
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	cr "github.com/shengyanli1982/orbit-contrib/pkg/compressor"
)

var (
	testUrl = "/test"
)

func testRequestFunc(idx int, router *gin.Engine, conf *cr.Config, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Get the response body
	bodyContent := resp.Body.String()

	// Create gzip reader
	gr, _ := gzip.NewReader(resp.Body)
	defer gr.Close()

	// Read the response
	plaintext, _ := io.ReadAll(gr)

	// Print the request information
	fmt.Println("[Request]", idx, resp.Code, url, string(plaintext), bodyContent)
}

func main() {
	// Create a new rate limiter
	conf := cr.NewConfig()
	compr := cr.NewCompressor(conf)
	defer compr.Stop()

	// Create a test context
	router := gin.New()
	router.Use(compr.HandlerFunc())
	router.GET(testUrl, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		testRequestFunc(i, router, conf, testUrl)
	}

	// Wait for to complete
	time.Sleep(time.Second)
}
```

**Result**

```bash
$ go run demo.go
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /test                     --> main.main.func1 (2 handlers)
[Request] 0 200 /test OK �����-�6�
[Request] 1 200 /test OK �����-�6�
[Request] 2 200 /test OK �����-�6�
[Request] 3 200 /test OK �����-�6�
[Request] 4 200 /test OK �����-�6�
[Request] 5 200 /test OK �����-�6�
[Request] 6 200 /test OK �����-�6�
[Request] 7 200 /test OK �����-�6�
[Request] 8 200 /test OK �����-�6�
[Request] 9 200 /test OK �����-�6�
```

#### 2. Deflate

**Methods**

-   `HandlerFunc` : return a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop` : stop the compressor. It is empty function, no need to call it.

**Example**

```go
package main

import (
	"compress/flate"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	cr "github.com/shengyanli1982/orbit-contrib/pkg/compressor"
)

var (
	testUrl = "/test"
)

func testNewDeflateWriterFunc(config *cr.Config, rw gin.ResponseWriter) any {
	return cr.NewDeflateWriter(config, rw)
}

func testRequestFunc(idx int, router *gin.Engine, conf *cr.Config, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Get the response body
	bodyContent := resp.Body.String()

	// Create gzip reader
	gr := flate.NewReader(resp.Body)
	defer gr.Close()

	// Read the response
	plaintext, _ := io.ReadAll(gr)

	// Print the request information
	fmt.Println("[Request]", idx, resp.Code, url, string(plaintext), bodyContent)
}

func main() {
	// Create a new rate limiter, use deflate writer create function
	conf := cr.NewConfig().WithWriterCreateFunc(testNewDeflateWriterFunc)
	compr := cr.NewCompressor(conf)
	defer compr.Stop()

	// Create a test context
	router := gin.New()
	router.Use(compr.HandlerFunc())
	router.GET(testUrl, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		testRequestFunc(i, router, conf, testUrl)
	}

	// Wait for to complete
	time.Sleep(time.Second)
}
```

**Result**

```bash
$ go run demo.go
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /test                     --> main.main.func1 (2 handlers)
[Request] 0 200 /test OK ����
[Request] 1 200 /test OK ����
[Request] 2 200 /test OK ����
[Request] 3 200 /test OK ����
[Request] 4 200 /test OK ����
[Request] 5 200 /test OK ����
[Request] 6 200 /test OK ����
[Request] 7 200 /test OK ����
[Request] 8 200 /test OK ����
[Request] 9 200 /test OK ����
```
