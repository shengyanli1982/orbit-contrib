package internal

import (
	"context"
	"sync"
	"time"
)

const (
	// SegmentSize 定义了 Segment 的大小
	// SegmentSize defines the size of the Segment
	SegmentSize = 1 << 8

	// segmentAndOpVal 是 SegmentSize 减 1 的结果，用于位运算
	// segmentAndOpVal is the result of SegmentSize minus 1, used for bitwise operations
	segmentAndOpVal = SegmentSize - 1
)

var (
	// defaultScanInterval 定义了默认的扫描间隔
	// defaultScanInterval defines the default scan interval
	defaultScanInterval = 10 * time.Second

	// defaultExpireTime 定义了默认的过期时间
	// defaultExpireTime defines the default expiration time
	defaultExpireTime = (defaultScanInterval).Milliseconds() * 3
)

// Segment 是一个结构体，包含数据、锁、上下文等字段
// Segment is a struct that contains fields such as data, lock, context, etc.
type Segment struct {
	// data 是一个 map，用于存储数据
	// data is a map used to store data
	data map[string]any

	// lock 是一个互斥锁，用于保护数据的并发访问
	// lock is a mutex used to protect concurrent access to data
	lock sync.Mutex

	// once 是一个 sync.Once 类型的变量，用于确保某些操作只执行一次
	// once is a variable of type sync.Once, used to ensure that certain operations are performed only once
	once sync.Once

	// wg 是一个等待组，用于等待 goroutine 完成
	// wg is a wait group used to wait for goroutines to finish
	wg sync.WaitGroup

	// ctx 是一个上下文，用于控制 goroutine 的生命周期
	// ctx is a context used to control the lifecycle of goroutines
	ctx context.Context

	// cancel 是一个函数，用于取消 ctx
	// cancel is a function used to cancel ctx
	cancel context.CancelFunc
}

// NewSegment 是一个函数，返回一个新的 Segment 结构体的指针
// NewSegment is a function that returns a new pointer to the Segment struct
func NewSegment() *Segment {
	// 创建一个新的 Segment 结构体的指针
	// Create a new pointer to the Segment struct
	s := &Segment{
		// 初始化 data 为一个新的 map
		// Initialize data as a new map
		data: make(map[string]any),

		// 初始化 lock 为一个新的 sync.Mutex
		// Initialize lock as a new sync.Mutex
		lock: sync.Mutex{},

		// 初始化 once 为一个新的 sync.Once
		// Initialize once as a new sync.Once
		once: sync.Once{},

		// 初始化 wg 为一个新的 sync.WaitGroup
		// Initialize wg as a new sync.WaitGroup
		wg: sync.WaitGroup{},
	}

	// 创建一个新的上下文和取消函数
	// Create a new context and cancel function
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// 增加等待组的计数
	// Increase the count of the wait group
	s.wg.Add(1)

	// 创建一个新的 goroutine
	// Create a new goroutine
	go func() {
		// 在 goroutine 结束时，减少等待组的计数
		// When the goroutine ends, decrease the count of the wait group
		defer s.wg.Done()

		// 创建一个新的定时器
		// Create a new timer
		ticker := time.NewTicker(defaultScanInterval)

		// 循环处理定时器的事件
		// Loop to handle the events of the timer
		for {
			select {
			// 如果上下文被取消，停止定时器并返回
			// If the context is cancelled, stop the timer and return
			case <-s.ctx.Done():
				// 停止定时器
				// Stop the timer
				ticker.Stop()

				return

			// 如果定时器触发，锁定数据并处理过期的元素
			// If the timer is triggered, lock the data and handle the expired elements
			case <-ticker.C:
				// 加锁，保护数据的并发访问
				// Lock to protect concurrent access to data
				s.lock.Lock()

				// 获取当前的 Unix 毫秒时间
				// Get the current Unix millisecond time
				now := time.Now().UnixMilli()

				// 遍历数据
				// Traverse the data
				for key, value := range s.data {
					// 如果元素的更新时间距离现在超过了默认的过期时间
					// If the update time of the element is more than the default expiration time from now
					if now-value.(*Element).GetUpdateAt() >= defaultExpireTime {
						// 将元素的值设置为 nil
						// Set the value of the element to nil
						value.(*Element).SetValue(nil)

						// 将元素放回到元素池
						// Put the element back into the element pool
						ElementPool.Put(value)

						// 从数据中删除这个元素
						// Delete this element from the data
						delete(s.data, key)
					}
				}

				// 解锁，释放数据的并发访问
				// Unlock to release concurrent access to data
				s.lock.Unlock()
			}
		}
	}()

	// 返回新创建的 Segment 结构体的指针
	// Return the newly created pointer to the Segment struct
	return s
}

// Get 方法用于获取指定键的值
// The Get method is used to get the value of the specified key
func (s *Segment) Get(key string) (any, bool) {
	// 加锁，防止并发操作
	// Lock to prevent concurrent operations
	s.lock.Lock()
	defer s.lock.Unlock()

	// 获取指定键的值
	// Get the value of the specified key
	value, ok := s.data[key]

	// 返回值和是否存在的标志
	// Return the value and the flag of whether it exists
	return value, ok
}

// GetOrCreate 方法用于获取指定键的值，如果不存在，则创建一个新的值
// The GetOrCreate method is used to get the value of the specified key, if it does not exist, create a new value
func (s *Segment) GetOrCreate(key string, fn func() any) (any, bool) {
	// 加锁，防止并发操作
	// Lock to prevent concurrent operations
	s.lock.Lock()
	defer s.lock.Unlock()

	// 获取指定键的值
	// Get the value of the specified key
	value, ok := s.data[key]

	// 如果值不存在
	// If the value does not exist
	if !ok {
		// 创建一个新的值
		// Create a new value
		value = fn()
		// 将新的值添加到数据中
		// Add the new value to the data
		s.data[key] = value
	}

	// 返回值和是否存在的标志
	// Return the value and the flag of whether it exists
	return value, ok
}

// Set 方法用于设置指定键的值
// The Set method is used to set the value of the specified key
func (s *Segment) Set(key string, value any) {
	// 加锁，防止并发操作
	// Lock to prevent concurrent operations
	s.lock.Lock()
	defer s.lock.Unlock()

	// 设置指定键的值
	// Set the value of the specified key
	s.data[key] = value
}

// Delete 方法用于删除指定键的值
// The Delete method is used to delete the value of the specified key
func (s *Segment) Delete(key string) {
	// 加锁，防止并发操作
	// Lock to prevent concurrent operations
	s.lock.Lock()
	defer s.lock.Unlock()

	// 删除指定键的值
	// Delete the value of the specified key
	delete(s.data, key)
}

// GetData 方法用于获取所有的数据
// The GetData method is used to get all data
func (s *Segment) GetData() map[string]any {
	// 返回所有的数据
	// Return all data
	return s.data
}

// Stop 方法用于停止段的操作，但在这里没有实现任何功能
// The Stop method is used to stop the operation of the segment, but it does not implement any functionality here
func (s *Segment) Stop() {
	s.once.Do(func() {
		s.cancel()
		s.wg.Wait()
	})
}
