package ratelimiter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
	"github.com/stretchr/testify/assert"
)

type testCallback struct {
	t *testing.T
}

func (c *testCallback) OnLimited(header *http.Request) {
	assert.Equal(c.t, com.TestEndpoint, header.RemoteAddr)
}

func testRequestFunc(t *testing.T, idx int, router *gin.Engine, conf *Config, ep, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.RemoteAddr = ep
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert that the response status code is 200 OK for the first "conf.burst" requests
	if idx < conf.burst {
		assert.Equal(t, http.StatusOK, resp.Code)
	} else {
		// Assert that the response status code is 429 Too Many Requests for the remaining requests
		assert.Equal(t, http.StatusTooManyRequests, resp.Code)
	}

	// Print the request information
	fmt.Println("[Request]", idx, ep, resp.Code, url)
}

func testWhitelistRequestFunc(t *testing.T, idx int, router *gin.Engine, ep, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.RemoteAddr = ep
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert that the response status code is 200 OK for the first "conf.burst" requests
	assert.Equal(t, http.StatusOK, resp.Code)

	// Print the request information
	fmt.Println("[Request]", idx, ep, resp.Code, url)
}

func TestLimiter_RateAndBurst(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithRate(2).WithBurst(5)
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		testRequestFunc(t, i, router, conf, com.TestEndpoint, com.TestUrlPath)
	}
}

func TestLimiter_Callback(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithCallback(&testCallback{t: t})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, com.TestEndpoint, com.TestUrlPath)
	}
}

func TestLimiter_MatchFunc(t *testing.T) {
	path := com.TestUrlPath + "2"

	// Create a new rate limiter
	conf := NewConfig().WithMatchFunc(func(header *http.Request) bool {
		return header.URL.Path == com.TestUrlPath
	})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET(path, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, com.TestEndpoint, com.TestUrlPath)
	}

	// Test the rate limiter/
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testWhitelistRequestFunc(t, i, router, com.TestEndpoint, path)
	}
}

func TestLimiter_IpWhitelistWithLocal(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithIpWhitelist([]string{com.DefaultLocalIpAddress})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testWhitelistRequestFunc(t, i, router, com.DefaultEndpoint, com.TestUrlPath)
	}
}

func TestLimiter_IpWhitelistWithOther(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithIpWhitelist([]string{com.TestIpAddress3})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testWhitelistRequestFunc(t, i, router, com.TestEndpoint3, com.TestUrlPath)
	}
}
