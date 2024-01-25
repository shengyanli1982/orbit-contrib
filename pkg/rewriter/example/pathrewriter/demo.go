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
