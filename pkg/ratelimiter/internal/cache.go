package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cespare/xxhash/v2"
)

const (
	SegmentSize     = 1 << 8
	segmentAndOpVal = SegmentSize - 1
)

var defaultExpireTime = (30 * time.Second).Milliseconds()

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

type Segment struct {
	data   map[string]any
	lock   sync.Mutex
	once   sync.Once
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSegment() *Segment {
	s := &Segment{
		data: make(map[string]any),
		lock: sync.Mutex{},
		once: sync.Once{},
		wg:   sync.WaitGroup{},
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-s.ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				s.lock.Lock()
				now := time.Now().UnixMilli()
				for key, value := range s.data {
					if now-value.(*Element).GetUpdateAt() >= defaultExpireTime {
						value.(*Element).SetValue(nil)
						ElementPool.Put(value)
						delete(s.data, key)
					}
				}
				s.lock.Unlock()
			}
		}
	}()

	return s
}

func (s *Segment) Get(key string) (any, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Segment) GetOrCreate(key string, fn func() any) (any, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.data[key]
	if !ok {
		value = fn()
		s.data[key] = value
	}
	return value, ok
}

func (s *Segment) Set(key string, value any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

func (s *Segment) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
}

func (s *Segment) GetData() map[string]any {
	return s.data
}

func (s *Segment) Stop() {
	s.once.Do(func() {
		s.cancel()
		s.wg.Wait()
	})
}

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
