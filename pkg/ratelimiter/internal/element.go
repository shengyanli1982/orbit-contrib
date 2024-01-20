package internal

import (
	"sync"
	"sync/atomic"
	"time"
)

type Element struct {
	value    any
	updateAt atomic.Int64
}

func NewElement() *Element {
	e := &Element{
		value:    nil,
		updateAt: atomic.Int64{},
	}
	e.updateAt.Store(time.Now().UnixMilli())
	return e
}

func (e *Element) SetValue(data any) {
	e.updateAt.Store(time.Now().UnixMilli())
	e.value = data
}

func (e *Element) GetValue() any {
	e.updateAt.Store(time.Now().UnixMilli())
	return e.value
}

func (e *Element) GetUpdateAt() int64 {
	return e.updateAt.Load()
}

var ElementPool = sync.Pool{
	New: func() any {
		return NewElement()
	},
}
