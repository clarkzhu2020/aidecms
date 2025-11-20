package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config Redis配置
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Client Redis客户端
type Client struct {
	client *redis.Client
	Config *Config
}

// NewClient 创建新的Redis客户端
func NewClient(config *Config) *Client {
	return &Client{
		Config: config,
	}
}

// Connect 连接Redis
func (c *Client) Connect() error {
	c.client = redis.NewClient(&redis.Options{
		Addr:     c.Config.Host + ":" + c.Config.Port,
		Password: c.Config.Password,
		DB:       c.Config.DB,
	})

	_, err := c.client.Ping(context.Background()).Result()
	return err
}

// Close 关闭Redis连接
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Ping 检查连接
func (c *Client) Ping(ctx context.Context) (string, error) {
	return c.client.Ping(ctx).Result()
}

// Get 获取键值
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 设置键值
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Del 删除键
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire 设置键过期时间
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// HGet 获取哈希字段值
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HSet 设置哈希字段值
func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.client.HSet(ctx, key, values...).Err()
}

// HGetAll 获取所有哈希字段
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// LPush 列表左推入
func (c *Client) LPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.LPush(ctx, key, values...).Err()
}

// RPop 列表右弹出
func (c *Client) RPop(ctx context.Context, key string) (string, error) {
	return c.client.RPop(ctx, key).Result()
}

// Publish 发布消息
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

// Subscribe 订阅频道
func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}
