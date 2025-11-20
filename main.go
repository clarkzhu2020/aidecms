package main

import (
	"context"

	_ "github.com/clarkzhu2020/aidecms/docs" // Swagger docs
	"github.com/clarkzhu2020/aidecms/pkg/framework"
	"github.com/clarkzhu2020/aidecms/pkg/swagger"
	"github.com/clarkzhu2020/aidecms/routes"
)

// @title           AideCMS API
// @version         1.0
// @description     AideCMS API文档 - 基于Hertz框架的高性能AI辅助内容管理系统
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@aidecms.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8888
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
		// Swagger文档路由 - 访问 http://localhost:8888/swagger/index.html
		router.GET("/swagger/*any", swagger.SwaggerHandler())

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
			c.String(200, "Welcome to AideCMS with AI capabilities! See /doc/ai.md for AI integration guide.")
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
