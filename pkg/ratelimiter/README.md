# Ratelimiter

**Ratelimiter** is a simple rate limiter for `gin` and `orbit`. It is based on [token bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm. It is designed to be used in a protected API endpoint.

`Ratelimiter` has two modes:

-   **Base on IP**: Each IP address has its own bucket
-   **All share one**: All requests share the same bucket

`Ratelimiter` based golang native package and other powerful packages:

-   [gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)
-   [xxhash](https://github.com/cespare/xxhash)
-   [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)

So it is safe to use in a concurrent environment.

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter
```

## Quick Start

### Config

`Ratelimiter` has a config object, which can be used to configure the batch process behavior. The config object can be used following methods to set.

-   `WithCallback` : set the callback function. The default is `&emptyCallback{}`.
-   `WithRate` : set the rate. The default is `float64(1)`.
-   `WithBurst` : set the burst. The default is `1`.
-   `WithMatchFunc` : set the match function. The default is `DefaultLimitMatchFunc`.
-   `WithIpWhitelist` : set the whitelist. The default is `DefaultIpWhitelist`.

### Components

#### 1. Ratelimiter

`Ratelimiter` is used the most component. It is used to limit the rate of requests.

**Methods**

-   `GetLimiter` : get the limiter.
-   `SetRate` : set the rate for limiter, thread safe.
-   `SetBurst` : set the burst for limiter, thread safe.
-   `HandlerFunc` : return a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop` : stop the limiter. It is empty function, no need to call it.

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
	testUrl = "/test"
)

func testRequestFunc(idx int, router *gin.Engine, conf *rl.Config, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the request information
	fmt.Println("[Request]", idx, resp.Code, url)
}

func main() {
	// Create a new rate limiter
	conf := rl.NewConfig().WithRate(2).WithBurst(5)
	limiter := rl.NewRateLimiter(conf)
	defer limiter.Stop()

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
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

-   `GetLimiter` : get the limiter by key.
-   `SetRate` : set the rate for limiter, thread safe.
-   `SetBurst` : set the burst for limiter, thread safe.
-   `HandlerFunc` : return a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop` : stop the limiter. It is used to release the resources.

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
	testUrl        = "/test"
	testPort       = 13143
	testIpAddress  = "192.168.0.11"
	testEndpoint   = fmt.Sprintf("%s:%d", testIpAddress, testPort)
	testIpAddress2 = "192.168.0.12"
	testEndpoint2  = fmt.Sprintf("%s:%d", testIpAddress2, testPort)
)

func testRequestFunc(idx int, router *gin.Engine, conf *rl.Config, ep, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.RemoteAddr = ep
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the request information
	fmt.Println("[Request]", idx, resp.Code, ep, url)
}

func main() {
	// Create a new rate limiter
	conf := rl.NewConfig().WithRate(2).WithBurst(5)
	limiter := rl.NewIpRateLimiter(conf)
	defer limiter.Stop()

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrl, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		testRequestFunc(i, router, conf, testEndpoint, testUrl)
	}

	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		testRequestFunc(i, router, conf, testEndpoint2, testUrl)
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
