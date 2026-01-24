package middleware

import (
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB // 假设全局db实例已定义

// DatabaseStore 数据库存储实现
type DatabaseStore struct {
	db *gorm.DB
}

// NewDatabaseStore 创建数据库存储实例
func NewDatabaseStore() SessionStore {
	return &DatabaseStore{
		db: db, // 使用全局数据库连接
	}
}

// Get 获取会话数据
func (s *DatabaseStore) Get(key string) (string, error) {
	// 简化实现，实际应该从数据库读取
	return "", nil
}

// Set 设置会话数据
func (s *DatabaseStore) Set(key string, value string, expire time.Duration) error {
	// 简化实现，实际应该写入数据库
	return nil
}

// Delete 删除会话数据
func (s *DatabaseStore) Delete(key string) error {
	// 简化实现，实际应该从数据库删除
	return nil
}

// Exists 检查会话是否存在
func (s *DatabaseStore) Exists(key string) bool {
	// 简化实现，实际应该查询数据库
	return false
}
