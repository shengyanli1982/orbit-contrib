package compressor

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Compressor 是一个通用压缩器，包含配置和同步池
// Compressor is a common compressor, containing configuration and sync pool
type Compressor struct {
	// 配置，包含压缩等级、IP白名单、匹配函数和创建压缩写入器的函数
	// Configuration, including compression level, IP whitelist, match function and function to create a compression writer
	config *Config

	// 同步池，用于存储和复用压缩写入器
	// Sync pool, used to store and reuse compression writers
	pool sync.Pool
}

// NewCompressor 创建一个新的压缩器，包含有效的配置和同步池
// NewCompressor creates a new compressor, including valid configuration and sync pool
func NewCompressor(config *Config) *Compressor {
	// 检查配置是否有效，如果无效则使用默认配置
	// Check if the configuration is valid, if not, use the default configuration
	config = isConfigValid(config)

	// 返回一个新的压缩器实例
	// Returns a new compressor instance
	return &Compressor{
		// 设置配置
		// Sets the configuration
		config: config,

		// 创建一个同步池，池中的新元素由 createFunc 函数创建
		// Creates a sync pool, new elements in the pool are created by the createFunc function
		pool: sync.Pool{
			New: func() interface{} {
				// 创建一个新的压缩写入器
				// Creates a new compression writer
				return config.createFunc(config, nil)
			},
		},
	}
}

// HandlerFunc 返回一个 gin.HandlerFunc，用于处理请求
// HandlerFunc returns a gin.HandlerFunc for processing requests
func (c *Compressor) HandlerFunc() gin.HandlerFunc {

	// 返回一个闭包函数，该函数接收一个 gin.Context 参数，用于处理 HTTP 请求
	// Returns a closure function that takes a gin.Context parameter to handle HTTP requests
	return func(ctx *gin.Context) {

		// 如果请求匹配配置的匹配函数，或者请求头允许压缩，则进行下一步处理
		// If the request matches the match function in the configuration, or the request header allows compression, then proceed to the next step
		if c.config.matchFunc(ctx.Request) || canCompressByHeader(ctx.Request) {

			// 获取客户端 IP 地址
			// Get the client IP address
			clientIP := ctx.ClientIP()

			// 如果客户端 IP 地址不在配置的 IP 白名单中，则进行下一步处理
			// If the client IP address is not in the IP whitelist in the configuration, then proceed to the next step
			if _, ok := c.config.ipWhitelist[clientIP]; !ok {

				// 从同步池中获取一个压缩写入器
				// Get a compression writer from the sync pool
				writer := c.pool.Get().(CodecWriter)

				// 使用 defer 语句在函数返回时执行一些清理操作
				// Use the defer statement to perform some cleanup operations when the function returns
				defer func() {
					// 重置压缩写入器的写入器为 io.Discard，忽略所有写入的数据
					// Reset the writer of the compression writer to io.Discard, ignoring all written data
					_ = writer.ResetCompressWriter(io.Discard)

					// 重置响应写入器为 nil
					// Reset the response writer to nil
					_ = writer.ResetResponseWriter(nil)

					// 将压缩写入器放回同步池
					// Put the compression writer back into the sync pool
					c.pool.Put(writer)
				}()

				// 重置压缩写入器的写入器为 ctx.Writer，如果出错则返回 500 错误
				// Reset the writer of the compression writer to ctx.Writer, if an error occurs, return a 500 error
				if err := writer.ResetCompressWriter(ctx.Writer); err != nil {
					// 中止请求处理
					// Abort the request processing
					ctx.Abort()

					// 返回 500 错误
					// Return a 500 error
					ctx.String(http.StatusInternalServerError, "[500] internal server error: compress writer error: "+err.Error()+", method: "+ctx.Request.Method+", path: "+ctx.Request.URL.Path)

					// 返回，不再执行后续代码
					// Return, no further code is executed
					return
				}

				// 重置响应写入器为 ctx.Writer，如果出错则返回 500 错误
				// Reset the response writer to ctx.Writer, if an error occurs, return a 500 error
				if err := writer.ResetResponseWriter(ctx.Writer); err != nil {
					// 中止请求处理
					// Abort the request processing
					ctx.Abort()

					// 返回 500 错误
					// Return a 500 error
					ctx.String(http.StatusInternalServerError, "[500] internal server error: response writer error: "+err.Error()+", method: "+ctx.Request.Method+", path: "+ctx.Request.URL.Path)

					// 返回，不再执行后续代码
					// Return, no further code is executed
					return
				}

				// 保存原来的响应写入器
				// Save the original response writer
				ctxWriter := ctx.Writer

				// 将 ctx.Writer 替换为压缩写入器
				// Replace ctx.Writer with the compression writer
				ctx.Writer = writer

				// 设置响应头的 "Content-Encoding" 字段为压缩写入器的 Content-Encoding
				// Set the "Content-Encoding" field of the response header to the Content-Encoding of the compression writer
				ctx.Header("Content-Encoding", writer.ContentEncoding())

				// 设置响应头的 "Vary" 字段为 "Accept-Encoding"
				// Set the "Vary" field of the response header to "Accept-Encoding"
				ctx.Header("Vary", "Accept-Encoding")

				// 执行后续的请求处理
				// Execute subsequent request processing
				ctx.Next()

				// 设置响应头的 "Content-Length" 字段为响应的大小
				// Set the "Content-Length" field of the response header to the size of the response
				ctx.Header("Content-Length", strconv.Itoa(ctx.Writer.Size()))

				// 停止压缩写入器的操作
				// Stop the operation of the compression writer
				writer.Stop()

				// 将 ctx.Writer 替换回原来的响应写入器
				// Replace ctx.Writer back to the original response writer
				ctx.Writer = ctxWriter
			}
		}
	}
}

// Stop 停止压缩器的操作，这个函数目前是空的，没有具体的实现
// Stop stops the operation of the compressor, this function is currently empty, without specific implementation
func (c *Compressor) Stop() {}

// canCompressByHeader 根据请求头判断是否可以压缩。首先获取请求头的 "Accept-Encoding"、"Connection" 和 "Accept" 字段，
// 如果 "Accept-Encoding" 字段不包含 "gzip"，或者 "Connection" 字段包含 "Upgrade"，或者 "Accept" 字段包含 "text/event-stream"，
// 则返回 false，表示不能压缩；否则返回 true，表示可以压缩。
// canCompressByHeader determines whether compression is possible based on the request header. It first gets the "Accept-Encoding", "Connection" and "Accept" fields of the request header,
// if the "Accept-Encoding" field does not contain "gzip", or the "Connection" field contains "Upgrade", or the "Accept" field contains "text/event-stream",
// then it returns false, indicating that compression is not possible; otherwise, it returns true, indicating that compression is possible.
func canCompressByHeader(req *http.Request) bool {
	// 获取请求头的 "Accept-Encoding" 字段
	// Gets the "Accept-Encoding" field of the request header
	acceptEncoding := req.Header.Get("Accept-Encoding")

	// 获取请求头的 "Connection" 字段
	// Gets the "Connection" field of the request header
	connection := req.Header.Get("Connection")

	// 获取请求头的 "Accept" 字段
	// Gets the "Accept" field of the request header
	accept := req.Header.Get("Accept")

	// 判断是否可以压缩
	// Determines whether compression is possible
	if !strings.Contains(acceptEncoding, "gzip") ||
		strings.Contains(connection, "Upgrade") ||
		strings.Contains(accept, "text/event-stream") {
		// 如果不能压缩，返回 false
		// If compression is not possible, return false
		return false
	}

	// 如果可以压缩，返回 true
	// If compression is possible, return true
	return true
}
