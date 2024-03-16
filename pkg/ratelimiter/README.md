# Ratelimiter

**Ratelimiter** is a lightweight rate limiter for `gin` and `orbit`. It implements the [token bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm and is designed to protect API endpoints.

`Ratelimiter` supports two modes:

-   **IP-based**: Each IP address has its own rate limit bucket.
-   **Shared**: All requests share the same rate limit bucket.

`Ratelimiter` is built using native Go packages and other powerful libraries:

-   [gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)
-   [xxhash](https://github.com/cespare/xxhash)
-   [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)

It is designed to be safe for concurrent usage.

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter
```

## Quick Start

### Configuration

The `Ratelimiter` component has a configuration object that can be used to customize its behavior. The configuration object provides the following methods:

-   `WithCallback`: Sets the callback function. The default is `&emptyCallback{}`.
-   `WithRate`: Sets the rate. The default is `float64(1)`.
-   `WithBurst`: Sets the burst. The default is `1`.
-   `WithMatchFunc`: Sets the match function. The default is `DefaultLimitMatchFunc`.
-   `WithIpWhitelist`: Sets the IP whitelist. The default is `DefaultIpWhitelist`.

### Components

#### 1. Ratelimiter

The `Ratelimiter` component is the main component used for rate limiting requests.

**Methods**

-   `GetLimiter`: Retrieves the limiter.
-   `SetRate`: Sets the rate for the limiter in a thread-safe manner.
-   `SetBurst`: Sets the burst for the limiter in a thread-safe manner.
-   `HandlerFunc`: Returns a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop`: Stops the limiter. This is an empty function and does not need to be called.

**Example**

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	rl "github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter"
)

var (
	// 测试URL
	// Test URL
	testUrl = "/test"
)

// testRequestFunc 是一个测试请求的函数
// testRequestFunc is a function to test the request
func testRequestFunc(idx int, router *gin.Engine, conf *rl.Config, url string) {
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
	fmt.Println("[Request]", idx, resp.Code, url)
}

func main() {
	// 创建一个新的配置，设置速率为2，突发数为5
	// Create a new configuration, set the rate to 2 and the burst to 5
	conf := rl.NewConfig().WithRate(2).WithBurst(5)

	// 创建一个新的速率限制器
	// Create a new rate limiter
	limiter := rl.NewRateLimiter(conf)

	// 在函数返回时停止速率限制器
	// Stop the rate limiter when the function returns
	defer limiter.Stop()

	// 创建一个新的路由器
	// Create a new router
	router := gin.New()

	// 使用速率限制器的处理函数
	// Use the handler function of the rate limiter
	router.Use(limiter.HandlerFunc())

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
[Request] 0 200 /test
[Request] 1 200 /test
[Request] 2 200 /test
[Request] 3 200 /test
[Request] 4 200 /test
[Request] 5 429 /test
[Request] 6 429 /test
[Request] 7 429 /test
[Request] 8 429 /test
[Request] 9 429 /test
```

#### 2. Ip Ratelimiter

**Methods**

-   `GetLimiter`: Retrieves the limiter by key.
-   `SetRate`: Sets the rate for the limiter in a thread-safe manner.
-   `SetBurst`: Sets the burst for the limiter in a thread-safe manner.
-   `HandlerFunc`: Returns a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop`: Stops the limiter and releases the associated resources.

**Example**

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	rl "github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter"
)

var (
	// 测试URL
	// Test URL
	testUrl = "/test"

	// 测试端口
	// Test port
	testPort = 13143

	// 测试IP地址1
	// Test IP address 1
	testIpAddress = "192.168.0.11"

	// 测试端点1
	// Test endpoint 1
	testEndpoint = fmt.Sprintf("%s:%d", testIpAddress, testPort)

	// 测试IP地址2
	// Test IP address 2
	testIpAddress2 = "192.168.0.12"

	// 测试端点2
	// Test endpoint 2
	testEndpoint2 = fmt.Sprintf("%s:%d", testIpAddress2, testPort)
)

// testRequestFunc 是一个测试请求的函数
// testRequestFunc is a function to test the request
func testRequestFunc(idx int, router *gin.Engine, conf *rl.Config, ep, url string) {
	// 创建一个新的请求
	// Create a new request
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// 设置请求的远程地址
	// Set the remote address of the request
	req.RemoteAddr = ep

	// 创建一个新的响应记录器
	// Create a new response recorder
	resp := httptest.NewRecorder()

	// 使用路由器处理HTTP请求
	// Use the router to handle the HTTP request
	router.ServeHTTP(resp, req)

	// 打印请求的信息
	// Print the information of the request
	fmt.Println("[Request]", idx, resp.Code, ep, url)
}

func main() {
	// 创建一个新的配置，设置速率为2，突发数为5
	// Create a new configuration, set the rate to 2 and the burst to 5
	conf := rl.NewConfig().WithRate(2).WithBurst(5)

	// 创建一个新的IP速率限制器
	// Create a new IP rate limiter
	limiter := rl.NewIpRateLimiter(conf)

	// 在函数返回时停止速率限制器
	// Stop the rate limiter when the function returns
	defer limiter.Stop()

	// 创建一个新的路由器
	// Create a new router
	router := gin.New()

	// 使用速率限制器的处理函数
	// Use the handler function of the rate limiter
	router.Use(limiter.HandlerFunc())

	// 添加一个GET路由
	// Add a GET route
	router.GET(testUrl, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 测试 10 次请求
	// Test 10 requests
	for i := 0; i < 10; i++ {
		// 调用测试请求函数，传入索引、路由器、配置、端点和测试URL
		// Call the test request function, passing in the index, router, configuration, endpoint, and test URL
		testRequestFunc(i, router, conf, testEndpoint, testUrl)
	}

	// 测试 10 次请求
	// Test 10 requests
	for i := 0; i < 10; i++ {
		// 调用测试请求函数，传入索引、路由器、配置、端点和测试URL
		// Call the test request function, passing in the index, router, configuration, endpoint, and test URL
		testRequestFunc(i, router, conf, testEndpoint2, testUrl)
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
[Request] 0 200 192.168.0.11:13143 /test
[Request] 1 200 192.168.0.11:13143 /test
[Request] 2 200 192.168.0.11:13143 /test
[Request] 3 200 192.168.0.11:13143 /test
[Request] 4 200 192.168.0.11:13143 /test
[Request] 5 429 192.168.0.11:13143 /test
[Request] 6 429 192.168.0.11:13143 /test
[Request] 7 429 192.168.0.11:13143 /test
[Request] 8 429 192.168.0.11:13143 /test
[Request] 9 429 192.168.0.11:13143 /test
[Request] 0 200 192.168.0.12:13143 /test
[Request] 1 200 192.168.0.12:13143 /test
[Request] 2 200 192.168.0.12:13143 /test
[Request] 3 200 192.168.0.12:13143 /test
[Request] 4 200 192.168.0.12:13143 /test
[Request] 5 429 192.168.0.12:13143 /test
[Request] 6 429 192.168.0.12:13143 /test
[Request] 7 429 192.168.0.12:13143 /test
[Request] 8 429 192.168.0.12:13143 /test
[Request] 9 429 192.168.0.12:13143 /test
```
