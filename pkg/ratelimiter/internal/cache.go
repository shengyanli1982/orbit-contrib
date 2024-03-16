package internal

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

// Cache 是一个结构体，包含多个 Segment 和一个 sync.Once 类型的变量
// Cache is a struct that contains multiple Segments and a variable of type sync.Once
type Cache struct {
	// segments 是一个 Segment 类型的切片，用于存储多个 Segment
	// segments is a slice of type Segment, used to store multiple Segments
	segments []*Segment

	// once 是一个 sync.Once 类型的变量，用于确保某些操作只执行一次
	// once is a variable of type sync.Once, used to ensure that certain operations are performed only once
	once sync.Once
}

// NewCache 是一个函数，返回一个新的 Cache 结构体的指针
// NewCache is a function that returns a new pointer to the Cache struct
func NewCache() *Cache {
	// 创建一个 Segment 类型的切片，长度为 SegmentSize
	// Create a slice of type Segment with a length of SegmentSize
	segments := make([]*Segment, SegmentSize)

	// 遍历切片
	// Traverse the slice
	for i := 0; i < SegmentSize; i++ {
		// 将切片的每个元素初始化为一个新的 Segment
		// Initialize each element of the slice as a new Segment
		segments[i] = NewSegment()
	}

	// 返回一个新的 Cache 结构体的指针，其中 segments 为刚刚创建的切片，once 为一个新的 sync.Once
	// Return a new pointer to the Cache struct, where segments is the slice just created and once is a new sync.Once
	return &Cache{segments: segments, once: sync.Once{}}
}

// Get 方法用于从缓存中获取指定键的值
// The Get method is used to get the value of the specified key from the cache
func (c *Cache) Get(key string) (any, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Get(key)
}

// GetOrCreate 方法用于从缓存中获取指定键的值，如果键不存在，则使用指定的函数创建一个新的值
// The GetOrCreate method is used to get the value of the specified key from the cache, if the key does not exist, a new value is created using the specified function
func (c *Cache) GetOrCreate(key string, fn func() any) (any, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].GetOrCreate(key, fn)
}

// Set 方法用于在缓存中设置指定键的值
// The Set method is used to set the value of the specified key in the cache
func (c *Cache) Set(key string, value any) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Set(key, value)
}

// Delete 方法用于从缓存中删除指定键的值
// The Delete method is used to delete the value of the specified key from the cache
func (c *Cache) Delete(key string) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Delete(key)
}

// Segments 方法用于获取缓存的所有段
// The Segments method is used to get all segments of the cache
func (c *Cache) Segments() []*Segment {
	return c.segments
}

// Stop 方法用于停止缓存，它会停止缓存的所有段
// The Stop method is used to stop the cache, it will stop all segments of the cache
func (c *Cache) Stop() {
	c.once.Do(func() {
		for i := 0; i < SegmentSize; i++ {
			c.segments[i].Stop()
		}
	})
}
