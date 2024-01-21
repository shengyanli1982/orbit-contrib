package ratelimiter

import "net/http"

var empty = struct{}{}

var (
	DefaultLocalIpAddress     = "127.0.0.1"                                     // 默认本地IP地址
	DefaultLocalIpv6Address   = "::1"                                           // 默认本地IPv6地址
	DefaultLimitRatePerSecond = float64(1)                                      // 默认每秒限制速率
	DefaultLimitBurst         = 1                                               // 默认限制突发数
	DefaultLimitMatchFunc     = func(header *http.Request) bool { return true } // 默认匹配函数
	DefaultIpWhitelist        = map[string]struct{}{
		DefaultLocalIpAddress:   empty,
		DefaultLocalIpv6Address: empty,
	} // 默认IP白名单
)

type HttpRequestHeaderMatchFunc func(header *http.Request) bool

type Config struct {
	rate      float64                    // 速率
	burst     int                        // 突发数
	whitelist map[string]struct{}        // 白名单
	match     HttpRequestHeaderMatchFunc // 匹配函数
	callback  Callback                   // 回调函数
}

// NewConfig 创建一个新的配置实例
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
func DefaultConfig() *Config {
	return NewConfig()
}

// WithCallback 设置回调函数
func (c *Config) WithCallback(callback Callback) *Config {
	c.callback = callback
	return c
}

// WithRate 设置速率
func (c *Config) WithRate(rate float64) *Config {
	c.rate = rate
	return c
}

// WithBurst 设置突发数
func (c *Config) WithBurst(burst int) *Config {
	c.burst = burst
	return c
}

// WithMatchFunc 设置匹配函数
func (c *Config) WithMatchFunc(match HttpRequestHeaderMatchFunc) *Config {
	c.match = match
	return c
}

// WithWhitelist 设置白名单
func (c *Config) WithWhitelist(whitelist []string) *Config {
	for _, ip := range whitelist {
		c.whitelist[ip] = empty
	}
	return c
}

// isConfigValid 检查配置是否有效，如果无效则返回默认配置实例
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
