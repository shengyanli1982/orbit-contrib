package internal

import (
	"sync"
	"sync/atomic"
	"time"
)

// Element 是一个结构体，包含一个任意类型的值和一个原子类型的更新时间
// Element is a struct that contains a value of any type and an update time of atomic type
type Element struct {
	// value 是一个任意类型的值
	// value is a value of any type
	value any

	// updateAt 是一个原子类型的更新时间，用于存储元素的最后更新时间
	// updateAt is an update time of atomic type, used to store the last update time of the element
	updateAt atomic.Int64
}

// NewElement 是一个函数，返回一个新的 Element 结构体的指针
// NewElement is a function that returns a new pointer to the Element struct
func NewElement() *Element {
	// 创建一个新的 Element 结构体的指针
	// Create a new pointer to the Element struct
	e := &Element{
		// 初始化 value 为 nil
		// Initialize value as nil
		value: nil,

		// 初始化 updateAt 为一个新的 atomic.Int64
		// Initialize updateAt as a new atomic.Int64
		updateAt: atomic.Int64{},
	}

	// 将 updateAt 的值设置为当前的 Unix 毫秒时间
	// Set the value of updateAt to the current Unix millisecond time
	e.updateAt.Store(time.Now().UnixMilli())

	// 返回新创建的 Element 结构体的指针
	// Return the newly created pointer to the Element struct
	return e
}

// SetValue 方法用于设置元素的值
// The SetValue method is used to set the value of the element
func (e *Element) SetValue(data any) {
	// 更新元素的更新时间为当前时间的毫秒数
	// Update the update time of the element to the current time in milliseconds
	e.updateAt.Store(time.Now().UnixMilli())

	// 设置元素的值
	// Set the value of the element
	e.value = data
}

// GetValue 方法用于获取元素的值
// The GetValue method is used to get the value of the element
func (e *Element) GetValue() any {
	// 更新元素的更新时间为当前时间的毫秒数
	// Update the update time of the element to the current time in milliseconds
	e.updateAt.Store(time.Now().UnixMilli())

	// 返回元素的值
	// Return the value of the element
	return e.value
}

// GetUpdateAt 方法用于获取元素的更新时间
// The GetUpdateAt method is used to get the update time of the element
func (e *Element) GetUpdateAt() int64 {
	// 返回元素的更新时间
	// Return the update time of the element
	return e.updateAt.Load()
}

// ElementPool 是一个同步池，用于存储元素
// ElementPool is a sync pool used to store elements
var ElementPool = sync.Pool{
	// New 方法用于创建一个新的元素
	// The New method is used to create a new element
	New: func() any {
		// 返回一个新的元素
		// Return a new element
		return NewElement()
	},
}
