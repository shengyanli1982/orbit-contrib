package rewriter

import (
	"net/url"

	com "github.com/shengyanli1982/orbit-contrib/internal/common"
)

// PathRewriteFunc 是一个路径重写函数
// PathRewriteFunc is a function to rewrite the path
type PathRewriteFunc func(header *url.URL) (bool, string)

// DefaultPathRewriteFunc 是一个默认的路径重写函数
// DefaultPathRewriteFunc is a default function to rewrite the path
var DefaultPathRewriteFunc = func(header *url.URL) (bool, string) {
	return false, ""
}

// Config 是一个配置结构体
// Config is a struct of config
type Config struct {
	// Ip 白名单
	// Ip whitelist
	ipWhitelist map[string]struct{}

	// 匹配函数
	// Match function
	matchFunc com.HttpRequestHeaderMatchFunc

	// 路径重写函数
	// Path rewrite function
	rewriteFunc PathRewriteFunc

	// 回调函数
	// Callback
	callback Callback
}

// NewConfig 创建一个新的配置实例
// NewConfig creates a new config instance
func NewConfig() *Config {
	return &Config{
		ipWhitelist: com.DefaultIpWhitelist,
		matchFunc:   com.DefaultLimitMatchFunc,
		rewriteFunc: DefaultPathRewriteFunc,
		callback:    &emptyCallback{},
	}
}

// DefaultConfig 创建一个默认的配置实例
// DefaultConfig creates a default config instance
func DefaultConfig() *Config {
	return NewConfig()
}

// WithCallback 设置回调函数
// WithCallback sets the callback function
func (c *Config) WithCallback(callback Callback) *Config {
	c.callback = callback
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

// WithPathRewriteFunc 设置路径重写函数
// WithPathRewriteFunc sets the path rewrite function
func (c *Config) WithPathRewriteFunc(fn PathRewriteFunc) *Config {
	c.rewriteFunc = fn
	return c
}

// isConfigValid 检查配置是否有效
// isConfigValid checks whether the config is valid
func isConfigValid(config *Config) *Config {
	if config != nil {
		if config.callback == nil {
			config.callback = &emptyCallback{}
		}
		if config.matchFunc == nil {
			config.matchFunc = com.DefaultLimitMatchFunc
		}
		if config.ipWhitelist == nil {
			config.ipWhitelist = com.DefaultIpWhitelist
		}
		if config.rewriteFunc == nil {
			config.rewriteFunc = DefaultPathRewriteFunc
		}
	} else {
		config = DefaultConfig()
	}
	return config
}
