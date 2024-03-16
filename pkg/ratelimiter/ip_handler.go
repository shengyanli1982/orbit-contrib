package ratelimiter

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	itl "github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter/internal"
	"golang.org/x/time/rate"
)

// IpRateLimiter 是一个结构体，包含配置和限流器
// IpRateLimiter is a struct that contains configuration and rate limiter
type IpRateLimiter struct {
	// cache 是一个指向 itl.Cache 的指针，用于存储限流器
	// cache is a pointer to itl.Cache, used to store rate limiters
	cache *itl.Cache

	// config 是一个指向 Config 的指针，用于存储限流器的配置
	// config is a pointer to Config, used to store the configuration of the rate limiter
	config *Config

	// once 是一个 sync.Once 类型的变量，用于确保某些操作只执行一次
	// once is a variable of type sync.Once, used to ensure that certain operations are performed only once
	once sync.Once
}

// NewIpRateLimiter 是一个函数，接收一个 Config 结构体的指针作为参数，返回一个新的 IpRateLimiter 结构体的指针
// NewIpRateLimiter is a function that takes a pointer to the Config struct as a parameter and returns a new pointer to the IpRateLimiter struct
func NewIpRateLimiter(config *Config) *IpRateLimiter {
	// 返回一个新的 IpRateLimiter 结构体的指针
	// Return a new pointer to the IpRateLimiter struct
	return &IpRateLimiter{
		// 初始化 cache 为一个新的 itl.Cache
		// Initialize cache as a new itl.Cache
		cache: itl.NewCache(),

		// 检查 config 是否有效，如果有效则返回 config，否则返回默认配置
		// Check if config is valid, if it is valid then return config, otherwise return the default configuration
		config: isConfigValid(config),

		// 初始化 once 为一个新的 sync.Once
		// Initialize once as a new sync.Once
		once: sync.Once{},
	}
}

// GetLimiter 方法用于根据键获取限流器
// The GetLimiter method is used to get the rate limiter based on the key
func (rl *IpRateLimiter) GetLimiter(key string) *rate.Limiter {
	// 从缓存中获取限流器
	// Get the rate limiter from the cache
	if rl, ok := rl.cache.Get(key); ok {
		// 如果存在，则返回限流器
		// If it exists, return the rate limiter
		return rl.(*RateLimiter).GetLimiter()
	}

	// 如果不存在，则返回 nil
	// If it does not exist, return nil
	return nil
}

// Stop 方法用于停止 IP 限流器
// The Stop method is used to stop the IP rate limiter
func (rl *IpRateLimiter) Stop() {
	// 使用 sync.Once 确保缓存只被停止一次
	// Use sync.Once to ensure that the cache is stopped only once
	rl.once.Do(func() {
		// 停止缓存
		// Stop the cache
		rl.cache.Stop()
	})
}

// SetRate 方法用于设置限流器的速率
// The SetRate method is used to set the rate of the rate limiter
func (rl *IpRateLimiter) SetRate(rate float64) {
	// 获取缓存的所有段
	// Get all segments of the cache
	segments := rl.cache.Segments()

	// 遍历所有段
	// Traverse all segments
	for i := 0; i < itl.SegmentSize; i++ {
		// 获取当前段的数据
		// Get the data of the current segment
		data := segments[i].GetData()

		// 遍历当前段的所有数据
		// Traverse all data of the current segment
		for key, value := range data {
			// 获取当前数据的限流器
			// Get the rate limiter of the current data
			limiter := value.(*RateLimiter)

			// 设置限流器的速率
			// Set the rate of the rate limiter
			limiter.SetRate(rate)

			// 更新当前数据的限流器
			// Update the rate limiter of the current data
			data[key] = limiter
		}
	}
}

// SetBurst 方法用于设置限流器的突发流量
// The SetBurst method is used to set the burst traffic of the rate limiter
func (rl *IpRateLimiter) SetBurst(burst int) {
	// 获取缓存的所有段
	// Get all segments of the cache
	segments := rl.cache.Segments()

	// 遍历所有段
	// Traverse all segments
	for i := 0; i < itl.SegmentSize; i++ {
		// 获取当前段的数据
		// Get the data of the current segment
		data := segments[i].GetData()

		// 遍历当前段的所有数据
		// Traverse all data of the current segment
		for key, value := range data {
			// 获取当前数据的限流器
			// Get the rate limiter of the current data
			limiter := value.(*RateLimiter)

			// 设置限流器的突发流量
			// Set the burst traffic of the rate limiter
			limiter.SetBurst(burst)

			// 更新当前数据的限流器
			// Update the rate limiter of the current data
			data[key] = limiter
		}
	}
}

// HandlerFunc 返回一个 gin.HandlerFunc，用于处理请求
// HandlerFunc returns a gin.HandlerFunc for processing requests
func (rl *IpRateLimiter) HandlerFunc() gin.HandlerFunc {

	// 返回一个闭包函数，该函数接收一个 gin.Context 参数，用于处理 HTTP 请求
	// Returns a closure function that takes a gin.Context parameter to handle HTTP requests
	return func(ctx *gin.Context) {

		// 如果请求匹配配置的匹配函数，则进行下一步处理
		// If the request matches the match function in the configuration, then proceed to the next step
		if rl.config.matchFunc(ctx.Request) {

			// 获取客户端 IP 地址
			// Get the client IP address
			clientIP := ctx.ClientIP()

			// 如果客户端 IP 地址不在配置的 IP 白名单中，则进行下一步处理
			// If the client IP address is not in the IP whitelist in the configuration, then proceed to the next step
			if _, ok := rl.config.ipWhitelist[clientIP]; !ok {

				// 从缓存中获取或创建一个限流器
				// Get or create a rate limiter from the cache
				limiter, _ := rl.cache.GetOrCreate(clientIP, func() any {
					// 从元素池中获取一个元素，并设置其值为一个新的限流器
					// Get an element from the element pool and set its value to a new rate limiter
					element := itl.ElementPool.Get()

					// 将元素的值设置为一个新的限流器，该限流器的配置为 rl.config
					// Set the value of the element to a new rate limiter, the configuration of this rate limiter is rl.config
					element.(*itl.Element).SetValue(NewRateLimiter(rl.config))

					// 返回元素，该元素将被添加到缓存中
					// Return the element, this element will be added to the cache
					return element
				})

				// 如果限流器不允许新的请求，则中止请求处理，并返回 429 错误
				// If the rate limiter does not allow new requests, abort the request processing and return a 429 error
				if !limiter.(*itl.Element).GetValue().(*RateLimiter).GetLimiter().Allow() {
					// 中止请求处理
					// Abort the request processing
					ctx.Abort()

					// 返回 429 错误
					// Return a 429 error
					ctx.String(http.StatusTooManyRequests, "[429] too many http requests, ip:"+clientIP+", method: "+ctx.Request.Method+", path: "+ctx.Request.URL.Path)

					// 调用回调函数，处理被限制的请求
					// Call the callback function to handle the limited request
					rl.config.callback.OnLimited(ctx.Request)

					// 返回，不再执行后续代码
					// Return, no further code is executed
					return
				}
			}
		}

		// 执行后续的请求处理
		// Execute subsequent request processing
		ctx.Next()
	}
}
