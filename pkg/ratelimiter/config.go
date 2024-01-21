package ratelimiter

import "net/http"

// empty 是一个空结构体，用于实现空的值
// empty is an empty struct for implementing empty value
var empty = struct{}{}

var (
	// 默认本地IP地址
	// Default local IP address
	DefaultLocalIpAddress = "127.0.0.1"
	// 默认本地IPv6地址
	// Default local IPv6 address
	DefaultLocalIpv6Address = "::1"
	// 默认每秒限制速率
	// Default limit rate per second
	DefaultLimitRatePerSecond = float64(1)
	// 默认限制突发数
	// Default limit burst
	DefaultLimitBurst = 1
	// 默认匹配函数
	// Default match function
	DefaultLimitMatchFunc = func(header *http.Request) bool { return true }
	// 默认IP白名单
	// Default IP whitelist
	DefaultIpWhitelist = map[string]struct{}{
		DefaultLocalIpAddress:   empty,
		DefaultLocalIpv6Address: empty,
	}
)

// HttpRequestHeaderMatchFunc 是一个匹配函数，用于匹配请求头
// HttpRequestHeaderMatchFunc is a match function for matching request headers
type HttpRequestHeaderMatchFunc func(header *http.Request) bool

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
	whitelist map[string]struct{}
	// 匹配函数
	// Match function
	match HttpRequestHeaderMatchFunc
	// 回调函数
	// Callback
	callback Callback
}

// NewConfig 创建一个新的配置实例
// NewConfig creates a new config instance
func NewConfig() *Config {
	return &Config{
		rate:      DefaultLimitRatePerSecond,
		burst:     DefaultLimitBurst,
		match:     DefaultLimitMatchFunc,
		whitelist: DefaultIpWhitelist,
		callback:  &emptyCallback{},
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
func (c *Config) WithMatchFunc(match HttpRequestHeaderMatchFunc) *Config {
	c.match = match
	return c
}

// WithWhitelist 设置白名单
// WithWhitelist sets the whitelist
func (c *Config) WithWhitelist(whitelist []string) *Config {
	for _, ip := range whitelist {
		c.whitelist[ip] = empty
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
		if config.match == nil {
			config.match = DefaultLimitMatchFunc
		}
		if config.whitelist == nil {
			config.whitelist = DefaultIpWhitelist
		}
		if config.callback == nil {
			config.callback = &emptyCallback{}
		}
	} else {
		config = DefaultConfig()
	}

	return config
}
