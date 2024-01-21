package internal

import (
	"context"
	"sync"
	"time"
)

const (
	// SegmentSize 表示 Segment 的大小，即 Segment 中最多可以存储的 key-value 对的数量
	// SegmentSize represents the size of Segment, that is, the maximum number of key-value pairs that can be stored in Segment
	SegmentSize = 1 << 8
	// SegmentMask 表示 Segment 的掩码，用于计算 key 对应的 Segment 的索引
	// SegmentMask represents the mask of Segment, which is used to calculate the index of Segment corresponding to key
	segmentAndOpVal = SegmentSize - 1
)

var (
	// defaultScanInterval 表示 Segment 的后台扫描任务的扫描间隔
	// defaultScanInterval represents the scan interval of the background scan task of Segment
	defaultScanInterval = 10 * time.Second
	// defaultExpireTime 表示 Segment 中 key-value 对的过期时间
	// defaultExpireTime represents the expiration time of key-value pairs in Segment
	defaultExpireTime = (defaultScanInterval).Milliseconds() * 3
)

// Segment 是一个分段结构，用于实现限流器的内部存储
// Segment is a segment structure used to implement internal storage of rate limiter
type Segment struct {
	data   map[string]any // key-value 对
	lock   sync.Mutex     // 用于保护 data 的读写
	once   sync.Once      // 用于保证 Stop 方法只执行一次
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewSegment 创建一个新的 Segment 实例
// NewSegment creates a new Segment instance
func NewSegment() *Segment {
	s := &Segment{
		data: make(map[string]any),
		lock: sync.Mutex{},
		once: sync.Once{},
		wg:   sync.WaitGroup{},
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// 启动后台扫描任务
	// Start the background scan task
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// 每隔 defaultScanInterval 扫描一次 Segment 中的所有 key-value 对
		// Scan all key-value pairs in Segment every defaultScanInterval
		ticker := time.NewTicker(defaultScanInterval)
		for {
			select {
			case <-s.ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				s.lock.Lock()
				now := time.Now().UnixMilli()
				// 遍历 Segment 中的所有 key-value 对，如果某个 key-value 对的过期时间已到，则将其删除
				// Traverse all key-value pairs in Segment, if the expiration time of a certain key-value pair has arrived, delete it
				for key, value := range s.data {
					if now-value.(*Element).GetUpdateAt() >= defaultExpireTime {
						// 重置 value 的值，放回到 ElementPool 中，以便复用
						// Reset the value of value, put it back into ElementPool for reuse
						value.(*Element).SetValue(nil)
						ElementPool.Put(value)
						// 从 Segment 中删除 key-value 对
						// Delete key-value pairs from Segment
						delete(s.data, key)
					}
				}
				s.lock.Unlock()
			}
		}
	}()

	return s
}

// Get 从 Segment 中获取指定 key 对应的值
// Get gets the value of the specified key from Segment
func (s *Segment) Get(key string) (any, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.data[key]
	return value, ok
}

// GetOrCreate 从 Segment 中获取指定 key 对应的值，如果不存在则使用给定的函数创建一个新值
// GetOrCreate gets the value of the specified key from Segment, and if it does not exist, use the given function to create a new value
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

// Set 将指定的 key-value 对存储到 Segment 中
// Set stores the specified key-value pair in Segment
func (s *Segment) Set(key string, value any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

// Delete 从 Segment 中删除指定的 key-value 对
// Delete deletes the specified key-value pair from Segment
func (s *Segment) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
}

// GetData 返回 Segment 中存储的所有 key-value 对
// GetData returns all key-value pairs stored in Segment
func (s *Segment) GetData() map[string]any {
	return s.data
}

// Stop 停止 Segment 的后台扫描任务
// Stop stops the background scan task of Segment
func (s *Segment) Stop() {
	s.once.Do(func() {
		s.cancel()
		s.wg.Wait()
	})
}
