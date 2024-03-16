# Rewriter

**Rewriter** is a lightweight middleware for rewriting request URL paths and returning a status code of `http.StatusTemporaryRedirect`. It can be seamlessly integrated with `gin` and `orbit` frameworks. This middleware is designed for redirect requests that require URL rewriting.

Why use `http.StatusTemporaryRedirect` instead of `http.StatusMovedPermanently`? By using `http.StatusTemporaryRedirect`, the response will not be cached by the browser, and subsequent requests will not be sent to the original URL again.

`Rewriter` is built on top of powerful Golang native packages and other libraries:

-   [gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/rewriter
```

## Quick Start

### Configuration

The `Rewriter` middleware provides a configuration object that allows you to customize its behavior. The configuration object offers the following methods for customization:

- `WithCallback`: Sets the callback function. The default is `&emptyCallback{}`.
- `WithPathRewriteFunc`: Sets the path rewrite function. The default is `DefaultPathRewriteFunc`.
- `WithMatchFunc`: Sets the match function. The default is `DefaultLimitMatchFunc`.
- `WithIpWhitelist`: Sets the IP whitelist. The default is `DefaultIpWhitelist`.

### Methods

- `HandlerFunc`: Returns a `gin.HandlerFunc` for `orbit` or `gin`.
- `Stop`: Stops the rewriter. This is an empty function and does not need to be called.

### Example

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	rw "github.com/shengyanli1982/orbit-contrib/pkg/rewriter"
)

var (
	// 测试URL
	// Test URL
	testUrl = "/test"
)

// testRewriteFunc 是一个测试重写函数
// testRewriteFunc is a test rewrite function
func testRewriteFunc(u *url.URL) (bool, string) {
	// 如果URL的路径等于测试URL
	// If the path of the URL equals the test URL
	if u.Path == testUrl {
		// 返回真和新的URL
		// Return true and the new URL
		return true, testUrl + "2"
	}
	// 否则返回假和空字符串
	// Otherwise return false and an empty string
	return false, ""
}

// testRequestFunc 是一个测试请求函数
// testRequestFunc is a test request function
func testRequestFunc(idx int, router *gin.Engine, conf *rw.Config, url string) {
	// 创建一个新的请求
	// Create a new request
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// 创建一个新的响应记录器
	// Create a new response recorder
	resp := httptest.NewRecorder()

	// 使用路由器处理HTTP请求
	// Use the router to handle the HTTP request
	router.ServeHTTP(resp, req)

	// 打印请求的信息
	// Print the information of the request
	fmt.Println("[Request]", idx, resp.Code, url, req.URL.Path, resp.Body.String())
}

func main() {
	// 创建一个新的配置
	// Create a new Config
	conf := rw.NewConfig().WithPathRewriteFunc(testRewriteFunc)

	// 创建一个新的路径重写器
	// Create a new Path Rewriter
	compr := rw.NewPathRewriter(conf)

	// 在函数返回时停止路径重写器
	// Stop the Path Rewriter when the function returns
	defer compr.Stop()

	// 创建一个新的Gin路由器
	// Create a new Gin router
	router := gin.New()

	// 使用路径重写器的处理函数
	// Use the handler function of the Path Rewriter
	router.Use(compr.HandlerFunc())

	// 添加一个新的路由
	// Add a new route
	router.GET(testUrl, func(c *gin.Context) {
		// 当请求成功时，返回"OK"
		// Return "OK" when the request is successful
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
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 0 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 1 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 2 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 3 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 4 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 5 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 6 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 7 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 8 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 307 with 200
[Request] 9 307 /test /test2 <a href="/test2">Temporary Redirect</a>.

OK
```
