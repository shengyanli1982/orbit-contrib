package ratelimiter

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestIpRateLimiter_GetLimiter(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	rl := NewIpRateLimiter(conf)
	defer rl.Stop()

	// Test case 1: Limiter exists in cache
	key := com.TestIpAddress
	limiter := rate.NewLimiter(rate.Limit(10), 100)
	rl.cache.Set(key, &RateLimiter{config: conf, limiter: limiter})

	result := rl.GetLimiter(key)
	assert.NotNil(t, result)
	assert.Equal(t, limiter, result)

	// Test case 2: Limiter does not exist in cache
	key = com.TestIpAddress2
	result = rl.GetLimiter(key)
	assert.Nil(t, result)
}

func TestIpRateLimiter_SetRate(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig()
	rl := NewIpRateLimiter(conf)
	defer rl.Stop()

	// Set up test data
	key1 := com.TestIpAddress
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
	key1 := com.TestIpAddress
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
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, com.TestEndpoint, com.TestUrlPath)
	}

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testRequestFunc(t, i, router, conf, com.TestEndpoint2, com.TestUrlPath)
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

func TestIpRateLimiter_MatchFunc(t *testing.T) {
	path := com.TestUrlPath + "2"

	// Create a new rate limiter
	conf := NewConfig().WithMatchFunc(func(header *http.Request) bool {
		return header.URL.Path == com.TestUrlPath
	})
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

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

	// Test the rate limiter
	// Send multiple requests to test the rate limiter
	for i := 0; i < 10; i++ {
		// Start a new test the rate limiter
		testWhitelistRequestFunc(t, i, router, com.TestEndpoint, path)
	}
}

func TestIpRateLimiter_IpWhitelistWithLocal(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithIpWhitelist([]string{com.DefaultLocalIpAddress})
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

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

func TestIpRateLimiter_IpWhitelistWithOther(t *testing.T) {
	// Create a new rate limiter
	conf := NewConfig().WithIpWhitelist([]string{com.TestIpAddress})
	limiter := NewIpRateLimiter(conf)
	defer limiter.Stop()

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
		testWhitelistRequestFunc(t, i, router, com.TestEndpoint, com.TestUrlPath)
	}
}
