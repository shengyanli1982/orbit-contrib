package ratelimiter

import (
	"golang.org/x/time/rate"
)

type Limiter struct {
	config *Config
	lr     *rate.Limiter
}

func NewLimiter(config *Config) *Limiter {
	return &Limiter{
		config: config,
		lr:     rate.NewLimiter(rate.Limit(config.rate), config.burst),
	}
}

func (l *Limiter) Allow() bool {
	return l.lr.Allow()
}

func (l *Limiter) SetRate(r float64) {
	l.config.rate = r
	l.lr.SetLimit(rate.Limit(l.config.rate))
}

func (l *Limiter) SetBurst(b int) {
	l.config.burst = b
	l.lr.SetBurst(l.config.burst)
}
