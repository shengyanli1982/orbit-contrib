package common

import "fmt"

var (
	// 测试远程IP地址
	// Test remote IP address
	TestIpAddress  = "192.168.0.1" // 测试IP地址1
	TestIpAddress2 = "192.168.0.2" // 测试IP地址2
	TestIpAddress3 = "192.168.0.3" // 测试IP地址3

	// 测试端口
	// Test port
	TestPort = 13143 // 测试端口号

	// 测试URL路径
	// Test URL path
	TestUrlPath  = "/test"  // 测试URL路径1
	TestUrlPath2 = "/test2" // 测试URL路径2

	// 测试端点
	// Test endpoint
	TestEndpoint    = fmt.Sprintf("%s:%d", TestIpAddress, TestPort)         // 测试端点1
	TestEndpoint2   = fmt.Sprintf("%s:%d", TestIpAddress2, TestPort)        // 测试端点2
	TestEndpoint3   = fmt.Sprintf("%s:%d", TestIpAddress3, TestPort)        // 测试端点3
	DefaultEndpoint = fmt.Sprintf("%s:%d", DefaultLocalIpAddress, TestPort) // 默认本地IP地址的测试端点

	// 测试响应文本
	// Test response text
	TestResponseText = "This is HelloWorld!!" // 测试响应文本
)
