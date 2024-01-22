package compressor

import (
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

	// 默认压缩阈值, 如果内容长度小于1KB, 不压缩
	// 1KB, if content length is less than 1KB, do not compress
	DefaultThreshold = 1024
)

// WriterCreateFunc 是一个创建压缩写入器的函数
// WriterCreateFunc is a function to create a compression writer
type WriterCreateFunc func(config *Config) any

// DefaultWriterCreateFunc 是一个默认的创建压缩写入器的函数
// DefaultWriterCreateFunc is a default function to create a compression writer
var DefaultWriterCreateFunc = func(config *Config) any { return nil }

// Config 是一个配置结构体
// Config is a struct of config
type Config struct {
	level       int
	threshold   int
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
		threshold:   DefaultThreshold,
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

// WithIpWhitelist 设置白名单
// WithIpWhitelist sets the whitelist
func (c *Config) WithIpWhitelist(whitelist []string) *Config {
	for _, ip := range whitelist {
		c.ipWhitelist[ip] = com.Empty
	}
	return c
}

// WithThreshold 设置内容大小压缩阈值 (单位: 字节)
// WithThreshold sets the compression threshold for content size (in bytes)
func (c *Config) WithThreshold(threshold int) *Config {
	c.threshold = threshold
	return c
}

// WithWriter 设置压缩写入器
// WithWriter sets the compression writer
func (c *Config) WithWriterCreateFunc(fn WriterCreateFunc) *Config {
	c.createFunc = fn
	return c
}

// isConfigValid 检查配置是否有效
// isConfigValid checks whether the config is valid
func isConfigValid(config *Config) *Config {
	if config != nil {
		if config.level < 0 || config.level > DefaultBestCompression {
			config.level = DefaultCompression
		}
		if config.threshold <= DefaultThreshold {
			config.threshold = DefaultThreshold
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
