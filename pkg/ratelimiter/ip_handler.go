package ratelimiter

import (
	"math"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lxzan/memorycache"
)

var (
	defaultBucketNum = 128
)

type limiterPool struct {
	p sync.Pool
}

func NewLimiterPool() *limiterPool {
	return &limiterPool{
		p: sync.Pool{
			New: func() interface{} {
				return &Limiter{}
			},
		},
	}
}

func (lp *limiterPool) Get() *Limiter {
	return lp.p.Get().(*Limiter)
}

func (lp *limiterPool) Put(l *Limiter) {
	lp.p.Put(l)
}

func NewIpRateLimiterHandlerFunc(config *Config) gin.HandlerFunc {
	pool := NewLimiterPool()

	cache := memorycache.New[string, *Limiter](
		memorycache.WithBucketNum(defaultBucketNum),
		memorycache.WithBucketSize(0, math.MaxInt64),
		memorycache.WithInterval(5*time.Second, 30*time.Second),
	)

	cbFunc := func(entry *memorycache.Element[string, *Limiter], reason memorycache.Reason) {
		if reason == memorycache.ReasonExpired {
			pool.Put(entry.Value)
		}
	}

	return func(c *gin.Context) {
		if config.match(c.Request) {
			clientIp := c.ClientIP()
			if _, ok := config.whitelist[clientIp]; ok {
				c.Next()
				return
			}
			lr, ok := cache.GetOrCreateWithCallback(clientIp, NewLimiter(config), 15*time.Second, cbFunc)
			if ok {
				pool.Put(lr)
			}
			if !lr.Allow() {
				c.AbortWithStatus(429)
				config.callback.OnLimited(c.Request)
				return
			}
		}

		c.Next()
	}
}
