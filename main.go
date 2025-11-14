package main

import (
	"context"

	"github.com/clarkgo/clarkgo/pkg/framework"
	"github.com/clarkgo/clarkgo/routes"
)

func main() {
	// 创建并启动应用实例
	app := framework.NewApplication().
		SetConfigPath("config").
		SetEnv("development").
		SetDebug(true).
		Boot()

	// 注册API路由
	routes.APIRoutes(app)

	// 注册其他路由
	app.RegisterRoutes(func(router *framework.Router) {
		// 基本API路由
		api := router.Group("/api")
		{
			api.GET("/ping", func(ctx context.Context, c *framework.RequestContext) {
				c.JSON(200, map[string]interface{}{
					"message": "pong",
				})
			})
		}

		// Web路由
		router.GET("/", func(ctx context.Context, c *framework.RequestContext) {
			c.String(200, "Welcome to ClarkGo with AI capabilities! See /doc/ai.md for AI integration guide.")
		})
	})

	// 注册中间件
	app.RegisterMiddleware(
		framework.Cors(),
		framework.Recovery(),
		framework.Logger(),
	)

	// 注册静态文件目录
	app.Static("/public", app.GetPublicPath())

	// 运行应用
	app.Run()
}
