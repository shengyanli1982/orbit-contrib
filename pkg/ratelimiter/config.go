package ratelimiter

import (
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
)

var (
	// 默认每秒限制速率
	// Default limit rate per second
	DefaultLimitRatePerSecond = float64(1)

	// 默认限制突发数
	// Default limit burst
	DefaultLimitBurst = 1
)

// HttpRequestHeaderMatchFunc 是一个匹配函数，用于匹配请求头

// Config 是一个配置结构体
// Config is a struct of config
type Config struct {
	// 速率
	// Rate
	rate float64

	// 突发数
	// Burst
	burst int

	// 白名单
	// Whitelist
	ipWhitelist map[string]struct{}

	// 匹配函数
	// Match function
	matchFunc com.HttpRequestHeaderMatchFunc

	// 回调函数
	// Callback
	callback Callback
}

// NewConfig 创建一个新的配置实例
// NewConfig creates a new config instance
func NewConfig() *Config {
	return &Config{
		rate:        DefaultLimitRatePerSecond,
		burst:       DefaultLimitBurst,
		matchFunc:   com.DefaultLimitMatchFunc,
		ipWhitelist: com.DefaultIpWhitelist,
		callback:    &emptyCallback{},
	}
}

// DefaultConfig 返回默认配置实例
// DefaultConfig returns the default config instance
func DefaultConfig() *Config {
	return NewConfig()
}

// WithCallback 设置回调函数
// WithCallback sets the callback function
func (c *Config) WithCallback(callback Callback) *Config {
	c.callback = callback
	return c
}

// WithRate 设置速率
// WithRate sets the rate
func (c *Config) WithRate(rate float64) *Config {
	c.rate = rate
	return c
}

// WithBurst 设置突发数
// WithBurst sets the burst
func (c *Config) WithBurst(burst int) *Config {
	c.burst = burst
	return c
}

// WithMatchFunc 设置匹配函数
// WithMatchFunc sets the match function
func (c *Config) WithMatchFunc(fn com.HttpRequestHeaderMatchFunc) *Config {
	c.matchFunc = fn
	return c
}

// WithIpWhitelist 设置白名单
// WithIpWhitelist sets the whitelist
func (c *Config) WithIpWhitelist(whitelist []string) *Config {
	for _, ip := range whitelist {
		c.ipWhitelist[ip] = com.Empty
	}
	return c
}

// isConfigValid 检查配置是否有效，如果无效则返回默认配置实例
// isConfigValid checks whether the config is valid, and returns the default config instance if it is invalid
func isConfigValid(config *Config) *Config {
	if config != nil {
		if config.rate <= 0 {
			config.rate = DefaultLimitRatePerSecond
		}
		if config.burst <= 0 {
			config.burst = DefaultLimitBurst
		}
		if config.matchFunc == nil {
			config.matchFunc = com.DefaultLimitMatchFunc
		}
		if config.ipWhitelist == nil {
			config.ipWhitelist = com.DefaultIpWhitelist
		}
		if config.callback == nil {
			config.callback = &emptyCallback{}
		}
	} else {
		config = DefaultConfig()
	}

	return config
}
