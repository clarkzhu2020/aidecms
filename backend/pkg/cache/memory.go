package cache

import (
	"errors"
	"sync"
	"time"
)

// MemoryItem 内存缓存项
type MemoryItem struct {
	Value      interface{}
	Expiration int64
}

// MemoryDriver 内存缓存驱动
type MemoryDriver struct {
	items map[string]MemoryItem
	mu    sync.RWMutex
}

// NewMemoryDriver 创建一个新的内存缓存驱动
func NewMemoryDriver() *MemoryDriver {
	driver := &MemoryDriver{
		items: make(map[string]MemoryItem),
	}

	// 启动过期清理
	go driver.startGC()

	return driver
}

// Get 获取缓存
func (d *MemoryDriver) Get(key string) (interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	item, found := d.items[key]
	if !found {
		return nil, errors.New("key not found")
	}

	// 检查是否过期
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		return nil, errors.New("key expired")
	}

	return item.Value, nil
}

// Set 设置缓存
func (d *MemoryDriver) Set(key string, value interface{}, ttl time.Duration) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	d.items[key] = MemoryItem{
		Value:      value,
		Expiration: expiration,
	}

	return nil
}

// Delete 删除缓存
func (d *MemoryDriver) Delete(key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, found := d.items[key]; !found {
		return errors.New("key not found")
	}

	delete(d.items, key)
	return nil
}

// Exists 检查缓存是否存在
func (d *MemoryDriver) Exists(key string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	item, found := d.items[key]
	if !found {
		return false
	}

	// 检查是否过期
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		return false
	}

	return true
}

// Clear 清空缓存
func (d *MemoryDriver) Clear() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.items = make(map[string]MemoryItem)
	return nil
}

// startGC 启动垃圾回收
func (d *MemoryDriver) startGC() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		d.deleteExpired()
	}
}

// deleteExpired 删除过期缓存
func (d *MemoryDriver) deleteExpired() {
	now := time.Now().UnixNano()

	d.mu.Lock()
	defer d.mu.Unlock()

	for key, item := range d.items {
		if item.Expiration > 0 && item.Expiration < now {
			delete(d.items, key)
		}
	}
}
