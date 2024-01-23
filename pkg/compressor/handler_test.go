package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestCompressorHandlerFunc_Default(t *testing.T) {
	// Create a new Config
	conf := NewConfig()

	// Create a new Compressor
	compr := NewCompressor(conf)
	defer compr.Stop()

	// Create a new Gin router
	router := gin.New()
	router.Use(compr.HandlerFunc())

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, testResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Encoding"), "gzip")
	assert.Equal(t, w.Header().Get("Vary"), "Accept-Encoding")

	// Create gzip reader
	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gr.Close()

	// Read the response
	plaintext, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, string(plaintext), testResponseText)
}

func TestCompressorHandlerFunc_MatchFunc(t *testing.T) {
	// Create a new Config
	conf := NewConfig().WithMatchFunc(func(header *http.Request) bool {
		return header.URL.Path == com.TestUrlPath
	})

	// Create a new Compressor
	compr := NewCompressor(conf)
	defer compr.Stop()

	// Create a new Gin router
	router := gin.New()
	router.Use(compr.HandlerFunc())

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, testResponseText)
	})
	router.GET(com.TestUrlPath2, func(c *gin.Context) {
		c.String(http.StatusOK, testResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath2, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testResponseText, w.Body.String())

	// Create a new recorder
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Encoding"), "gzip")
	assert.Equal(t, w.Header().Get("Vary"), "Accept-Encoding")

	// Create gzip reader
	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gr.Close()

	// Read the response
	plaintext, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, string(plaintext), testResponseText)
}

func TestCompressorHandlerFunc_IpWhitelist(t *testing.T) {
	// Create a new Config
	conf := NewConfig()

	// Create a new Compressor
	compr := NewCompressor(conf)
	defer compr.Stop()

	// Create a new Gin router
	router := gin.New()
	router.Use(compr.HandlerFunc())

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, testResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)
	req.RemoteAddr = com.DefaultEndpoint

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testResponseText, w.Body.String())
}

func TestCompressorHandlerFunc_CustomWhitelist(t *testing.T) {
	// Create a new Config
	conf := NewConfig().WithIpWhitelist([]string{com.TestIpAddress})

	// Create a new Compressor
	compr := NewCompressor(conf)
	defer compr.Stop()

	// Create a new Gin router
	router := gin.New()
	router.Use(compr.HandlerFunc())

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, testResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)
	req.RemoteAddr = com.TestEndpoint

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testResponseText, w.Body.String())
}
