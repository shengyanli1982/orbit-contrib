package compressor

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
	"github.com/stretchr/testify/assert"
)

func testNewDeflateWriterFunc(config *Config, rw gin.ResponseWriter) any {
	return NewDeflateWriter(config, rw)
}

func TestGZipWriter_Write(t *testing.T) {
	// Create a new Gin router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Create a new Config
		conf := NewConfig()

		// Create a new GZipWriter
		gw := NewGZipWriter(conf, c.Writer)

		// Set the underlying ResponseWriter
		c.Writer = gw

		// Set the Content-Encoding and Vary headers
		c.Header("Content-Encoding", GZipContentEncoding)
		c.Header("Vary", "Accept-Encoding")

		// Call the Next method
		c.Next()

		// Set the Content-Length header
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))

		// Stop the GZipWriter
		gw.Stop()
	})

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, com.TestResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Encoding"), GZipContentEncoding)
	assert.Equal(t, w.Header().Get("Vary"), "Accept-Encoding")

	// Create gzip reader
	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gr.Close()

	// Read the response
	plaintext, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, string(plaintext), com.TestResponseText)
}

func TestGZipWriter_Reset(t *testing.T) {
	// Create a new Gin router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Create a new Config
		conf := NewConfig()

		// Create a new GZipWriter
		gw := NewGZipWriter(conf, nil)

		// Reset the underlying ResponseWriter
		err := gw.ResetCompressWriter(c.Writer)
		assert.NoError(t, err)
		err = gw.ResetResponseWriter(c.Writer)
		assert.NoError(t, err)

		// Set the underlying ResponseWriter
		c.Writer = gw

		// Set the Content-Encoding and Vary headers
		c.Header("Content-Encoding", GZipContentEncoding)
		c.Header("Vary", "Accept-Encoding")

		// Call the Next method
		c.Next()

		// Set the Content-Length header
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))

		// Stop the GZipWriter
		gw.Stop()
	})

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, com.TestResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Encoding"), GZipContentEncoding)
	assert.Equal(t, w.Header().Get("Vary"), "Accept-Encoding")

	// Create gzip reader
	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gr.Close()

	// Read the response
	plaintext, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, string(plaintext), com.TestResponseText)
}

func TestDeflateWriter_Write(t *testing.T) {
	// Create a new Gin router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Create a new Config
		conf := NewConfig().WithWriterCreateFunc(testNewDeflateWriterFunc)

		// Create a new DeflateWriter
		dw := NewDeflateWriter(conf, c.Writer)

		// Set the underlying ResponseWriter
		c.Writer = dw

		// Set the Content-Encoding and Vary headers
		c.Header("Content-Encoding", DeflateContentEncoding)
		c.Header("Vary", "Accept-Encoding")

		// Call the Next method
		c.Next()

		// Set the Content-Length header
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))

		// Stop the DeflateWriter
		dw.Stop()
	})

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, com.TestResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Encoding"), DeflateContentEncoding)
	assert.Equal(t, w.Header().Get("Vary"), "Accept-Encoding")

	// Create deflate reader
	dr := flate.NewReader(w.Body)
	defer dr.Close()

	// Read the response
	plaintext, err := io.ReadAll(dr)
	assert.NoError(t, err)
	assert.Equal(t, string(plaintext), com.TestResponseText)
}

func TestDeflateWriter_Reset(t *testing.T) {
	// Create a new Gin router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Create a new Config
		conf := NewConfig().WithWriterCreateFunc(testNewDeflateWriterFunc)

		// Create a new DeflateWriter
		dw := NewDeflateWriter(conf, nil)

		// Reset the underlying ResponseWriter
		err := dw.ResetCompressWriter(c.Writer)
		assert.NoError(t, err)
		err = dw.ResetResponseWriter(c.Writer)
		assert.NoError(t, err)

		// Set the underlying ResponseWriter
		c.Writer = dw

		// Set the Content-Encoding and Vary headers
		c.Header("Content-Encoding", DeflateContentEncoding)
		c.Header("Vary", "Accept-Encoding")

		// Call the Next method
		c.Next()

		// Set the Content-Length header
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))

		// Stop the DeflateWriter
		dw.Stop()
	})

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, com.TestResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Encoding"), DeflateContentEncoding)
	assert.Equal(t, w.Header().Get("Vary"), "Accept-Encoding")

	// Create deflate reader
	dr := flate.NewReader(w.Body)
	defer dr.Close()

	// Read the response
	plaintext, err := io.ReadAll(dr)
	assert.NoError(t, err)
	assert.Equal(t, string(plaintext), com.TestResponseText)
}
