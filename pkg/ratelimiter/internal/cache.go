package internal

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

type Cache struct {
	segments []*Segment
	once     sync.Once
}

func NewCache() *Cache {
	segments := make([]*Segment, SegmentSize)
	for i := 0; i < SegmentSize; i++ {
		segments[i] = NewSegment()
	}
	return &Cache{segments: segments, once: sync.Once{}}
}

func (c *Cache) Get(key string) (any, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Get(key)
}

func (c *Cache) GetOrCreate(key string, fn func() any) (any, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].GetOrCreate(key, fn)
}

func (c *Cache) Set(key string, value any) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Set(key, value)
}

func (c *Cache) Delete(key string) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Delete(key)
}

func (c *Cache) Segments() []*Segment {
	return c.segments
}

func (c *Cache) Stop() {
	c.once.Do(func() {
		for i := 0; i < SegmentSize; i++ {
			c.segments[i].Stop()
		}
	})
}
