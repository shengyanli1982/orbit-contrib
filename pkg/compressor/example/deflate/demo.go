package main

import (
	"compress/flate"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	cr "github.com/shengyanli1982/orbit-contrib/pkg/compressor"
)

var (
	testUrl = "/test"
)

func testNewDeflateWriterFunc(config *cr.Config, rw gin.ResponseWriter) any {
	return cr.NewDeflateWriter(config, rw)
}

func testRequestFunc(idx int, router *gin.Engine, conf *cr.Config, url string) {
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Get the response body
	bodyContent := resp.Body.String()

	// Create deflate reader
	gr := flate.NewReader(resp.Body)
	defer gr.Close()

	// Read the response
	plaintext, _ := io.ReadAll(gr)

	// Print the request information
	fmt.Println("[Request]", idx, resp.Code, url, string(plaintext), bodyContent)
}

func main() {
	// Create a new rate limiter
	conf := cr.NewConfig().WithWriterCreateFunc(testNewDeflateWriterFunc)
	compr := cr.NewCompressor(conf)
	defer compr.Stop()

	// Create a test context
	router := gin.New()
	router.Use(compr.HandlerFunc())
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
