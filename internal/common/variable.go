package common

import "net/http"

// Empty 是一个空结构体，用于实现空的值
// Empty is an Empty struct for implementing Empty value
var Empty = struct{}{}

// HttpRequestHeaderMatchFunc 是一个匹配函数，用于匹配请求头
// HttpRequestHeaderMatchFunc is a match function for matching request headers
type HttpRequestHeaderMatchFunc func(header *http.Request) bool

var (
	// 默认本地IP地址
	// Default local IP address
	DefaultLocalIpAddress = "127.0.0.1"

	// 默认本地IPv6地址
	// Default local IPv6 address
	DefaultLocalIpv6Address = "::1"

	// 默认IP白名单
	// Default IP whitelist
	DefaultIpWhitelist = map[string]struct{}{}

	// 默认匹配函数
	// Default match function
	DefaultLimitMatchFunc = func(header *http.Request) bool { return true }
)