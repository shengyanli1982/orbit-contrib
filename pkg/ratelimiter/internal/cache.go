package internal

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

const (
	SegmentSize     = 1 << 8
	segmentAndOpVal = SegmentSize - 1
)

type Segment struct {
	data map[string]interface{}
	lock sync.Mutex
}

func NewSegment() *Segment {
	return &Segment{
		data: make(map[string]interface{}),
		lock: sync.Mutex{},
	}
}

func (s *Segment) Get(key string) (interface{}, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Segment) GetOrCreate(key string, fn func() interface{}) (interface{}, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.data[key]
	if !ok {
		value = fn()
		s.data[key] = value
	}
	return value, ok
}

func (s *Segment) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

func (s *Segment) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
}

func (s *Segment) GetData() map[string]interface{} {
	return s.data
}

type Cache struct {
	segments []*Segment
}

func NewCache() *Cache {
	segments := make([]*Segment, SegmentSize)
	for i := 0; i < SegmentSize; i++ {
		segments[i] = NewSegment()
	}
	return &Cache{segments: segments}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Get(key)
}

func (c *Cache) GetOrCreate(key string, fn func() interface{}) (interface{}, bool) {
	return c.segments[xxhash.Sum64String(key)&segmentAndOpVal].GetOrCreate(key, fn)
}

func (c *Cache) Set(key string, value interface{}) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Set(key, value)
}

func (c *Cache) Delete(key string) {
	c.segments[xxhash.Sum64String(key)&segmentAndOpVal].Delete(key)
}

func (c *Cache) Segments() []*Segment {
	return c.segments
}
