package ratelimiter

import "net/http"

type Callback interface {
	OnLimited(header *http.Request)
}

type emptyCallback struct{}

func (e *emptyCallback) OnLimited(header *http.Request) {}
