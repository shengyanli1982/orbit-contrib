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
	// 默认不进行重写，返回假和空字符串
	// By default, do not rewrite, return false and an empty string
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
		// 默认的IP白名单
		// Default IP whitelist
		ipWhitelist: com.DefaultIpWhitelist,

		// 默认的限制匹配函数
		// Default limit match function
		matchFunc: com.DefaultLimitMatchFunc,

		// 默认的路径重写函数
		// Default path rewrite function
		rewriteFunc: DefaultPathRewriteFunc,

		// 空的回调函数
		// Empty callback function
		callback: &emptyCallback{},
	}
}

// DefaultConfig 创建一个默认的配置实例
// DefaultConfig creates a default config instance
func DefaultConfig() *Config {
	// 调用 NewConfig 函数创建一个新的配置实例
	// Call the NewConfig function to create a new config instance
	return NewConfig()
}

// WithCallback 设置回调函数
// WithCallback sets the callback function
func (c *Config) WithCallback(callback Callback) *Config {
	// 设置回调函数
	// Set the callback function
	c.callback = callback

	// 返回配置实例
	// Return the config instance
	return c
}

// WithMatchFunc 设置匹配函数
// WithMatchFunc sets the match function
func (c *Config) WithMatchFunc(fn com.HttpRequestHeaderMatchFunc) *Config {
	// 设置匹配函数
	// Set the match function
	c.matchFunc = fn

	// 返回配置实例
	// Return the config instance
	return c
}

// WithIpWhitelist 设置白名单
// WithIpWhitelist sets the whitelist
func (c *Config) WithIpWhitelist(whitelist []string) *Config {
	// 遍历白名单列表
	// Iterate over the whitelist
	for _, ip := range whitelist {
		// 将 IP 添加到白名单
		// Add the IP to the whitelist
		c.ipWhitelist[ip] = com.Empty
	}

	// 返回配置实例
	// Return the config instance
	return c
}

// WithPathRewriteFunc 设置路径重写函数
// WithPathRewriteFunc sets the path rewrite function
func (c *Config) WithPathRewriteFunc(fn PathRewriteFunc) *Config {
	// 设置路径重写函数
	// Set the path rewrite function
	c.rewriteFunc = fn

	// 返回配置实例
	// Return the config instance
	return c
}

// isConfigValid 检查配置是否有效
// isConfigValid checks whether the config is valid
func isConfigValid(config *Config) *Config {
	// 如果配置不为空
	// If the config is not null
	if config != nil {
		// 如果回调函数为空，设置为默认的空回调函数
		// If the callback function is null, set it to the default empty callback function
		if config.callback == nil {
			config.callback = &emptyCallback{}
		}

		// 如果匹配函数为空，设置为默认的限制匹配函数
		// If the match function is null, set it to the default limit match function
		if config.matchFunc == nil {
			config.matchFunc = com.DefaultLimitMatchFunc
		}

		// 如果 IP 白名单为空，设置为默认的 IP 白名单
		// If the IP whitelist is null, set it to the default IP whitelist
		if config.ipWhitelist == nil {
			config.ipWhitelist = com.DefaultIpWhitelist
		}

		// 如果路径重写函数为空，设置为默认的路径重写函数
		// If the path rewrite function is null, set it to the default path rewrite function
		if config.rewriteFunc == nil {
			config.rewriteFunc = DefaultPathRewriteFunc
		}
	} else {
		// 如果配置为空，设置为默认配置
		// If the config is null, set it to the default config
		config = DefaultConfig()
	}

	// 返回配置
	// Return the config
	return config
}
