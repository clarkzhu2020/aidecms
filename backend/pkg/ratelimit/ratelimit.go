package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter 限流器接口
type Limiter interface {
	// Allow 检查是否允许请求
	Allow(key string) bool
	// AllowN 检查是否允许 n 个请求
	AllowN(key string, n int) bool
	// Reset 重置指定键的限制
	Reset(key string)
}

// TokenBucket 令牌桶算法实现
type TokenBucket struct {
	rate       int // 每秒生成的令牌数
	capacity   int // 桶容量
	buckets    map[string]*bucket
	mu         sync.RWMutex
	gcInterval time.Duration // 垃圾回收间隔
	ctx        context.Context
	cancel     context.CancelFunc
}

type bucket struct {
	tokens    float64
	lastCheck time.Time
	mu        sync.Mutex
}

// NewTokenBucket 创建令牌桶限流器
func NewTokenBucket(rate, capacity int) *TokenBucket {
	ctx, cancel := context.WithCancel(context.Background())
	tb := &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		buckets:    make(map[string]*bucket),
		gcInterval: 5 * time.Minute,
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动垃圾回收
	go tb.gc()

	return tb
}

// Allow 检查是否允许请求
func (tb *TokenBucket) Allow(key string) bool {
	return tb.AllowN(key, 1)
}

// AllowN 检查是否允许 n 个请求
func (tb *TokenBucket) AllowN(key string, n int) bool {
	tb.mu.RLock()
	b, exists := tb.buckets[key]
	tb.mu.RUnlock()

	if !exists {
		tb.mu.Lock()
		// 双重检查
		if b, exists = tb.buckets[key]; !exists {
			b = &bucket{
				tokens:    float64(tb.capacity),
				lastCheck: time.Now(),
			}
			tb.buckets[key] = b
		}
		tb.mu.Unlock()
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// 计算应该添加的令牌数
	now := time.Now()
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * float64(tb.rate)

	// 限制令牌数不超过容量
	if b.tokens > float64(tb.capacity) {
		b.tokens = float64(tb.capacity)
	}

	b.lastCheck = now

	// 检查是否有足够的令牌
	if b.tokens >= float64(n) {
		b.tokens -= float64(n)
		return true
	}

	return false
}

// Reset 重置指定键的限制
func (tb *TokenBucket) Reset(key string) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	delete(tb.buckets, key)
}

// Close 关闭限流器
func (tb *TokenBucket) Close() {
	tb.cancel()
}

// gc 垃圾回收
func (tb *TokenBucket) gc() {
	ticker := time.NewTicker(tb.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tb.ctx.Done():
			return
		case <-ticker.C:
			tb.mu.Lock()
			now := time.Now()
			for key, b := range tb.buckets {
				b.mu.Lock()
				// 如果桶超过 10 分钟没有使用，删除它
				if now.Sub(b.lastCheck) > 10*time.Minute {
					delete(tb.buckets, key)
				}
				b.mu.Unlock()
			}
			tb.mu.Unlock()
		}
	}
}

// SlidingWindow 滑动窗口算法实现
type SlidingWindow struct {
	limit      int           // 时间窗口内的最大请求数
	window     time.Duration // 时间窗口大小
	windows    map[string]*windowData
	mu         sync.RWMutex
	gcInterval time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
}

type windowData struct {
	requests []time.Time
	mu       sync.Mutex
}

// NewSlidingWindow 创建滑动窗口限流器
func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	ctx, cancel := context.WithCancel(context.Background())
	sw := &SlidingWindow{
		limit:      limit,
		window:     window,
		windows:    make(map[string]*windowData),
		gcInterval: 5 * time.Minute,
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动垃圾回收
	go sw.gc()

	return sw
}

// Allow 检查是否允许请求
func (sw *SlidingWindow) Allow(key string) bool {
	return sw.AllowN(key, 1)
}

// AllowN 检查是否允许 n 个请求
func (sw *SlidingWindow) AllowN(key string, n int) bool {
	sw.mu.RLock()
	wd, exists := sw.windows[key]
	sw.mu.RUnlock()

	if !exists {
		sw.mu.Lock()
		// 双重检查
		if wd, exists = sw.windows[key]; !exists {
			wd = &windowData{
				requests: make([]time.Time, 0),
			}
			sw.windows[key] = wd
		}
		sw.mu.Unlock()
	}

	wd.mu.Lock()
	defer wd.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	// 移除过期的请求
	validRequests := make([]time.Time, 0)
	for _, reqTime := range wd.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	wd.requests = validRequests

	// 检查是否超过限制
	if len(wd.requests)+n <= sw.limit {
		for i := 0; i < n; i++ {
			wd.requests = append(wd.requests, now)
		}
		return true
	}

	return false
}

// Reset 重置指定键的限制
func (sw *SlidingWindow) Reset(key string) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	delete(sw.windows, key)
}

// Close 关闭限流器
func (sw *SlidingWindow) Close() {
	sw.cancel()
}

// gc 垃圾回收
func (sw *SlidingWindow) gc() {
	ticker := time.NewTicker(sw.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sw.ctx.Done():
			return
		case <-ticker.C:
			sw.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-sw.window * 2) // 保留2倍窗口时间

			for key, wd := range sw.windows {
				wd.mu.Lock()
				if len(wd.requests) == 0 || wd.requests[len(wd.requests)-1].Before(cutoff) {
					delete(sw.windows, key)
				}
				wd.mu.Unlock()
			}
			sw.mu.Unlock()
		}
	}
}

// GetStats 获取统计信息
func (sw *SlidingWindow) GetStats(key string) map[string]interface{} {
	sw.mu.RLock()
	wd, exists := sw.windows[key]
	sw.mu.RUnlock()

	if !exists {
		return map[string]interface{}{
			"requests":  0,
			"limit":     sw.limit,
			"window":    sw.window.String(),
			"remaining": sw.limit,
		}
	}

	wd.mu.Lock()
	defer wd.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	// 计算有效请求数
	count := 0
	for _, reqTime := range wd.requests {
		if reqTime.After(cutoff) {
			count++
		}
	}

	return map[string]interface{}{
		"requests":  count,
		"limit":     sw.limit,
		"window":    sw.window.String(),
		"remaining": sw.limit - count,
	}
}

// FixedWindow 固定窗口算法实现
type FixedWindow struct {
	limit   int
	window  time.Duration
	windows map[string]*fixedWindowData
	mu      sync.RWMutex
}

type fixedWindowData struct {
	count     int
	resetTime time.Time
	mu        sync.Mutex
}

// NewFixedWindow 创建固定窗口限流器
func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		limit:   limit,
		window:  window,
		windows: make(map[string]*fixedWindowData),
	}
}

// Allow 检查是否允许请求
func (fw *FixedWindow) Allow(key string) bool {
	return fw.AllowN(key, 1)
}

// AllowN 检查是否允许 n 个请求
func (fw *FixedWindow) AllowN(key string, n int) bool {
	fw.mu.RLock()
	fwd, exists := fw.windows[key]
	fw.mu.RUnlock()

	if !exists {
		fw.mu.Lock()
		if fwd, exists = fw.windows[key]; !exists {
			fwd = &fixedWindowData{
				count:     0,
				resetTime: time.Now().Add(fw.window),
			}
			fw.windows[key] = fwd
		}
		fw.mu.Unlock()
	}

	fwd.mu.Lock()
	defer fwd.mu.Unlock()

	now := time.Now()

	// 检查是否需要重置窗口
	if now.After(fwd.resetTime) {
		fwd.count = 0
		fwd.resetTime = now.Add(fw.window)
	}

	// 检查是否超过限制
	if fwd.count+n <= fw.limit {
		fwd.count += n
		return true
	}

	return false
}

// Reset 重置指定键的限制
func (fw *FixedWindow) Reset(key string) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	delete(fw.windows, key)
}

// GetResetTime 获取重置时间
func (fw *FixedWindow) GetResetTime(key string) time.Time {
	fw.mu.RLock()
	fwd, exists := fw.windows[key]
	fw.mu.RUnlock()

	if !exists {
		return time.Now().Add(fw.window)
	}

	fwd.mu.Lock()
	defer fwd.mu.Unlock()
	return fwd.resetTime
}

// LimiterFactory 限流器工厂
type LimiterFactory struct{}

// CreateTokenBucket 创建令牌桶限流器
func (f *LimiterFactory) CreateTokenBucket(rate, capacity int) Limiter {
	return NewTokenBucket(rate, capacity)
}

// CreateSlidingWindow 创建滑动窗口限流器
func (f *LimiterFactory) CreateSlidingWindow(limit int, window time.Duration) Limiter {
	return NewSlidingWindow(limit, window)
}

// CreateFixedWindow 创建固定窗口限流器
func (f *LimiterFactory) CreateFixedWindow(limit int, window time.Duration) Limiter {
	return NewFixedWindow(limit, window)
}

// DefaultFactory 默认限流器工厂
var DefaultFactory = &LimiterFactory{}

// KeyGenerator 键生成器函数类型
type KeyGenerator func(ctx interface{}) string

// IPKeyGenerator IP地址键生成器
func IPKeyGenerator(ip string) string {
	return fmt.Sprintf("ip:%s", ip)
}

// UserKeyGenerator 用户键生成器
func UserKeyGenerator(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

// EndpointKeyGenerator 端点键生成器
func EndpointKeyGenerator(method, path string) string {
	return fmt.Sprintf("endpoint:%s:%s", method, path)
}

// CompositeKeyGenerator 组合键生成器
func CompositeKeyGenerator(parts ...string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ":"
		}
		result += part
	}
	return result
}
