package ratelimiter

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

var (
	testIpAddress2 = "192.168.0.12"
	testEndpoint2  = fmt.Sprintf("%s:%d", testIpAddress2, testPort)
)

func TestIpRateLimiter_GetLimiter(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	rl := NewIpRateLimiter(conf)

	// Test case 1: Limiter exists in cache
	key := testIpAddress
	limiter := rate.NewLimiter(rate.Limit(10), 100)
	rl.cache.Set(key, &RateLimiter{config: conf, limiter: limiter})

	result := rl.GetLimiter(key)
	assert.NotNil(t, result)
	assert.Equal(t, limiter, result)

	// Test case 2: Limiter does not exist in cache
	key = testIpAddress2
	result = rl.GetLimiter(key)
	assert.Nil(t, result)
}

func TestIpRateLimiter_SetRate(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	rl := NewIpRateLimiter(conf)

	// Set up test data
	key1 := testIpAddress
	limiter := rate.NewLimiter(rate.Limit(10), 100)
	rl.cache.Set(key1, &RateLimiter{config: conf, limiter: limiter})

	// Set rate for all limiters
	rate := float64(10)
	rl.SetRate(rate)

	// Check if rate is set correctly for all limiters
	result := limiter.Limit()
	assert.Equal(t, rate, float64(result))
}

func TestIpRateLimiter_SetBurst(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	rl := NewIpRateLimiter(conf)

	// Set up test data
	key1 := testIpAddress
	limiter := rate.NewLimiter(rate.Limit(10), 100)
	rl.cache.Set(key1, &RateLimiter{config: conf, limiter: limiter})

	// Set rate for all limiters
	burst := 10
	rl.SetBurst(burst)

	// Check if rate is set correctly for all limiters
	result := limiter.Burst()
	assert.Equal(t, burst, int(result))
}

func TestIpRateLimiter_RateAndBurst(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithRate(2).WithBurst(5)
	rl := NewIpRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(rl.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter - 1
	// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testRequestFunc(t, i, router, conf, &wg, testEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}

	// Test the rate limiter - 2
	// Send multiple requests to test the rate limiter
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testRequestFunc(t, i, router, conf, &wg, testEndpoint2, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}

func TestIpRateLimiter_Callback(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithCallback(&testCallback{t: t})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testRequestFunc(t, i, router, conf, &wg, testEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}

func TestIpRateLimiter_MatchFunc(t *testing.T) {
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

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
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

func TestIpRateLimiter_Whitelist(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testWhitelistRequestFunc(t, i, router, &wg, defaultEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}

func TestIpRateLimiter_CustomWhitelist(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithWhitelist([]string{testIpAddress})
	limiter := NewRateLimiter(conf)

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go testWhitelistRequestFunc(t, i, router, &wg, testEndpoint, testUrlPath)
		// Wait for all goroutines to finish
		wg.Wait()
	}
}
