# Compressor

**Compressor** is a lightweight middleware for compressing response data. It can be used with `gin` and `orbit` frameworks to improve network transmission efficiency.

`Compressor` utilizes the `CodecWriter` interface to compress response data. It currently supports the `gzip` and `deflate` algorithms. You can also implement your own `CodecWriter` to support other compression algorithms.

`Compressor` is built on top of the Go standard library and other powerful packages:

-   [gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)
-   [compress/gzip](https://pkg.go.dev/compress/gzip)
-   [compress/flate](https://pkg.go.dev/compress/flate)

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/compressor
```

## Quick Start

### Config

The `Compressor` has a config object that can be used to configure the batch process behavior. The config object provides the following methods for configuration:

-   `WithCompressLevel`: Sets the compression level. The default level is `6`.
-   `WithWriterCreateFunc`: Sets the writer create function. The default function is `DefaultWriterCreateFunc`.
-   `WithMatchFunc`: Sets the match function. The default function is `DefaultLimitMatchFunc`.
-   `WithIpWhitelist`: Sets the IP whitelist. The default whitelist is `DefaultIpWhitelist`.

### Compressor

#### 1. GZip

The `gzip` algorithm is the default algorithm used by the `Compressor`.

**Methods**

-   `HandlerFunc`: Returns a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop`: Stops the compressor. This is an empty function and does not need to be called.

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
	// 测试URL
	// Test URL
	testUrl = "/test"
)

// testRequestFunc 是一个测试请求的函数
// testRequestFunc is a function to test the request
func testRequestFunc(idx int, router *gin.Engine, conf *cr.Config, url string) {
	// 创建一个新的请求
	// Create a new request
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// 创建一个新的响应记录器
	// Create a new response recorder
	resp := httptest.NewRecorder()

	// 使用路由器处理HTTP请求
	// Use the router to handle the HTTP request
	router.ServeHTTP(resp, req)

	// 获取响应体内容
	// Get the content of the response body
	bodyContent := resp.Body.String()

	// 创建一个新的gzip读取器
	// Create a new gzip reader
	gr, _ := gzip.NewReader(resp.Body)
	defer gr.Close()

	// 读取响应的全部内容
	// Read all the content of the response
	plaintext, _ := io.ReadAll(gr)

	// 打印请求的信息
	// Print the information of the request
	fmt.Println("[Request]", idx, resp.Code, url, string(plaintext), bodyContent)
}

func main() {
	// 创建一个新的配置
	// Create a new configuration
	conf := cr.NewConfig()

	// 创建一个新的压缩器
	// Create a new compressor
	compr := cr.NewCompressor(conf)

	// 在函数返回时停止压缩器
	// Stop the compressor when the function returns
	defer compr.Stop()

	// 创建一个新的路由器
	// Create a new router
	router := gin.New()

	// 使用压缩器的处理函数
	// Use the handler function of the compressor
	router.Use(compr.HandlerFunc())

	// 添加一个GET路由
	// Add a GET route
	router.GET(testUrl, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 测试 10 次请求
	// Test 10 requests
	for i := 0; i < 10; i++ {
		// 调用测试请求函数，传入索引、路由器、配置和测试URL
		// Call the test request function, passing in the index, router, configuration, and test URL
		testRequestFunc(i, router, conf, testUrl)
	}

	// 等待所有请求任务执行完毕
	// Wait for all request tasks to complete
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

-   `HandlerFunc`: Returns a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop`: Stops the compressor. It is an empty function and does not need to be called.

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
	// 测试URL
	// Test URL
	testUrl = "/test"
)

// testNewDeflateWriterFunc 是一个创建新的DeflateWriter的函数
// testNewDeflateWriterFunc is a function to create a new DeflateWriter
func testNewDeflateWriterFunc(config *cr.Config, rw gin.ResponseWriter) any {
	return cr.NewDeflateWriter(config, rw)
}

// testRequestFunc 是一个测试请求的函数
// testRequestFunc is a function to test the request
func testRequestFunc(idx int, router *gin.Engine, conf *cr.Config, url string) {
	// 创建一个新的请求
	// Create a new request
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// 创建一个新的记录器
	// Create a new recorder
	resp := httptest.NewRecorder()

	// 服务HTTP
	// Serve HTTP
	router.ServeHTTP(resp, req)

	// 获取响应体内容
	// Get the body content of the response
	bodyContent := resp.Body.String()

	// 创建一个新的flate读取器
	// Create a new flate reader
	gr := flate.NewReader(resp.Body)
	defer gr.Close()

	// 读取所有的明文
	// Read all plaintext
	plaintext, _ := io.ReadAll(gr)

	// 打印请求信息
	// Print request information
	fmt.Println("[Request]", idx, resp.Code, url, string(plaintext), bodyContent)
}

func main() {
	// 创建新的配置
	// Create new configuration
	conf := cr.NewConfig().WithWriterCreateFunc(testNewDeflateWriterFunc)

	// 创建新的压缩器
	// Create new compressor
	compr := cr.NewCompressor(conf)

	// 停止压缩器
	// Stop the compressor
	defer compr.Stop()

	// 创建新的路由器
	// Create new router
	router := gin.New()

	// 使用压缩器的处理函数
	// Use the handler function of the compressor
	router.Use(compr.HandlerFunc())

	// 添加GET路由
	// Add GET route
	router.GET(testUrl, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 测试 10 次请求
	// Test 10 requests
	for i := 0; i < 10; i++ {
		// 调用测试请求函数，传入索引、路由器、配置和测试URL
		// Call the test request function, passing in the index, router, configuration, and test URL
		testRequestFunc(i, router, conf, testUrl)
	}

	// 等待所有请求任务执行完毕
	// Wait for all request tasks to complete
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
