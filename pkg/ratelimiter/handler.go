package ratelimiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	gr "golang.org/x/time/rate"
)

type RateLimiter struct {
	config  *Config
	limiter *rate.Limiter
}

func NewRateLimiter(config *Config) *RateLimiter {
	config = isConfigValid(config)
	return &RateLimiter{
		config:  config,
		limiter: rate.NewLimiter(rate.Limit(config.rate), config.burst),
	}
}

func (rl *RateLimiter) GetLimiter() *rate.Limiter {
	return rl.limiter
}

func (rl *RateLimiter) SetRate(rate float64) {
	rl.config.rate = rate
	rl.limiter.SetLimit(gr.Limit(rl.config.rate))
}

func (rl *RateLimiter) SetBurst(burst int) {
	rl.config.burst = burst
	rl.limiter.SetBurst(rl.config.burst)
}

func (rl *RateLimiter) HandlerFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		if rl.config.match(context.Request) {
			clientIP := context.ClientIP()
			if _, ok := rl.config.whitelist[clientIP]; !(ok || rl.limiter.Allow()) {
				context.Abort()
				context.String(http.StatusTooManyRequests, "[429] too many http requests, method: "+context.Request.Method+", path: "+context.Request.URL.Path)
				rl.config.callback.OnLimited(context.Request)
				return
			}
		}
		context.Next()
	}
}
