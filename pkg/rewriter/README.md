# Rewriter

**Rewriter** is a simple middleware for rewriting request url path and return status code `http.StatusTemporaryRedirect`. It can be used in `gin` and `orbit`. And designed to be used in redirect request which is need to be rewrite.

Why use `http.StatusTemporaryRedirect` instead of `http.StatusMovedPermanently`? Because `http.StatusTemporaryRedirect` will not be cached by the browser, and the browser will not send the request to the original url again.

`Rewriter` based golang native package and other powerful packages:

-   [gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/rewriter
```

## Quick Start

### Config

`Rewriter` has a config object, which can be used to configure the batch process behavior. The config object can be used following methods to set.

-   `WithCallback` : set the callback function. The default is `&emptyCallback{}`.
-   `WithPathRewriteFunc` : set the path rewrite function. The default is `DefaultPathRewriteFunc`.
-   `WithMatchFunc` : set the match function. The default is `DefaultLimitMatchFunc`.
-   `WithIpWhitelist` : set the whitelist. The default is `DefaultIpWhitelist`.

### Methods

-   `HandlerFunc` : return a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop` : stop the rewriter. It is empty function, no need to call it.

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
	testUrl = "/test"
)

// testRewriteFunc is a test rewrite function
func testRewriteFunc(u *url.URL) (bool, string) {
	if u.Path == testUrl {
		return true, testUrl + "2"
	}
	return false, ""
}

// testRequestFunc is a test request function
func testRequestFunc(idx int, router *gin.Engine, conf *rw.Config, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the request information
	fmt.Println("[Request]", idx, resp.Code, url, req.URL.Path, resp.Body.String())
}

func main() {
	// Create a new Config
	conf := rw.NewConfig().WithPathRewriteFunc(testRewriteFunc)

	// Create a new Compressor
	compr := rw.NewPathRewriter(conf)
	defer compr.Stop()

	// Create a new Gin router
	router := gin.New()
	router.Use(compr.HandlerFunc())

	// Add a new route
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
