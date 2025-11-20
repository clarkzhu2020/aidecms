package adapters

import (
	"context"

	"github.com/chenyusolar/aidecms/pkg/framework"
	"github.com/cloudwego/hertz/pkg/app"
)

// HertzToFramework 将Hertz的处理函数转换为framework的处理函数
func HertzToFramework(handler app.HandlerFunc) framework.HandlerFunc {
	return func(ctx context.Context, c *framework.RequestContext) {
		// framework.RequestContext 包含了 *app.RequestContext
		handler(ctx, c.RequestContext)
	}
}

// ControllerToFramework 将控制器方法转换为framework的处理函数
func ControllerToFramework(handler func(ctx context.Context, c *framework.RequestContext)) framework.HandlerFunc {
	return handler
}
