package main

import (
	"compress/flate"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	cr "github.com/shengyanli1982/orbit-contrib/pkg/compressor"
)

var (
	// 测试URL
	// Test URL
	testUrl = "/test"
)

// testNewDeflateWriterFunc 是一个创建新的DeflateWriter的函数
// testNewDeflateWriterFunc is a function to create a new DeflateWriter
func testNewDeflateWriterFunc(config *cr.Config, rw gin.ResponseWriter) any {
	return cr.NewDeflateWriter(config, rw)
}

// testRequestFunc 是一个测试请求的函数
// testRequestFunc is a function to test the request
func testRequestFunc(idx int, router *gin.Engine, conf *cr.Config, url string) {
	// 创建一个新的请求
	// Create a new request
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// 创建一个新的记录器
	// Create a new recorder
	resp := httptest.NewRecorder()

	// 服务HTTP
	// Serve HTTP
	router.ServeHTTP(resp, req)

	// 获取响应体内容
	// Get the body content of the response
	bodyContent := resp.Body.String()

	// 创建一个新的flate读取器
	// Create a new flate reader
	gr := flate.NewReader(resp.Body)
	defer gr.Close()

	// 读取所有的明文
	// Read all plaintext
	plaintext, _ := io.ReadAll(gr)

	// 打印请求信息
	// Print request information
	fmt.Println("[Request]", idx, resp.Code, url, string(plaintext), bodyContent)
}

func main() {
	// 创建新的配置
	// Create new configuration
	conf := cr.NewConfig().WithWriterCreateFunc(testNewDeflateWriterFunc)

	// 创建新的压缩器
	// Create new compressor
	compr := cr.NewCompressor(conf)

	// 停止压缩器
	// Stop the compressor
	defer compr.Stop()

	// 创建新的路由器
	// Create new router
	router := gin.New()

	// 使用压缩器的处理函数
	// Use the handler function of the compressor
	router.Use(compr.HandlerFunc())

	// 添加GET路由
	// Add GET route
	router.GET(testUrl, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 测试 10 次请求
	// Test 10 requests
	for i := 0; i < 10; i++ {
		// 调用测试请求函数，传入索引、路由器、配置和测试URL
		// Call the test request function, passing in the index, router, configuration, and test URL
		testRequestFunc(i, router, conf, testUrl)
	}

	// 等待所有请求任务执行完毕
	// Wait for all request tasks to complete
	time.Sleep(time.Second)
}
