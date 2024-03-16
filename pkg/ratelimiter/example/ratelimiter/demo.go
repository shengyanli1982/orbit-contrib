package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	rl "github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter"
)

var (
	// 测试URL
	// Test URL
	testUrl = "/test"
)

// testRequestFunc 是一个测试请求的函数
// testRequestFunc is a function to test the request
func testRequestFunc(idx int, router *gin.Engine, conf *rl.Config, url string) {
	// 创建一个新的请求
	// Create a new request
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// 创建一个新的响应记录器
	// Create a new response recorder
	resp := httptest.NewRecorder()

	// 使用路由器处理HTTP请求
	// Use the router to handle the HTTP request
	router.ServeHTTP(resp, req)

	// 打印请求的信息
	// Print the information of the request
	fmt.Println("[Request]", idx, resp.Code, url)
}

func main() {
	// 创建一个新的配置，设置速率为2，突发数为5
	// Create a new configuration, set the rate to 2 and the burst to 5
	conf := rl.NewConfig().WithRate(2).WithBurst(5)

	// 创建一个新的速率限制器
	// Create a new rate limiter
	limiter := rl.NewRateLimiter(conf)

	// 在函数返回时停止速率限制器
	// Stop the rate limiter when the function returns
	defer limiter.Stop()

	// 创建一个新的路由器
	// Create a new router
	router := gin.New()

	// 使用速率限制器的处理函数
	// Use the handler function of the rate limiter
	router.Use(limiter.HandlerFunc())

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
