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
