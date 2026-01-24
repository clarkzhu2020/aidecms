package cache

import (
	"errors"
	"sync"
	"time"
)

// Driver 缓存驱动接口
type Driver interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) bool
	Clear() error
}

// Cache 缓存管理器
type Cache struct {
	driver Driver
	mu     sync.RWMutex
}

// NewCache 创建一个新的缓存管理器
func NewCache(driver Driver) *Cache {
	return &Cache{
		driver: driver,
	}
}

// Get 获取缓存
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.driver.Get(key)
}

// GetString 获取字符串缓存
func (c *Cache) GetString(key string) (string, error) {
	value, err := c.Get(key)
	if err != nil {
		return "", err
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	return "", errors.New("value is not a string")
}

// GetInt 获取整数缓存
func (c *Cache) GetInt(key string) (int, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	}

	return 0, errors.New("value is not an integer")
}

// GetBool 获取布尔缓存
func (c *Cache) GetBool(key string) (bool, error) {
	value, err := c.Get(key)
	if err != nil {
		return false, err
	}

	if b, ok := value.(bool); ok {
		return b, nil
	}

	return false, errors.New("value is not a boolean")
}

// Set 设置缓存
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.driver.Set(key, value, ttl)
}

// Delete 删除缓存
func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.driver.Delete(key)
}

// Exists 检查缓存是否存在
func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.driver.Exists(key)
}

// Clear 清空缓存
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.driver.Clear()
}
