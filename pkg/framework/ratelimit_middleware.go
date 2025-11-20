package framework

import (
	"context"
	"fmt"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/ratelimit"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// Limiter 限流器实例
	Limiter ratelimit.Limiter

	// KeyFunc 键生成函数
	KeyFunc func(ctx context.Context, c *app.RequestContext) string

	// ErrorHandler 错误处理函数
	ErrorHandler func(ctx context.Context, c *app.RequestContext)

	// SkipFunc 跳过函数
	SkipFunc func(ctx context.Context, c *app.RequestContext) bool
}

// RateLimit 限流中间件
func RateLimit(config RateLimitConfig) app.HandlerFunc {
	// 设置默认值
	if config.KeyFunc == nil {
		config.KeyFunc = defaultKeyFunc
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultErrorHandler
	}

	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否跳过
		if config.SkipFunc != nil && config.SkipFunc(ctx, c) {
			c.Next(ctx)
			return
		}

		// 生成键
		key := config.KeyFunc(ctx, c)

		// 检查是否允许
		if !config.Limiter.Allow(key) {
			config.ErrorHandler(ctx, c)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// defaultKeyFunc 默认键生成函数（基于 IP）
func defaultKeyFunc(ctx context.Context, c *app.RequestContext) string {
	return ratelimit.IPKeyGenerator(c.ClientIP())
}

// defaultErrorHandler 默认错误处理函数
func defaultErrorHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusTooManyRequests, map[string]interface{}{
		"success": false,
		"message": "Too many requests, please try again later",
		"error":   "rate_limit_exceeded",
	})
}

// RateLimitByIP IP限流中间件
func RateLimitByIP(rate, capacity int) app.HandlerFunc {
	limiter := ratelimit.NewTokenBucket(rate, capacity)
	return RateLimit(RateLimitConfig{
		Limiter: limiter,
		KeyFunc: func(ctx context.Context, c *app.RequestContext) string {
			return ratelimit.IPKeyGenerator(c.ClientIP())
		},
	})
}

// RateLimitByUser 用户限流中间件
func RateLimitByUser(rate, capacity int, getUserID func(context.Context, *app.RequestContext) string) app.HandlerFunc {
	limiter := ratelimit.NewTokenBucket(rate, capacity)
	return RateLimit(RateLimitConfig{
		Limiter: limiter,
		KeyFunc: func(ctx context.Context, c *app.RequestContext) string {
			userID := getUserID(ctx, c)
			if userID == "" {
				// 如果没有用户ID，使用IP
				return ratelimit.IPKeyGenerator(c.ClientIP())
			}
			return ratelimit.UserKeyGenerator(userID)
		},
	})
}

// RateLimitByEndpoint 端点限流中间件
func RateLimitByEndpoint(limit int, window time.Duration) app.HandlerFunc {
	limiter := ratelimit.NewSlidingWindow(limit, window)
	return RateLimit(RateLimitConfig{
		Limiter: limiter,
		KeyFunc: func(ctx context.Context, c *app.RequestContext) string {
			method := string(c.Method())
			path := string(c.Path())
			return ratelimit.EndpointKeyGenerator(method, path)
		},
	})
}

// RateLimitCombined 组合限流中间件（IP + 端点）
func RateLimitCombined(ipRate, ipCapacity, endpointLimit int, endpointWindow time.Duration) app.HandlerFunc {
	ipLimiter := ratelimit.NewTokenBucket(ipRate, ipCapacity)
	endpointLimiter := ratelimit.NewSlidingWindow(endpointLimit, endpointWindow)

	return func(ctx context.Context, c *app.RequestContext) {
		// 先检查IP限流
		ipKey := ratelimit.IPKeyGenerator(c.ClientIP())
		if !ipLimiter.Allow(ipKey) {
			c.JSON(consts.StatusTooManyRequests, map[string]interface{}{
				"success": false,
				"message": "Too many requests from your IP",
				"error":   "ip_rate_limit_exceeded",
			})
			c.Abort()
			return
		}

		// 再检查端点限流
		method := string(c.Method())
		path := string(c.Path())
		endpointKey := ratelimit.EndpointKeyGenerator(method, path)
		if !endpointLimiter.Allow(endpointKey) {
			c.JSON(consts.StatusTooManyRequests, map[string]interface{}{
				"success": false,
				"message": "Too many requests to this endpoint",
				"error":   "endpoint_rate_limit_exceeded",
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// RateLimitWithStats 带统计信息的限流中间件
func RateLimitWithStats(limit int, window time.Duration) app.HandlerFunc {
	limiter := ratelimit.NewSlidingWindow(limit, window)

	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		key := ratelimit.IPKeyGenerator(ip)

		// 检查是否允许
		if !limiter.Allow(key) {
			// 获取统计信息
			stats := limiter.GetStats(key)

			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

			c.JSON(consts.StatusTooManyRequests, map[string]interface{}{
				"success":   false,
				"message":   "Too many requests",
				"error":     "rate_limit_exceeded",
				"limit":     stats["limit"],
				"window":    stats["window"],
				"remaining": 0,
			})
			c.Abort()
			return
		}

		// 获取当前统计信息
		stats := limiter.GetStats(key)

		// 设置响应头
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", stats["remaining"]))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		c.Next(ctx)
	}
}

// RateLimitPerSecond 每秒限流（简化版）
func RateLimitPerSecond(limit int) app.HandlerFunc {
	return RateLimitByIP(limit, limit*2)
}

// RateLimitPerMinute 每分钟限流（简化版）
func RateLimitPerMinute(limit int) app.HandlerFunc {
	limiter := ratelimit.NewSlidingWindow(limit, time.Minute)
	return RateLimit(RateLimitConfig{
		Limiter: limiter,
	})
}

// RateLimitPerHour 每小时限流（简化版）
func RateLimitPerHour(limit int) app.HandlerFunc {
	limiter := ratelimit.NewSlidingWindow(limit, time.Hour)
	return RateLimit(RateLimitConfig{
		Limiter: limiter,
	})
}
