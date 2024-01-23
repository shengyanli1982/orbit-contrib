package internal

import (
	"sync"
	"sync/atomic"
	"time"
)

// Element 是 Segment 中的元素，用于存储 key-value 对
// Element is an element in Segment, used to store key-value pairs
type Element struct {
	value    any
	updateAt atomic.Int64
}

// NewElement 返回一个新的 Element 实例。
// NewElement returns a new Element instance.
func NewElement() *Element {
	e := &Element{
		value:    nil,
		updateAt: atomic.Int64{},
	}
	// 初始化 updateAt 的值为当前时间戳
	// Initialize the value of updateAt to the current timestamp
	e.updateAt.Store(time.Now().UnixMilli())
	return e
}

// SetValue 设置 Element 的值
// SetValue sets the value of Element
func (e *Element) SetValue(data any) {
	e.updateAt.Store(time.Now().UnixMilli())
	e.value = data
}

// GetValue 获取 Element 的值
// GetValue gets the value of Element
func (e *Element) GetValue() any {
	e.updateAt.Store(time.Now().UnixMilli())
	return e.value
}

// GetUpdateAt 获取 Element 的更新时间
// GetUpdateAt gets the update time of Element
func (e *Element) GetUpdateAt() int64 {
	return e.updateAt.Load()
}

// ElementPool 是 Element 的对象池，用于复用 Element 实例
// ElementPool is the object pool of Element, used to reuse Element instances
var ElementPool = sync.Pool{
	New: func() any {
		return NewElement()
	},
}
