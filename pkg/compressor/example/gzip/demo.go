package main

import (
	"compress/gzip"
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

// testRequestFunc 是一个测试请求的函数
// testRequestFunc is a function to test the request
func testRequestFunc(idx int, router *gin.Engine, conf *cr.Config, url string) {
	// 创建一个新的请求
	// Create a new request
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// 创建一个新的响应记录器
	// Create a new response recorder
	resp := httptest.NewRecorder()

	// 使用路由器处理HTTP请求
	// Use the router to handle the HTTP request
	router.ServeHTTP(resp, req)

	// 获取响应体内容
	// Get the content of the response body
	bodyContent := resp.Body.String()

	// 创建一个新的gzip读取器
	// Create a new gzip reader
	gr, _ := gzip.NewReader(resp.Body)
	defer gr.Close()

	// 读取响应的全部内容
	// Read all the content of the response
	plaintext, _ := io.ReadAll(gr)

	// 打印请求的信息
	// Print the information of the request
	fmt.Println("[Request]", idx, resp.Code, url, string(plaintext), bodyContent)
}

func main() {
	// 创建一个新的配置
	// Create a new configuration
	conf := cr.NewConfig()

	// 创建一个新的压缩器
	// Create a new compressor
	compr := cr.NewCompressor(conf)

	// 在函数返回时停止压缩器
	// Stop the compressor when the function returns
	defer compr.Stop()

	// 创建一个新的路由器
	// Create a new router
	router := gin.New()

	// 使用压缩器的处理函数
	// Use the handler function of the compressor
	router.Use(compr.HandlerFunc())

	// 添加一个GET路由
	// Add a GET route
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
