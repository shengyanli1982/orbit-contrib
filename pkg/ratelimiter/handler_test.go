package ratelimiter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	testIpAddress = "192.168.0.1"
)

func TestNewRateLimiterHandlerFunc_Request(t *testing.T) {
	// Create a test config
	config := NewConfig().WithRate(1)

	// Test when the request matches the config
	config.match = func(req *http.Request) bool {
		assert.Equal(t, req.URL.Path, "/test")
		assert.Equal(t, req.Method, "GET")
		return true
	}

	// Create a new rate limiter handler function
	handler := NewRateLimiterHandlerFunc(config)

	// Create a new Gin router
	router := gin.New()

	// Use the BodyBuffer middleware
	router.Use(handler)

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})

	//Define a stress function
	stressFunc := func(addr string) {
		// Create a test request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = addr + ":1234"

		// Create a Test recorder
		rec := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(rec, req)
	}

	// Define a test function
	testFunc := func(addr string, code int) {
		// Create a test request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = addr + ":1234"

		// Create a Test recorder
		rec := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(rec, req)

		// Assert the response
		// If the response code is not 200, it means the request is limited
		assert.Equal(t, code, rec.Code)
	}

	t.Run("LowSpeedRequest", func(t *testing.T) {
		// Test the handler ( should return 200 ), low speed
		testFunc(testIpAddress, http.StatusOK)
	})

	// continue to run stressFunc
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				stressFunc(testIpAddress)
			}
		}
	}()

	t.Run("WhitelistRequest", func(t *testing.T) {
		// Test the handler ( should return 200 ), default whitelist
		for i := 0; i < 10; i++ {
			go testFunc(DefaultLocalIpAddress, http.StatusOK)
		}
	})

	t.Run("HighSpeedRequest", func(t *testing.T) {
		// Test the handler ( should return 429 ), high speed
		for i := 0; i < 10; i++ {
			go testFunc(testIpAddress, http.StatusTooManyRequests)
		}
	})

	// Wait for the goroutine to finish
	time.Sleep(time.Second)
}

func TestNewRateLimiterHandlerFunc_MatchFunc(t *testing.T) {
	// Create a test config
	config := NewConfig().WithRate(1).WithMatchFunc(func(req *http.Request) bool {
		return req.URL.Path == "/test"
	})

	// Create a new rate limiter handler function
	handler := NewRateLimiterHandlerFunc(config)

	// Create a new Gin router
	router := gin.New()

	// Use the BodyBuffer middleware
	router.Use(handler)

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})
	router.GET("/test2", func(c *gin.Context) {
		c.String(http.StatusOK, "Test2")
	})

	//Define a stress function
	stressFunc := func(addr string) {
		// Create a test request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = addr + ":1234"

		// Create a Test recorder
		rec := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(rec, req)
	}

	// Define a test function
	testFunc := func(addr string, code int, path string) {
		// Create a test request
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		req.RemoteAddr = addr + ":1234"

		// Create a Test recorder
		rec := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(rec, req)

		// Assert the response
		// If the response code is not 200, it means the request is limited
		assert.Equal(t, code, rec.Code)
	}

	// continue to run stressFunc
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				stressFunc(testIpAddress)
			}
		}
	}()

	t.Run("MatchPath", func(t *testing.T) {
		// Test the handler ( should return 200 )
		for i := 0; i < 10; i++ {
			go testFunc(testIpAddress, http.StatusTooManyRequests, "/test")
		}
	})

	t.Run("DontMatch", func(t *testing.T) {
		// Test the handler ( should return 429 ), high speed
		for i := 0; i < 10; i++ {
			go testFunc(testIpAddress, http.StatusOK, "/test2")
		}
	})

	// Wait for the goroutine to finish
	time.Sleep(time.Second)
}
