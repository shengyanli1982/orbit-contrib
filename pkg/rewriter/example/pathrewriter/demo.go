package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	rw "github.com/shengyanli1982/orbit-contrib/pkg/rewriter"
)

var (
	// 测试URL
	// Test URL
	testUrl = "/test"
)

// testRewriteFunc 是一个测试重写函数
// testRewriteFunc is a test rewrite function
func testRewriteFunc(u *url.URL) (bool, string) {
	// 如果URL的路径等于测试URL
	// If the path of the URL equals the test URL
	if u.Path == testUrl {
		// 返回真和新的URL
		// Return true and the new URL
		return true, testUrl + "2"
	}
	// 否则返回假和空字符串
	// Otherwise return false and an empty string
	return false, ""
}

// testRequestFunc 是一个测试请求函数
// testRequestFunc is a test request function
func testRequestFunc(idx int, router *gin.Engine, conf *rw.Config, url string) {
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
	fmt.Println("[Request]", idx, resp.Code, url, req.URL.Path, resp.Body.String())
}

func main() {
	// 创建一个新的配置
	// Create a new Config
	conf := rw.NewConfig().WithPathRewriteFunc(testRewriteFunc)

	// 创建一个新的路径重写器
	// Create a new Path Rewriter
	compr := rw.NewPathRewriter(conf)

	// 在函数返回时停止路径重写器
	// Stop the Path Rewriter when the function returns
	defer compr.Stop()

	// 创建一个新的Gin路由器
	// Create a new Gin router
	router := gin.New()

	// 使用路径重写器的处理函数
	// Use the handler function of the Path Rewriter
	router.Use(compr.HandlerFunc())

	// 添加一个新的路由
	// Add a new route
	router.GET(testUrl, func(c *gin.Context) {
		// 当请求成功时，返回"OK"
		// Return "OK" when the request is successful
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
