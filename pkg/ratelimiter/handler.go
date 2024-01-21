package ratelimiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	gr "golang.org/x/time/rate"
)

// RateLimiter 结构体用于实现限流功能
// RateLimiter is a struct for implementing rate limiting
type RateLimiter struct {
	// 配置信息
	// config is used to store config information
	config *Config
	// 限流器
	// limiter is used to store rate limiter
	limiter *rate.Limiter
}

// NewRateLimiter 创建一个新的 RateLimiter 实例
// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(config *Config) *RateLimiter {
	config = isConfigValid(config)
	return &RateLimiter{
		config:  config,
		limiter: rate.NewLimiter(rate.Limit(config.rate), config.burst),
	}
}

// GetLimiter 返回当前 RateLimiter 实例的限流器
// GetLimiter returns the limiter of the current RateLimiter instance
func (rl *RateLimiter) GetLimiter() *rate.Limiter {
	return rl.limiter
}

// SetRate 设置限流速率
// SetRate sets the rate of rate limiter
func (rl *RateLimiter) SetRate(rate float64) {
	rl.config.rate = rate
	// 重新设置限流器的速率
	// Reset the rate of rate limiter
	rl.limiter.SetLimit(gr.Limit(rl.config.rate))
}

// SetBurst 设置限流突发值
// SetBurst sets the burst of rate limiter
func (rl *RateLimiter) SetBurst(burst int) {
	rl.config.burst = burst
	// 重新设置限流器的突发值
	// Reset the burst of rate limiter
	rl.limiter.SetBurst(rl.config.burst)
}

// HandlerFunc 返回一个 gin.HandlerFunc，用于处理请求
// HandlerFunc returns a gin.HandlerFunc for processing requests
func (rl *RateLimiter) HandlerFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 如果请求匹配限流器的配置，则进行限流
		// If the request matches the configuration of the rate limiter, rate limiting is performed
		if rl.config.match(context.Request) {
			// 获取客户端 IP
			// Get the client IP
			clientIP := context.ClientIP()
			// 如果客户端 IP 不在白名单中，且限流器不允许该请求通过，则返回 429 状态码
			// If client IP is not in the whitelist and the rate limiter does not allow the request to pass, return 429 status code
			if _, ok := rl.config.whitelist[clientIP]; !ok && !rl.limiter.Allow() {
				// 退出请求链
				// Exit the request chain
				context.Abort()
				context.String(http.StatusTooManyRequests, "[429] too many http requests, method: "+context.Request.Method+", path: "+context.Request.URL.Path)
				// 调用回调函数
				// Call the callback function
				rl.config.callback.OnLimited(context.Request)
				return
			}
		}
		// 继续请求链
		// Continue the request chain
		context.Next()
	}
}

// Stop 停止限流器
// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {}
