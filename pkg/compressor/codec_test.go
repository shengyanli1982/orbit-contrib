package compressor

// func TestGZipWriter_WriteHeader(t *testing.T) {
// 	rw := gin.CreateTestContext(w http.ResponseWriter)

// 	// Create a new Config
// 	conf := NewConfig().WithThreshold(10)

// 	// Create a new GZipWriter
// 	gw := NewGZipWriter(conf)
// 	defer gw.Stop()

// 	// Set the underlying ResponseWriter
// 	rec := httptest.NewRecorder()
// 	err := gw.Reset(rec)

// 	// Check if there is no error
// 	assert.Nil(t, err)

// 	// Call the WriteHeader method
// 	gw.WriteHeader(http.StatusOK)

// 	// Check if the underlying ResponseWriter has the same status code
// 	assert.Equal(t, http.StatusOK, rec.Code)
// }

// func TestGZipWriter_Write(t *testing.T) {
// 	// Create a new Config
// 	conf := NewConfig().WithThreshold(10)

// 	// Create a new GZipWriter
// 	gw := NewGZipWriter(conf)
// 	defer gw.Stop()

// 	// Set the underlying ResponseWriter
// 	rec := httptest.NewRecorder()
// 	err := gw.Reset(rec)

// 	// Check if there is no error
// 	assert.Nil(t, err)

// 	// Set the input message
// 	msg := []byte("Hello, World!")

// 	// Call the Write method
// 	n, err := gw.Write(msg)

// 	// Check if there is no error
// 	assert.Nil(t, err)

// 	// Check if the number of bytes written is correct
// 	assert.Equal(t, len(msg), n)

// 	// Check if the underlying ResponseWriter has the same content
// 	assert.Equal(t, string(msg), rec.Body.String())
// }

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	testHttpStatusCode = 444
)

func TestGZipWriter_WriteHeader(t *testing.T) {
	// Create a new Gin router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Save the original writer
		origWriter := c.Writer

		// Create a new Config
		conf := NewConfig().WithThreshold(10)

		// Create a new GZipWriter
		gw := NewGZipWriter(conf, c.Writer)
		defer gw.Stop()

		// Set the underlying ResponseWriter
		c.Writer = gw

		// Call the Next method
		c.Next()

		// Call the WriteHeader method
		gw.WriteHeader(testHttpStatusCode)

		// Restore the original writer
		c.Writer = origWriter
	})

	// Add a new route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, testHttpStatusCode, w.Code)
	assert.Equal(t, "\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\xf2\xf7\x06\x04\x00\x00\xff\xff-\xd96\xd7\x02\x00\x00\x00", w.Body.String())
}
