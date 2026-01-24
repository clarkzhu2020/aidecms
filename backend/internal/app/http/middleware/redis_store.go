package middleware

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStore Redis存储实现
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore 创建Redis存储实例
func NewRedisStore() *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // 从配置读取
			Password: "",               // 无密码
			DB:       0,                // 默认DB
		}),
	}
}

func (s *RedisStore) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *RedisStore) Set(key, value string, expire time.Duration) error {
	ctx := context.Background()
	return s.client.Set(ctx, key, value, expire).Err()
}

func (s *RedisStore) Delete(key string) error {
	ctx := context.Background()
	return s.client.Del(ctx, key).Err()
}

func (s *RedisStore) Exists(key string) bool {
	ctx := context.Background()
	exists, _ := s.client.Exists(ctx, key).Result()
	return exists > 0
}
