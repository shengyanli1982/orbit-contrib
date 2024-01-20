package internal

import (
	"context"
	"sync"
	"time"
)

const (
	SegmentSize     = 1 << 8
	segmentAndOpVal = SegmentSize - 1
)

var (
	defaultScanInterval = 10 * time.Second
	defaultExpireTime   = (defaultScanInterval).Milliseconds() * 3
)

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
		ticker := time.NewTicker(defaultScanInterval)
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
