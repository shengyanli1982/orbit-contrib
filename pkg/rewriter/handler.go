package rewriter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PathRewriter 结构体用于实现路径重写功能
// PathRewriter is a struct for implementing path rewriting
type PathRewriter struct {
	config *Config
}

// NewPathRewriter 创建一个新的 PathRewriter 实例
// NewPathRewriter creates a new PathRewriter instance
func NewPathRewriter(config *Config) *PathRewriter {
	config = isConfigValid(config)
	return &PathRewriter{
		config: config,
	}
}

// HandlerFunc 返回一个 gin.HandlerFunc，用于处理请求
// HandlerFunc returns a gin.HandlerFunc for processing requests
func (p *PathRewriter) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果请求匹配限流器的配置，则进行限流
		// If the request matches the configuration of the compressor, compression is performed
		if p.config.matchFunc(ctx.Request) {
			// 如果请求的 IP 不在白名单中，则进行限流
			// If the IP of the request is not in the whitelist, compression is performed
			clientIP := ctx.ClientIP()
			if _, ok := p.config.ipWhitelist[clientIP]; !ok {
				// 如果请求的路径需要重写，则进行重写
				// If the path of the request needs to be rewritten, the path is rewritten
				if ok, newPath := p.config.rewriteFunc(ctx.Request.URL); ok {
					// 保存旧的请求路径
					// Save the old request path
					oldPath := ctx.Request.URL.Path

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
func (p *PathRewriter) Stop() {}
