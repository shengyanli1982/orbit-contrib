package ratelimiter

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	itl "github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter/internal"
	"golang.org/x/time/rate"
)

// IpRateLimiter 是一个IP限流器的结构体
// IpRateLimiter is a struct of IP rate limiter
type IpRateLimiter struct {
	// 缓存用于存储IP限流器
	// cache is used to store IP rate limiter
	cache *itl.Cache
	// 配置信息
	// config is used to store config information
	config *Config
	// 用于确保Stop方法只执行一次的同步锁
	// once is used to ensure that the Stop method is executed only once
	once sync.Once
}

// NewIpRateLimiter 创建一个新的IP限流器
// NewIpRateLimiter creates a new IP rate limiter
func NewIpRateLimiter(config *Config) *IpRateLimiter {
	return &IpRateLimiter{
		cache:  itl.NewCache(),        // 创建一个新的缓存
		config: isConfigValid(config), // 验证配置信息的有效性
		once:   sync.Once{},           // 初始化同步锁
	}
}

// GetLimiter 根据给定的key获取对应的限流器
// GetLimiter gets the corresponding limiter according to the given key
func (rl *IpRateLimiter) GetLimiter(key string) *rate.Limiter {
	if rl, ok := rl.cache.Get(key); ok {
		return rl.(*RateLimiter).GetLimiter()
	}
	return nil
}

// SetRate 设置所有限流器的速率
// SetRate sets the rate of all limiters
func (rl *IpRateLimiter) SetRate(rate float64) {
	// 获取缓存中的所有分段
	// Get all segments in the cache
	segments := rl.cache.Segments()
	// 遍历所有分段
	// Traverse all segments
	for i := 0; i < itl.SegmentSize; i++ {
		// 获取分段中的对应存储数据
		// Get the corresponding stored data in the segment
		data := segments[i].GetData()
		// 遍历所有key-value对
		// Traverse all key-value pairs
		for key, value := range data {
			// 获取限流器
			// Get the limiter
			limiter := value.(*RateLimiter)
			// 更新限流器的速率
			// Update the rate of the limiter
			limiter.SetRate(rate)
			// 更新key-value对
			// Update key-value pairs
			data[key] = limiter
		}
	}
}

// SetBurst 设置所有限流器的突发值
// SetBurst sets the burst of all limiters
func (rl *IpRateLimiter) SetBurst(burst int) {
	// 获取缓存中的所有分段
	// Get all segments in the cache
	segments := rl.cache.Segments()
	// 遍历所有分段
	// Traverse all segments
	for i := 0; i < itl.SegmentSize; i++ {
		// 获取分段中的对应存储数据
		// Get the corresponding stored data in the segment
		data := segments[i].GetData()
		// 遍历所有key-value对
		// Traverse all key-value pairs
		for key, value := range data {
			// 获取限流器
			// Get the limiter
			limiter := value.(*RateLimiter)
			// 更新限流器的突发值
			// Update the burst of the limiter
			limiter.SetBurst(burst)
			// 更新key-value对
			// Update key-value pairs
			data[key] = limiter
		}
	}
}

// HandlerFunc 返回一个gin中间件函数，用于处理请求
// HandlerFunc returns a gin middleware function for processing requests
func (rl *IpRateLimiter) HandlerFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 如果请求匹配限流器的配置，则进行限流
		// If the request matches the configuration of the limiter, then limit it
		if rl.config.match(context.Request) {
			// 获取客户端IP
			// Get the client IP
			clientIP := context.ClientIP()
			// 判断客户端IP在不在白名单中
			// Determine whether the client IP is in the whitelist
			if _, ok := rl.config.whitelist[clientIP]; !ok {
				// 获取限流器，如果不存在则创建一个新的限流器
				// Get the limiter, if it does not exist, create a new one
				limiter, _ := rl.cache.GetOrCreate(clientIP, func() any {
					// 从对象池中获取一个Element实例
					// Get an Element instance from the object pool
					element := itl.ElementPool.Get()
					// 设置Element的值为一个新的限流器
					// Set the value of Element to a new limiter
					element.(*itl.Element).SetValue(NewRateLimiter(rl.config))
					return element
				})
				// 限流器不允许该请求通过，如果不允许则返回429状态码
				// Limiter does not allow the request to pass, if not allowed, return 429 status code
				if !limiter.(*itl.Element).GetValue().(*RateLimiter).GetLimiter().Allow() {
					// 退出请求链
					// Exit the request chain
					context.Abort()
					context.String(http.StatusTooManyRequests, "[429] too many http requests, ip:"+clientIP+", method: "+context.Request.Method+", path: "+context.Request.URL.Path)
					// 触发限流回调函数
					// Trigger the rate limit callback function
					rl.config.callback.OnLimited(context.Request)
					return
				}
			}
		}
		// 继续请求链
		// Continue the request chain
		context.Next()
	}
}

// Stop 停止IP限流器
// Stop stops the IP rate limiter
func (rl *IpRateLimiter) Stop() {
	rl.once.Do(func() {
		// 停止缓存中的所有分段
		// Stop all segments in the cache
		rl.cache.Stop()
	})
}
