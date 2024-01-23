package compressor

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Compressor 是一个通用压缩
// Compressor is a common compressor
type Compressor struct {
	config *Config
	pool   sync.Pool
}

// NewCompressor 创建一个新的压缩器
// NewCompressor creates a new compressor
func NewCompressor(config *Config) *Compressor {
	config = isConfigValid(config)
	return &Compressor{
		config: config,
		pool: sync.Pool{
			New: func() interface{} {
				return config.createFunc(config, nil)
			},
		},
	}
}

// HandlerFunc 返回一个 gin.HandlerFunc，用于处理请求
// HandlerFunc returns a gin.HandlerFunc for processing requests
func (c *Compressor) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 如果请求匹配压缩器的配置，则进行压缩
		// If the request matches the configuration of the compressor, compression is performed
		if c.config.matchFunc(ctx.Request) || canCompressByHeader(ctx.Request) {

			// 如果请求的 IP 不在白名单中，则进行压缩
			// If the IP of the request is not in the whitelist, compression is performed
			clientIP := ctx.ClientIP()
			if _, ok := c.config.ipWhitelist[clientIP]; !ok {

				// 从池中获取一个压缩写入器
				// Get a compression writer from the pool
				writer := c.pool.Get().(CodecWriter)

				// 重置压缩写入器，归还给池
				// Reset the compression writer and return it to the pool
				defer func() {
					_ = writer.ResetCompressWriter(io.Discard)
					_ = writer.ResetResponseWriter(nil)
					c.pool.Put(writer)
				}()

				// 重置压缩写入器
				// Reset the compression writer
				if err := writer.ResetCompressWriter(ctx.Writer); err != nil {
					ctx.Abort()
					ctx.String(http.StatusInternalServerError, "[500] internal server error: compress writer error: "+err.Error()+", method: "+ctx.Request.Method+", path: "+ctx.Request.URL.Path)
					return
				}

				// 重置响应写入器
				// Reset the response writer
				if err := writer.ResetResponseWriter(ctx.Writer); err != nil {
					ctx.Abort()
					ctx.String(http.StatusInternalServerError, "[500] internal server error: response writer error: "+err.Error()+", method: "+ctx.Request.Method+", path: "+ctx.Request.URL.Path)
					return
				}

				// 保存原有响应写入器，替换成 GZip 写入器
				// Save the original response writer and replace it with the GZip writer
				ctxWriter := ctx.Writer
				ctx.Writer = writer

				// 设置响应头
				// Set the response header
				ctx.Header("Content-Encoding", "gzip")
				ctx.Header("Vary", "Accept-Encoding")

				// 调用下一个中间件
				// Call the next middleware
				ctx.Next()

				// 设置响应头
				// Set the response header
				ctx.Header("Content-Length", strconv.Itoa(ctx.Writer.Size()))

				// 停止压缩器，刷新数据到底层
				// Stop the compressor and flush the data to the underlying
				writer.Stop()

				// 恢复原有响应写入器
				// Restore the original response writer
				ctx.Writer = ctxWriter
			}
		}
	}
}

// Stop 停止压缩器
// Stop stops the compressor
func (c *Compressor) Stop() {}

// canCompressByHeader 根据请求头判断是否可以压缩
// canCompressByHeader determines whether compression is possible based on the request header
func canCompressByHeader(req *http.Request) bool {
	acceptEncoding := req.Header.Get("Accept-Encoding")
	connection := req.Header.Get("Connection")
	accept := req.Header.Get("Accept")

	if !strings.Contains(acceptEncoding, "gzip") ||
		strings.Contains(connection, "Upgrade") ||
		strings.Contains(accept, "text/event-stream") {
		return false
	}
	return true
}
