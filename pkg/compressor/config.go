package compressor

import (
	"github.com/gin-gonic/gin"
	com "github.com/shengyanli1982/orbit-contrib/internal/common"
)

const (
	// 默认压缩等级
	// Default compression level
	DefaultBestCompression = 9 // 9 is the best compression level
	DefaultCompression     = 6 // 6 is the default compression level
	DefaultSpeed           = 3 // 3 is the fast speed compression level
	DefaultBestSpeed       = 1 // 1 is the best speed compression level
	DefaultNoCompression   = 0 // 0 is no compression
)

// WriterCreateFunc 是一个创建压缩写入器的函数
// WriterCreateFunc is a function to create a compression writer
type WriterCreateFunc func(config *Config, rw gin.ResponseWriter) any

// DefaultWriterCreateFunc 是一个默认的创建压缩写入器的函数
// DefaultWriterCreateFunc is a default function to create a compression writer
var DefaultWriterCreateFunc = func(config *Config, rw gin.ResponseWriter) any {
	// 默认使用 GZipWriter
	// Default to use GZipWriter
	return NewGZipWriter(config, rw)
}

// Config 是一个配置结构体
// Config is a struct of config
type Config struct {
	level       int
	ipWhitelist map[string]struct{}
	// 匹配函数
	// Match function
	matchFunc  com.HttpRequestHeaderMatchFunc
	createFunc WriterCreateFunc
}

// NewConfig 创建一个新的配置实例
// NewConfig creates a new config instance
func NewConfig() *Config {
	return &Config{
		level:       DefaultCompression,
		ipWhitelist: com.DefaultIpWhitelist,
		matchFunc:   com.DefaultLimitMatchFunc,
		createFunc:  DefaultWriterCreateFunc,
	}
}

// DefaultConfig 创建一个默认的配置实例
// DefaultConfig creates a default config instance
func DefaultConfig() *Config {
	return NewConfig()
}

// WithThreshold 设置压缩阈值
// WithThreshold sets the compression threshold
func (c *Config) WithCompressLevel(level int) *Config {
	c.level = level
	return c
}

// WithWriter 设置压缩写入器
// WithWriter sets the compression writer
func (c *Config) WithWriterCreateFunc(fn WriterCreateFunc) *Config {
	c.createFunc = fn
	return c
}

// WithMatchFunc 设置匹配函数
// WithMatchFunc sets the match function
func (c *Config) WithMatchFunc(match com.HttpRequestHeaderMatchFunc) *Config {
	c.matchFunc = match
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

// isConfigValid 检查配置是否有效
// isConfigValid checks whether the config is valid
func isConfigValid(config *Config) *Config {
	if config != nil {
		if config.level < 0 || config.level > DefaultBestCompression {
			config.level = DefaultCompression
		}
		if config.createFunc == nil {
			config.createFunc = DefaultWriterCreateFunc
		}
		if config.matchFunc == nil {
			config.matchFunc = com.DefaultLimitMatchFunc
		}
		if config.ipWhitelist == nil {
			config.ipWhitelist = com.DefaultIpWhitelist
		}
	} else {
		config = DefaultConfig()
	}
	return config
}
