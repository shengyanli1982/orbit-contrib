package ratelimiter

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	itl "github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter/internal"
	"golang.org/x/time/rate"
)

type IpRateLimiter struct {
	cache  *itl.Cache
	config *Config
	once   sync.Once
}

func NewIpRateLimiter(config *Config) *IpRateLimiter {
	return &IpRateLimiter{
		cache:  itl.NewCache(),
		config: isConfigValid(config),
		once:   sync.Once{},
	}
}

func (rl *IpRateLimiter) GetLimiter(key string) *rate.Limiter {
	if rl, ok := rl.cache.Get(key); ok {
		return rl.(*RateLimiter).GetLimiter()
	}
	return nil
}

func (rl *IpRateLimiter) SetRate(rate float64) {
	segments := rl.cache.Segments()
	for i := 0; i < itl.SegmentSize; i++ {
		data := segments[i].GetData()
		for key, value := range data {
			limiter := value.(*RateLimiter)
			limiter.SetRate(rate)
			data[key] = limiter
		}
	}
}

func (rl *IpRateLimiter) SetBurst(burst int) {
	segments := rl.cache.Segments()
	for i := 0; i < itl.SegmentSize; i++ {
		data := segments[i].GetData()
		for key, value := range data {
			limiter := value.(*RateLimiter)
			limiter.SetBurst(burst)
			data[key] = limiter
		}
	}
}

func (rl *IpRateLimiter) HandlerFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		if rl.config.match(context.Request) {
			clientIP := context.ClientIP()
			if _, ok := rl.config.whitelist[clientIP]; !ok {
				limiter, _ := rl.cache.GetOrCreate(clientIP, func() any {
					element := itl.ElementPool.Get()
					element.(*itl.Element).SetValue(NewRateLimiter(rl.config))
					return element
				})
				if !limiter.(*itl.Element).GetValue().(*RateLimiter).GetLimiter().Allow() {
					context.Abort()
					context.String(http.StatusTooManyRequests, "[429] too many http requests, ip:"+clientIP+", method: "+context.Request.Method+", path: "+context.Request.URL.Path)
					rl.config.callback.OnLimited(context.Request)
					return
				}
			}
		}
		context.Next()
	}
}

func (rl *IpRateLimiter) Stop() {
	rl.once.Do(func() {
		rl.cache.Stop()
	})
}
