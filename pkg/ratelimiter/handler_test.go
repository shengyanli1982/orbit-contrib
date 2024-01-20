package ratelimiter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	testIpAddress   = "192.168.0.11"
	testPort        = 13143
	testUrlPath     = "/test"
	testEndpoint    = fmt.Sprintf("%s:%d", testIpAddress, testPort)
	defaultEndpoint = fmt.Sprintf("%s:%d", DefaultLocalIpAddress, testPort)
)

type testCallback struct {
	t *testing.T
}

func (c *testCallback) OnLimited(header *http.Request) {
	assert.Equal(c.t, testEndpoint, header.RemoteAddr)
}

func testRequestFunc(t *testing.T, idx int, router *gin.Engine, conf *Config, wg *sync.WaitGroup, ep, url string) {
	// Defer the wait group to ensure that the goroutine is finished
	defer wg.Done()

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

func testWhitelistRequestFunc(t *testing.T, idx int, router *gin.Engine, wg *sync.WaitGroup, ep, url string) {
	// Defer the wait group to ensure that the goroutine is finished
	defer wg.Done()

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

func TestLimiterHandlerFunc_RateAndBurst(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithRate(2).WithBurst(5)
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testRequestFunc(t, i, router, conf, &wg, testEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}

func TestLimiterHandlerFunc_Callback(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithRate(2).WithBurst(5)
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testRequestFunc(t, i, router, conf, &wg, testEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}

func TestLimiterHandlerFunc_MatchFunc(t *testing.T) {
	path := testUrlPath + "2"

	// Create a new rate limiter
	conf := NewConfig().WithMatchFunc(func(header *http.Request) bool {
		return header.URL.Path == testUrlPath
	})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET(path, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testRequestFunc(t, i, router, conf, &wg, testEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}

	// Test the rate limiter// Send multiple requests to test the rate limiter
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testWhitelistRequestFunc(t, i, router, &wg, testEndpoint, path)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}

func TestLimiterHandlerFunc_Whitelist(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testWhitelistRequestFunc(t, i, router, &wg, defaultEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}

func TestLimiterHandlerFunc_CustomWhitelist(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithWhitelist([]string{testIpAddress})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testWhitelistRequestFunc(t, i, router, &wg, testEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}
