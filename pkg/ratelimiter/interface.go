package ratelimiter

import "net/http"

// Callback 是一个限流回调接口，用于处理被限流的请求
// Callback is a rate limiting callback interface for processing requests that are rate limited
type Callback interface {
	// OnLimited 在请求被限流时被调用
	// OnLimited is called when the request is rate limited
	OnLimited(header *http.Request)
}

// emptyCallback 是一个空回调函数，不执行任何操作
// emptyCallback is an empty callback function that does nothing
type emptyCallback struct{}

// OnLimited 是空回调函数，不执行任何操作
// OnLimited is an empty callback function that does nothing
func (e *emptyCallback) OnLimited(header *http.Request) {}
