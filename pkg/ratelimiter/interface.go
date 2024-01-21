package ratelimiter

import "net/http"

// Callback 是一个限流回调接口，用于处理被限流的请求
type Callback interface {
	// OnLimited 在请求被限流时被调用
	OnLimited(header *http.Request)
}

// emptyCallback 是一个空回调函数，不执行任何操作
type emptyCallback struct{}

// OnLimited 是空回调函数，不执行任何操作
func (e *emptyCallback) OnLimited(header *http.Request) {}
