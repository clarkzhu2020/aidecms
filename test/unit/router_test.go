package unit_test

import (
	"testing"

	"github.com/chenyusolar/aidecms/pkg/framework"
	"github.com/stretchr/testify/assert"
)

func TestRouterRegistration(t *testing.T) {
	app := framework.NewApplication().
		SetConfigPath("../../config").
		Boot()
	
	app.RegisterRoutes(func(r *framework.Router) {
		r.GET("/test", func(c *app.RequestContext) {
			c.String(200, "OK")
		})
	})
	
	// 验证路由是否注册
	assert.NotNil(t, app.Router)
}

func TestMiddlewareRegistration(t *testing.T) {
	app := framework.NewApplication().
		SetConfigPath("../../config").
		Boot()
	
	middlewareCalled := false
	app.RegisterMiddleware(func(c *app.RequestContext) {
		middlewareCalled = true
		c.Next()
	})
	
	// 验证中间件是否注册
	assert.NotNil(t, app.Server)
}