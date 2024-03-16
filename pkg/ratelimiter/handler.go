package ratelimiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	gr "golang.org/x/time/rate"
)

// RateLimiter 是一个结构体，包含配置和限流器
// RateLimiter is a struct that contains configuration and rate limiter
type RateLimiter struct {
	// config 是一个指向 Config 结构体的指针，用于存储配置信息
	// config is a pointer to the Config struct, used to store configuration information
	config *Config

	// limiter 是一个指向 rate.Limiter 结构体的指针，用于存储限流器
	// limiter is a pointer to the rate.Limiter struct, used to store rate limiter
	limiter *rate.Limiter
}

// NewRateLimiter 是一个函数，接收一个 Config 结构体的指针作为参数，返回一个新的 RateLimiter 结构体的指针
// NewRateLimiter is a function that takes a pointer to the Config struct as a parameter and returns a new pointer to the RateLimiter struct
func NewRateLimiter(config *Config) *RateLimiter {
	// 验证并获取有效的配置
	// Validate and get valid configuration
	config = isConfigValid(config)

	// 返回一个新的 RateLimiter 结构体的指针，其中包含配置和新的限流器
	// Return a new pointer to the RateLimiter struct, which includes configuration and a new rate limiter
	return &RateLimiter{
		// 设置配置
		// Set configuration
		config: config,

		// 创建并设置新的限流器，其中速率和突发来自配置
		// Create and set a new rate limiter, where the rate and burst come from the configuration
		limiter: rate.NewLimiter(rate.Limit(config.rate), config.burst),
	}
}

// GetLimiter 方法用于获取限流器
// The GetLimiter method is used to get the rate limiter
func (rl *RateLimiter) GetLimiter() *rate.Limiter {
	// 返回限流器
	// Return the rate limiter
	return rl.limiter
}

// SetRate 方法用于设置限流器的速率
// The SetRate method is used to set the rate of the rate limiter
func (rl *RateLimiter) SetRate(rate float64) {
	// 设置配置的速率
	// Set the rate of the configuration
	rl.config.rate = rate

	// 设置限流器的速率
	// Set the rate of the rate limiter
	rl.limiter.SetLimit(gr.Limit(rl.config.rate))
}

// SetBurst 方法用于设置限流器的突发流量
// The SetBurst method is used to set the burst traffic of the rate limiter
func (rl *RateLimiter) SetBurst(burst int) {
	// 设置配置的突发流量
	// Set the burst traffic of the configuration
	rl.config.burst = burst

	// 设置限流器的突发流量
	// Set the burst traffic of the rate limiter
	rl.limiter.SetBurst(rl.config.burst)
}

// Stop 方法用于停止限流器，但在这里没有实现任何功能
// The Stop method is used to stop the rate limiter, but it does not implement any functionality here
func (rl *RateLimiter) Stop() {}

// HandlerFunc 返回一个 gin.HandlerFunc，用于处理请求
// HandlerFunc returns a gin.HandlerFunc for processing requests
func (rl *RateLimiter) HandlerFunc() gin.HandlerFunc {

	// 返回一个闭包函数，该函数接收一个 gin.Context 参数，用于处理 HTTP 请求
	// Returns a closure function that takes a gin.Context parameter to handle HTTP requests
	return func(ctx *gin.Context) {

		// 如果请求匹配配置的匹配函数，则进行下一步处理
		// If the request matches the match function in the configuration, then proceed to the next step
		if rl.config.matchFunc(ctx.Request) {

			// 获取客户端 IP 地址
			// Get the client IP address
			clientIP := ctx.ClientIP()

			// 如果客户端 IP 地址不在配置的 IP 白名单中，并且限流器不允许新的请求，则进行下一步处理
			// If the client IP address is not in the IP whitelist in the configuration, and the rate limiter does not allow new requests, then proceed to the next step
			if _, ok := rl.config.ipWhitelist[clientIP]; !ok && !rl.limiter.Allow() {

				// 中止请求处理
				// Abort the request processing
				ctx.Abort()

				// 返回 429 错误，表示请求过多
				// Return a 429 error, indicating too many requests
				ctx.String(http.StatusTooManyRequests, "[429] too many http requests, method: "+ctx.Request.Method+", path: "+ctx.Request.URL.Path)

				// 调用配置的回调函数，处理限流事件
				// Call the callback function in the configuration to handle the rate limiting event
				rl.config.callback.OnLimited(ctx.Request)

				// 返回，不再执行后续代码
				// Return, no further code is executed
				return
			}
		}

		// 执行后续的请求处理
		// Execute subsequent request processing
		ctx.Next()
	}
}
