package rewriter

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
	"github.com/stretchr/testify/assert"
)

var (
	newContext          = com.TestResponseText + "2"
	testRedirectContext = "<a href=\"/test2\">Temporary Redirect</a>.\n\n"
)

func TestPathRewriter_PathRewrite(t *testing.T) {
	// Create a new Config
	conf := NewConfig().WithPathRewriteFunc(func(u *url.URL) (bool, string) {
		if u.Path == com.TestUrlPath {
			return true, com.TestUrlPath2
		}
		return false, ""
	})

	// Create a new Compressor
	compr := NewPathRewriter(conf)
	defer compr.Stop()

	// Create a new Gin router
	router := gin.New()
	router.Use(compr.HandlerFunc())

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, com.TestResponseText)
	})

	router.GET(com.TestUrlPath2, func(c *gin.Context) {
		c.String(http.StatusOK, newContext)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, com.TestUrlPath2, req.URL.Path)
	assert.Equal(t, testRedirectContext+com.TestResponseText, w.Body.String())

	path := req.URL.Path

	// Create a new recorder
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, path, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	// assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, com.TestUrlPath2, req.URL.Path)
	assert.Equal(t, newContext, w.Body.String())
}

func TestPathRewriter_MatchFunc(t *testing.T) {
	// Create a new Config
	conf := NewConfig().WithMatchFunc(func(r *http.Request) bool {
		return r.URL.Path == com.TestUrlPath && r.Method == http.MethodGet
	}).WithPathRewriteFunc(func(u *url.URL) (bool, string) {
		return true, com.TestUrlPath2
	})

	// Create a new Compressor
	compr := NewPathRewriter(conf)
	defer compr.Stop()

	// Create a new Gin router
	router := gin.New()
	router.Use(compr.HandlerFunc())

	// Add a new route
	router.GET(com.TestUrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, com.TestResponseText)
	})

	router.GET(com.TestUrlPath2, func(c *gin.Context) {
		c.String(http.StatusOK, newContext)
	})

	// Create a new recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, com.TestUrlPath, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, com.TestUrlPath2, req.URL.Path)
	assert.Equal(t, testRedirectContext+com.TestResponseText, w.Body.String())

	path := req.URL.Path

	// Create a new recorder
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, path, nil)

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	// assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, com.TestUrlPath2, req.URL.Path)
	assert.Equal(t, newContext, w.Body.String())
}

func TestPathRewriter_IpWhitelistWithLocal(t *testing.T) {}

func TestPathRewriter_IpWhitelistWithOther(t *testing.T) {}
