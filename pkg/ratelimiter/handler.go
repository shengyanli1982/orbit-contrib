package ratelimiter

import (
	"github.com/gin-gonic/gin"
)

func NewRateLimiterHandlerFunc(config *Config) gin.HandlerFunc {
	config = isConfigValid(config)
	limiter := NewLimiter(config)

	return func(c *gin.Context) {
		if config.match(c.Request) {
			clientIp := c.ClientIP()
			if _, ok := config.whitelist[clientIp]; !ok && !limiter.Allow() {
				c.AbortWithStatus(429)
				config.callback.OnLimited(c.Request)
				return
			}
		}
		c.Next()
	}
}
