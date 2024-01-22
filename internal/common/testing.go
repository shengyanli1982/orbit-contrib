package common

import "fmt"

var (
	// 测试远程IP地址
	// Test remote IP address
	TestIpAddress  = "192.168.0.1"
	TestIpAddress2 = "192.168.0.2"
	TestIpAddress3 = "192.168.0.3"
	// 测试端口
	// Test port
	TestPort = 13143
	// 测试URL路径
	// Test URL path
	TestUrlPath  = "/test"
	TestUrlPath2 = "/test2"
	// 测试端点
	// Test endpoint
	TestEndpoint    = fmt.Sprintf("%s:%d", TestIpAddress, TestPort)
	TestEndpoint2   = fmt.Sprintf("%s:%d", TestIpAddress2, TestPort)
	TestEndpoint3   = fmt.Sprintf("%s:%d", TestIpAddress3, TestPort)
	DefaultEndpoint = fmt.Sprintf("%s:%d", DefaultLocalIpAddress, TestPort)
)
