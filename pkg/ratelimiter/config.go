package ratelimiter

import "net/http"

var empty = struct{}{}

var (
	DefaultLocalIpAddress     = "127.0.0.1"
	DefaultLocalIpv6Address   = "::1"
	DefaultLimitRatePerSecond = float64(1)
	DefaultLimitBurst         = 5
	DefaultLimitMatchFunc     = func(header *http.Request) bool { return true }
	DefaultIpWhitelist        = map[string]struct{}{
		DefaultLocalIpAddress:   empty,
		DefaultLocalIpv6Address: empty,
	}
)

type HttpRequestHeaderMatchFunc func(header *http.Request) bool

type Config struct {
	rate      float64
	burst     int
	whitelist map[string]struct{}
	match     HttpRequestHeaderMatchFunc
	callback  Callback
}

func NewConfig() *Config {
	return &Config{
		rate:      DefaultLimitRatePerSecond,
		burst:     DefaultLimitBurst,
		match:     DefaultLimitMatchFunc,
		whitelist: DefaultIpWhitelist,
		callback:  &emptyCallback{},
	}
}

func DefaultConfig() *Config {
	return NewConfig()
}

func (c *Config) WithCallback(callback Callback) *Config {
	c.callback = callback
	return c
}

func (c *Config) WithRate(rate float64) *Config {
	c.rate = rate
	return c
}

func (c *Config) WithBurst(burst int) *Config {
	c.burst = burst
	return c
}

func (c *Config) WithMatchFunc(match HttpRequestHeaderMatchFunc) *Config {
	c.match = match
	return c
}

func (c *Config) WithWhitelist(whitelist []string) *Config {
	for _, ip := range whitelist {
		c.whitelist[ip] = empty
	}
	return c
}

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
