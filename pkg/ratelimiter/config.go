package ratelimiter

import (
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
)

// DefaultLimitRatePerSecond 是默认的每秒限制速率
// DefaultLimitRatePerSecond is the default limit rate per second
var DefaultLimitRatePerSecond = float64(1)

// DefaultLimitBurst 是默认的限制突发
// DefaultLimitBurst is the default limit burst
var DefaultLimitBurst = 1

// Config 是配置结构体，包含速率、突发、IP白名单、匹配函数和回调
// Config is the configuration structure, including rate, burst, IP whitelist, match function and callback
type Config struct {
	// rate 是每秒限制速率
	// rate is the limit rate per second
	rate float64

	// burst 是限制突发
	// burst is the limit burst
	burst int

	// ipWhitelist 是IP白名单
	// ipWhitelist is the IP whitelist
	ipWhitelist map[string]struct{}

	// matchFunc 是匹配函数
	// matchFunc is the match function
	matchFunc com.HttpRequestHeaderMatchFunc

	// callback 是回调
	// callback is the callback
	callback Callback
}

// NewConfig 创建一个新的配置，包含默认的速率、突发、匹配函数、IP白名单和回调
// NewConfig creates a new configuration, including default rate, burst, match function, IP whitelist and callback
func NewConfig() *Config {
	return &Config{
		// 设置速率为默认的每秒限制速率
		// Sets the rate to the default limit rate per second
		rate: DefaultLimitRatePerSecond,

		// 设置突发为默认的限制突发
		// Sets the burst to the default limit burst
		burst: DefaultLimitBurst,

		// 设置匹配函数为默认的限制匹配函数
		// Sets the match function to the default limit match function
		matchFunc: com.DefaultLimitMatchFunc,

		// 设置IP白名单为默认的IP白名单
		// Sets the IP whitelist to the default IP whitelist
		ipWhitelist: com.DefaultIpWhitelist,

		// 设置回调为空回调
		// Sets the callback to the empty callback
		callback: &emptyCallback{},
	}
}

// DefaultConfig 是一个函数，返回一个新的默认配置
// DefaultConfig is a function that returns a new default configuration
func DefaultConfig() *Config {
	return NewConfig()
}

// WithCallback 是一个方法，接收一个回调作为参数，设置配置的回调，并返回配置
// WithCallback is a method that takes a callback as a parameter, sets the callback of the configuration, and returns the configuration
func (c *Config) WithCallback(callback Callback) *Config {
	c.callback = callback
	return c
}

// WithRate 是一个方法，接收一个浮点数作为参数，设置配置的速率，并返回配置
// WithRate is a method that takes a float as a parameter, sets the rate of the configuration, and returns the configuration
func (c *Config) WithRate(rate float64) *Config {
	c.rate = rate
	return c
}

// WithBurst 是一个方法，接收一个整数作为参数，设置配置的突发，并返回配置
// WithBurst is a method that takes an integer as a parameter, sets the burst of the configuration, and returns the configuration
func (c *Config) WithBurst(burst int) *Config {
	c.burst = burst
	return c
}

// WithMatchFunc 是一个方法，接收一个匹配函数作为参数，设置配置的匹配函数，并返回配置
// WithMatchFunc is a method that takes a match function as a parameter, sets the match function of the configuration, and returns the configuration
func (c *Config) WithMatchFunc(fn com.HttpRequestHeaderMatchFunc) *Config {
	c.matchFunc = fn
	return c
}

// WithIpWhitelist 是一个方法，接收一个字符串切片作为参数，设置配置的 IP 白名单，并返回配置
// WithIpWhitelist is a method that takes a slice of strings as a parameter, sets the IP whitelist of the configuration, and returns the configuration
func (c *Config) WithIpWhitelist(whitelist []string) *Config {
	for _, ip := range whitelist {
		c.ipWhitelist[ip] = com.Empty
	}
	return c
}

// isConfigValid 是一个函数，它接收一个 Config 指针作为参数，检查配置是否有效，如果无效则设置为默认值，最后返回有效的配置
// isConfigValid is a function that takes a pointer to Config as a parameter, checks if the configuration is valid, if not, sets it to the default value, and finally returns the valid configuration
func isConfigValid(config *Config) *Config {
	// 如果配置不为 nil，则进行检查和设置
	// If the configuration is not nil, then check and set it
	if config != nil {
		// 如果速率小于等于 0，则设置为默认的每秒限制速率
		// If the rate is less than or equal to 0, set it to the default limit rate per second
		if config.rate <= 0 {
			config.rate = DefaultLimitRatePerSecond
		}

		// 如果突发小于等于 0，则设置为默认的限制突发
		// If the burst is less than or equal to 0, set it to the default limit burst
		if config.burst <= 0 {
			config.burst = DefaultLimitBurst
		}

		// 如果匹配函数为 nil，则设置为默认的限制匹配函数
		// If the match function is nil, set it to the default limit match function
		if config.matchFunc == nil {
			config.matchFunc = com.DefaultLimitMatchFunc
		}

		// 如果 IP 白名单为 nil，则设置为默认的 IP 白名单
		// If the IP whitelist is nil, set it to the default IP whitelist
		if config.ipWhitelist == nil {
			config.ipWhitelist = com.DefaultIpWhitelist
		}

		// 如果回调为 nil，则设置为空回调
		// If the callback is nil, set it to the empty callback
		if config.callback == nil {
			config.callback = &emptyCallback{}
		}
	} else {
		// 如果配置为 nil，则设置为默认配置
		// If the configuration is nil, set it to the default configuration
		config = DefaultConfig()
	}

	// 返回有效的配置
	// Return the valid configuration
	return config
}
