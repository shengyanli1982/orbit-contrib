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
	testResponseText     = "OK"
	testGZipResponseText = "\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\xf2\xf7\x06\x04\x00\x00\xff\xff-\xd96\xd7\x02\x00\x00\x00"
)

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

		// Call the Next method
		c.Next()

		// Stop the GZipWriter
		gw.Stop()
	})

	// Add a new route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, testResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testGZipResponseText, w.Body.String())
}

func TestGZipWriter_Reset(t *testing.T) {
	// Create a new Gin router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Create a new Config
		conf := NewConfig()

		// Create a new GZipWriter
		testCtx := gin.CreateTestContextOnly(httptest.NewRecorder(), router)
		gw := NewGZipWriter(conf, testCtx.Writer)

		// Reset the GZipWriter with the underlying ResponseWriter
		err := gw.Reset(c.Writer)

		// Check if there is no error
		assert.Nil(t, err)

		// Set the underlying ResponseWriter
		c.Writer = gw

		// Call the Next method
		c.Next()

		// Stop the GZipWriter
		gw.Stop()
	})

	// Add a new route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, testResponseText)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testGZipResponseText, w.Body.String())
}
