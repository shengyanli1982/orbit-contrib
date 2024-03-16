package ratelimiter

import "net/http"

// Callback 是一个限流回调接口，用于处理被限流的请求
// Callback is a rate limiting callback interface for handling requests that are rate limited
type Callback interface {
	// OnLimited 是一个方法，当请求被限流时被调用，接收一个 http.Request 指针作为参数
	// OnLimited is a method that is called when a request is rate limited, it takes a pointer to http.Request as a parameter
	OnLimited(header *http.Request)
}

// emptyCallback 是一个实现了 Callback 接口的结构体，它的方法不执行任何操作
// emptyCallback is a struct that implements the Callback interface, its methods do not perform any operations
type emptyCallback struct{}

// OnLimited 是 emptyCallback 结构体的方法，它不执行任何操作，接收一个 http.Request 指针作为参数
// OnLimited is a method of the emptyCallback struct, it does not perform any operations, it takes a pointer to http.Request as a parameter
func (e *emptyCallback) OnLimited(header *http.Request) {}
