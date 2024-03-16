package compressor

import (
	"github.com/gin-gonic/gin"
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
)

const (
	// DefaultBestCompression 是默认的最佳压缩等级，值为 9
	// DefaultBestCompression is the default best compression level, the value is 9
	DefaultBestCompression = 9

	// DefaultCompression 是默认的压缩等级，值为 6
	// DefaultCompression is the default compression level, the value is 6
	DefaultCompression = 6

	// DefaultSpeed 是默认的快速压缩等级，值为 3
	// DefaultSpeed is the default fast speed compression level, the value is 3
	DefaultSpeed = 3

	// DefaultBestSpeed 是默认的最佳速度压缩等级，值为 1
	// DefaultBestSpeed is the default best speed compression level, the value is 1
	DefaultBestSpeed = 1

	// DefaultNoCompression 是没有压缩，值为 0
	// DefaultNoCompression is no compression, the value is 0
	DefaultNoCompression = 0
)

// WriterCreateFunc 是一个创建压缩写入器的函数类型
// WriterCreateFunc is a function type to create a compression writer
type WriterCreateFunc func(config *Config, rw gin.ResponseWriter) any

// DefaultWriterCreateFunc 是一个默认的创建压缩写入器的函数，它返回一个 GZipWriter 实例
// DefaultWriterCreateFunc is a default function to create a compression writer, it returns a GZipWriter instance
var DefaultWriterCreateFunc = func(config *Config, rw gin.ResponseWriter) any {
	// 默认使用 GZipWriter
	// Default to use GZipWriter
	return NewGZipWriter(config, rw)
}

// Config 是一个配置结构体，包含压缩等级、IP白名单、匹配函数和创建压缩写入器的函数
// Config is a struct of config, including compression level, IP whitelist, match function and function to create a compression writer
type Config struct {
	// 压缩等级
	// Compression level
	level int

	// Ip 白名单
	// Ip whitelist
	ipWhitelist map[string]struct{}

	// 匹配函数，用于匹配 HTTP 请求头
	// Match function, used to match HTTP request headers
	matchFunc com.HttpRequestHeaderMatchFunc

	// 创建压缩写入器的函数
	// Function to create a compression writer
	createFunc WriterCreateFunc
}

// NewConfig 创建一个新的配置实例，包括默认的压缩等级、IP白名单、匹配函数和创建压缩写入器的函数
// NewConfig creates a new config instance, including default compression level, IP whitelist, match function and function to create a compression writer
func NewConfig() *Config {
	// 返回一个新的配置实例
	// Returns a new config instance
	return &Config{
		// 设置默认的压缩等级
		// Sets the default compression level
		level: DefaultCompression,

		// 设置默认的 IP 白名单
		// Sets the default IP whitelist
		ipWhitelist: com.DefaultIpWhitelist,

		// 设置默认的匹配函数
		// Sets the default match function
		matchFunc: com.DefaultLimitMatchFunc,

		// 设置默认的创建压缩写入器的函数
		// Sets the default function to create a compression writer
		createFunc: DefaultWriterCreateFunc,
	}
}

// DefaultConfig 创建一个默认的配置实例，实际上就是调用 NewConfig 函数
// DefaultConfig creates a default config instance, which is actually calling the NewConfig function
func DefaultConfig() *Config {
	return NewConfig()
}

// WithCompressLevel 设置压缩等级，并返回配置实例
// WithCompressLevel sets the compression level and returns the config instance
func (c *Config) WithCompressLevel(level int) *Config {
	c.level = level
	return c
}

// WithWriterCreateFunc 设置创建压缩写入器的函数，并返回配置实例
// WithWriterCreateFunc sets the function to create a compression writer and returns the config instance
func (c *Config) WithWriterCreateFunc(fn WriterCreateFunc) *Config {
	c.createFunc = fn
	return c
}

// WithMatchFunc 设置匹配函数，并返回配置实例
// WithMatchFunc sets the match function and returns the config instance
func (c *Config) WithMatchFunc(fn com.HttpRequestHeaderMatchFunc) *Config {
	c.matchFunc = fn
	return c
}

// WithIpWhitelist 设置 IP 白名单，并返回配置实例
// WithIpWhitelist sets the IP whitelist and returns the config instance
func (c *Config) WithIpWhitelist(whitelist []string) *Config {
	for _, ip := range whitelist {
		c.ipWhitelist[ip] = com.Empty
	}
	return c
}

// isConfigValid 检查配置是否有效
// isConfigValid checks whether the config is valid
func isConfigValid(config *Config) *Config {
	// 如果配置不为空
	// If the config is not null
	if config != nil {
		// 如果压缩等级小于 0 或者大于默认的最佳压缩等级
		// If the compression level is less than 0 or greater than the default best compression level
		if config.level < 0 || config.level > DefaultBestCompression {
			// 设置压缩等级为默认的压缩等级
			// Sets the compression level to the default compression level
			config.level = DefaultCompression
		}

		// 如果创建压缩写入器的函数为空
		// If the function to create a compression writer is null
		if config.createFunc == nil {
			// 设置创建压缩写入器的函数为默认的函数
			// Sets the function to create a compression writer to the default function
			config.createFunc = DefaultWriterCreateFunc
		}

		// 如果匹配函数为空
		// If the match function is null
		if config.matchFunc == nil {
			// 设置匹配函数为默认的函数
			// Sets the match function to the default function
			config.matchFunc = com.DefaultLimitMatchFunc
		}

		// 如果 IP 白名单为空
		// If the IP whitelist is null
		if config.ipWhitelist == nil {
			// 设置 IP 白名单为默认的白名单
			// Sets the IP whitelist to the default whitelist
			config.ipWhitelist = com.DefaultIpWhitelist
		}
	} else {
		// 如果配置为空，设置配置为默认的配置
		// If the config is null, sets the config to the default config
		config = DefaultConfig()
	}

	// 返回配置
	// Returns the config
	return config
}
