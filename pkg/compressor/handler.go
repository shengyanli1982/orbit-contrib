package compressor

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var emptyContextWriter = gin.CreateTestContextOnly(httptest.NewRecorder(), gin.New()).Writer

type Compressor struct {
	config *Config
	pool   sync.Pool
}

func NewCompressor(config *Config) *Compressor {
	config = isConfigValid(config)
	compr := Compressor{
		config: config,
		pool: sync.Pool{
			New: func() interface{} {
				return config.createFunc(config, emptyContextWriter)
			},
		},
	}
	return &compr
}

// GetHandlerFunc returns a gin.HandlerFunc for processing requests
func (c *Compressor) GetHandlerFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !canCompressByHeader(context.Request) {
			context.Next()
			return
		}

		// 如果请求匹配限流器的配置，则进行限流
		// If the request matches the configuration of the rate limiter, rate limiting is performed
		if c.config.matchFunc(context.Request) {
			// 获取客户端 IP
			// Get the client IP
			clientIP := context.ClientIP()
			// 如果客户端 IP 不在白名单中，且限流器不允许该请求通过，则返回 429 状态码
			// If client IP is not in the whitelist and the rate limiter does not allow the request to pass, return 429 status code
			if _, ok := c.config.ipWhitelist[clientIP]; !ok {
				writer := c.pool.Get().(CodecWriter)
				if err := writer.Reset(context.Writer); err != nil {
					context.Abort()
					context.String(http.StatusInternalServerError, "[500] internal server error: compress error: "+err.Error()+", method: "+context.Request.Method+", path: "+context.Request.URL.Path)
					return
				}
				ctxWriter := context.Writer
				context.Writer = writer

				context.Header("Content-Encoding", "gzip")
				context.Header("Vary", "Accept-Encoding")

				context.Next()

				context.Header("Content-Length", strconv.Itoa(context.Writer.Size()))

				context.Writer = ctxWriter
				c.pool.Put(writer)
			}
		}
	}
}

// Stop stops the compressor
func (c *Compressor) Stop() {}

func canCompressByHeader(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Accept"), "text/event-stream") {
		return false
	}
	return true
}
