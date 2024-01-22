package ratelimiter

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestIpRateLimiter_GetLimiter(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	rl := NewIpRateLimiter(conf)
	defer rl.Stop()

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
	defer rl.Stop()

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
	defer rl.Stop()

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
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, testEndpoint, testUrlPath)
	}

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, testEndpoint2, testUrlPath)
	}
}

func TestIpRateLimiter_Callback(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithCallback(&testCallback{t: t})
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, testEndpoint, testUrlPath)
	}
}

func TestIpRateLimiter_MatchFunc(t *testing.T) {
	path := testUrlPath + "2"

	// Create a new rate limiter
	conf := NewConfig().WithMatchFunc(func(header *http.Request) bool {
		return header.URL.Path == testUrlPath
	})
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

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
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, testEndpoint, testUrlPath)
	}

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testWhitelistRequestFunc(t, i, router, testEndpoint, path)
	}
}

func TestIpRateLimiter_Whitelist(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testWhitelistRequestFunc(t, i, router, defaultEndpoint, testUrlPath)
	}
}

func TestIpRateLimiter_CustomWhitelist(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithIpWhitelist([]string{testIpAddress})
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

	// Create a test context
	router := gin.New()
	router.Use(limiter.HandlerFunc())
	router.GET(testUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testWhitelistRequestFunc(t, i, router, testEndpoint, testUrlPath)
	}
}
