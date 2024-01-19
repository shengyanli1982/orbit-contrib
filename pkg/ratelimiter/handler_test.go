package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	testIpAddress = "192.168.0.1"
)

func TestNewRateLimiterHandlerFunc(t *testing.T) {
	// Create a test config
	config := NewConfig().WithRate(1)

	// Create a new rate limiter handler function
	handler := NewRateLimiterHandlerFunc(config)

	// Test when the request matches the config
	config.match = func(req *http.Request) bool {
		assert.Equal(t, req.URL.Path, "/test")
		assert.Equal(t, req.Method, "GET")
		return true
	}

	// Create a new Gin router
	router := gin.New()

	// Use the BodyBuffer middleware
	router.Use(handler)

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test response")
	})

	// Define a test function
	testFunc := func(addr string, code int) {
		// Create a test request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = addr + ":1234"

		// Create a test response recorder
		rec := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(rec, req)

		// Assert the response
		assert.Equal(t, code, rec.Code)
		assert.Equal(t, "Test response", rec.Body.String())
	}

	// Test the handler ( should return 200 ), low speed
	testFunc(testIpAddress, http.StatusOK)

	// Test the handler ( should return 200 ), default whitelist
	for i := 0; i < 10; i++ {
		go testFunc(DefaultLocalIpAddress, http.StatusTooManyRequests)
	}

	// Test the handler ( should return 429 ), high speed
	for i := 0; i < 10; i++ {
		go testFunc(testIpAddress, http.StatusTooManyRequests)
	}
}
