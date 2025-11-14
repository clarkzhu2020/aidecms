package routes

import (
	"context"

	"github.com/clarkgo/clarkgo/pkg/framework"
)

// TestRoutes 测试路由
func TestRoutes(app *framework.Application) {
	app.RegisterRoutes(func(r *framework.Router) {
		r.GET("/test", func(ctx context.Context, c *framework.RequestContext) {
			c.JSON(200, map[string]interface{}{
				"message": "测试成功",
			})
		})
	})
}
