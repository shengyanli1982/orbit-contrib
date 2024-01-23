package internal

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

// Cache 是一个缓存结构体
// Cache is a struct of cache
type Cache struct {
	segments []*Segment
	once     sync.Once
}

// NewCache 创建一个新的缓存实例
// NewCache creates a new cache instance
func NewCache() *Cache {
	// 创建 SegmentSize 个 Segment 实例
	// Create SegmentSize Segment instances
	segments := make([]*Segment, SegmentSize)
	for i := 0; i < SegmentSize; i++ {
		segments[i] = NewSegment()
	}
	// 返回 Cache 实例
	// Return Cache instance
	return &Cache{segments: segments, once: sync.Once{}}
}

// Get 从缓存中获取指定键的值
// Get gets the value of the specified key from the cache
func (c *Cache) Get(key string) (any, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Get(key)
}

// GetOrCreate 从缓存中获取指定键的值，如果键不存在则创建并返回默认值
// GetOrCreate gets the value of the specified key from the cache, and creates and returns the default value if the key does not exist
func (c *Cache) GetOrCreate(key string, fn func() any) (any, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].GetOrCreate(key, fn)
}

// Set 向缓存中设置指定键的值
// Set sets the value of the specified key in the cache
func (c *Cache) Set(key string, value any) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Set(key, value)
}

// Delete 从缓存中删除指定键的值
// Delete deletes the value of the specified key from the cache
func (c *Cache) Delete(key string) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Delete(key)
}

// Segments 返回缓存的所有分段
// Segments returns all segments of the cache
func (c *Cache) Segments() []*Segment {
	return c.segments
}

// Stop 停止缓存的所有分段
// Stop stops all segments of the cache
func (c *Cache) Stop() {
	c.once.Do(func() {
		// 停止所有分段的后台扫描任务
		// Stop the background scan task of all segments
		for i := 0; i < SegmentSize; i++ {
			c.segments[i].Stop()
		}
	})
}
