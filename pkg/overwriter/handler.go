package overwriter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PathOverwriter 结构体用于实现路径重写功能
// PathOverwriter is a struct for implementing path rewriting
type PathOverwriter struct {
	config *Config
}

// NewPathOverwriter 创建一个新的 PathOverwriter 实例
// NewPathOverwriter creates a new PathOverwriter instance
func NewPathOverwriter(config *Config) *PathOverwriter {
	config = isConfigValid(config)
	return &PathOverwriter{
		config: config,
	}
}

// HandlerFunc 返回一个 gin.HandlerFunc，用于处理请求
// HandlerFunc returns a gin.HandlerFunc for processing requests
func (p *PathOverwriter) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果请求匹配限流器的配置，则进行限流
		// If the request matches the configuration of the compressor, compression is performed
		if p.config.matchFunc(ctx.Request) {
			// 如果请求的 IP 不在白名单中，则进行限流
			// If the IP of the request is not in the whitelist, compression is performed
			clientIP := ctx.ClientIP()
			if _, ok := p.config.ipWhitelist[clientIP]; !ok {
				// 保存旧的请求路径
				// Save the old request path
				oldPath := ctx.Request.URL.Path

				// 重写请求路径, 如果重写失败, 则返回 500 错误
				// Rewrite the request path, if the rewrite fails, return a 500 error
				if ok, newPath := p.config.rewriteFunc(ctx.Request); ok {
					// 重定向到新的请求路径, 并修改请求路径, 以便后续中间件可以正确处理
					// Redirect to the new request path, and modify the request path so that subsequent middleware can handle it correctly
					ctx.Redirect(http.StatusTemporaryRedirect, newPath)
					ctx.Request.URL.Path = newPath

					// 调用回调函数
					// Call the callback function
					p.config.callback.OnPathRewrited(oldPath, newPath)
				}
			}
		}

		// 调用下一个中间件
		// Call the next middleware
		ctx.Next()
	}
}

// Stop 停止压缩器
// Stop stops the compressor
func (p *PathOverwriter) Stop() {}
