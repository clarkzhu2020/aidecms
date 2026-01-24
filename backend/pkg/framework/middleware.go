package framework

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Middleware 中间件管理器
type Middleware struct {
	handlers []app.HandlerFunc
}

// NewMiddleware 创建一个新的中间件管理器
func NewMiddleware() *Middleware {
	return &Middleware{
		handlers: make([]app.HandlerFunc, 0),
	}
}

// Add 添加中间件
func (m *Middleware) Add(handler app.HandlerFunc) *Middleware {
	m.handlers = append(m.handlers, handler)
	return m
}

// Remove 移除中间件
func (m *Middleware) Remove(handler app.HandlerFunc) *Middleware {
	for i, h := range m.handlers {
		if &h == &handler {
			m.handlers = append(m.handlers[:i], m.handlers[i+1:]...)
			break
		}
	}
	return m
}

// GetHandlers 获取所有中间件
func (m *Middleware) GetHandlers() []app.HandlerFunc {
	return m.handlers
}

// Clear 清空中间件
func (m *Middleware) Clear() *Middleware {
	m.handlers = make([]app.HandlerFunc, 0)
	return m
}

// Cors 跨域中间件
func Cors() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

		if string(ctx.Request.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next(c)
	}
}

// Recovery 恢复中间件
func Recovery() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				ctx.JSON(500, map[string]interface{}{
					"code":    500,
					"message": "Internal Server Error",
				})
				ctx.Abort()
			}
		}()
		ctx.Next(c)
	}
}

// Logger 日志中间件
func Logger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		path := string(ctx.Request.URI().Path())
		method := string(ctx.Request.Method())

		ctx.Next(c)

		latency := time.Since(start)
		statusCode := ctx.Response.StatusCode()

		hlog.Infof("[%s] %s %d %s", method, path, statusCode, latency)
	}
}
