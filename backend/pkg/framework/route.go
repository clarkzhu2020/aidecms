package framework

import (
	"context"
	"fmt"
	"sort"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RouteInfo 存储路由信息
type RouteInfo struct {
	Method  string
	Path    string
	Handler string
}

// Router 路由管理器
type Router struct {
	server *server.Hertz
	prefix string
	routes []RouteInfo // 存储所有注册的路由
}

// HandlerFunc 路由处理函数类型
type HandlerFunc func(context.Context, *RequestContext)

// NewRouter 创建一个新的路由管理器
func NewRouter(server *server.Hertz) *Router {
	return &Router{
		server: server,
		prefix: "",
		routes: []RouteInfo{},
	}
}

// PrintRoutes 打印所有已注册的路由
func (r *Router) PrintRoutes() {
	if len(r.routes) == 0 {
		fmt.Println("No routes registered.")
		return
	}

	// 按路径排序
	sort.Slice(r.routes, func(i, j int) bool {
		return r.routes[i].Path < r.routes[j].Path
	})

	// 打印表头
	fmt.Println("\n+--------+------------------------------------+------------------------------------+")
	fmt.Println("| METHOD | PATH                               | HANDLER                            |")
	fmt.Println("+--------+------------------------------------+------------------------------------+")

	// 打印路由
	for _, route := range r.routes {
		method := fmt.Sprintf("%-6s", route.Method)
		path := route.Path
		if len(path) > 34 {
			path = path[:31] + "..."
		} else {
			path = fmt.Sprintf("%-34s", path)
		}

		handler := route.Handler
		if len(handler) > 34 {
			handler = handler[:31] + "..."
		} else {
			handler = fmt.Sprintf("%-34s", handler)
		}

		fmt.Printf("| %s | %s | %s |\n", method, path, handler)
	}

	fmt.Println("+--------+------------------------------------+------------------------------------+")
	fmt.Printf("\nTotal routes: %d\n\n", len(r.routes))
}

// GetRoutes 获取所有已注册的路由
func (r *Router) GetRoutes() []RouteInfo {
	return r.routes
}

// Group 创建一个路由组
func (r *Router) Group(prefix string, handlers ...HandlerFunc) *Router {
	// 将我们的 HandlerFunc 转换为 Hertz 的 HandlerFunc
	h := make([]app.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		h[i] = func(ctx context.Context, c *app.RequestContext) {
			handler(ctx, NewRequestContext(c))
		}
	}

	// 创建路由组
	r.server.Group(r.prefix+prefix, h...)
	return &Router{
		server: r.server,
		prefix: r.prefix + prefix,
	}
}

// GET 注册GET路由
func (r *Router) GET(path string, handler HandlerFunc) {
	r.server.GET(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	r.routes = append(r.routes, RouteInfo{
		Method:  "GET",
		Path:    r.prefix + path,
		Handler: handlerName,
	})
}

// POST 注册POST路由
func (r *Router) POST(path string, handler HandlerFunc) {
	r.server.POST(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	r.routes = append(r.routes, RouteInfo{
		Method:  "POST",
		Path:    r.prefix + path,
		Handler: handlerName,
	})
}

// PUT 注册PUT路由
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.server.PUT(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	r.routes = append(r.routes, RouteInfo{
		Method:  "PUT",
		Path:    r.prefix + path,
		Handler: handlerName,
	})
}

// DELETE 注册DELETE路由
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.server.DELETE(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	r.routes = append(r.routes, RouteInfo{
		Method:  "DELETE",
		Path:    r.prefix + path,
		Handler: handlerName,
	})
}

// PATCH 注册PATCH路由
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.server.PATCH(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	r.routes = append(r.routes, RouteInfo{
		Method:  "PATCH",
		Path:    r.prefix + path,
		Handler: handlerName,
	})
}

// OPTIONS 注册OPTIONS路由
func (r *Router) OPTIONS(path string, handler HandlerFunc) {
	r.server.OPTIONS(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	r.routes = append(r.routes, RouteInfo{
		Method:  "OPTIONS",
		Path:    r.prefix + path,
		Handler: handlerName,
	})
}

// HEAD 注册HEAD路由
func (r *Router) HEAD(path string, handler HandlerFunc) {
	r.server.HEAD(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	r.routes = append(r.routes, RouteInfo{
		Method:  "HEAD",
		Path:    r.prefix + path,
		Handler: handlerName,
	})
}

// Any 注册所有HTTP方法的路由
func (r *Router) Any(path string, handler HandlerFunc) {
	r.server.Any(r.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		handler(ctx, NewRequestContext(c))
	})

	// 收集路由信息
	handlerName := fmt.Sprintf("%T", handler)
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, method := range methods {
		r.routes = append(r.routes, RouteInfo{
			Method:  method,
			Path:    r.prefix + path,
			Handler: handlerName,
		})
	}
}

// Static 注册静态文件路由
func (r *Router) Static(path, root string) {
	r.server.Static(r.prefix+path, root)

	// 收集路由信息
	r.routes = append(r.routes, RouteInfo{
		Method:  "GET",
		Path:    r.prefix + path + "/*filepath",
		Handler: "Static(" + root + ")",
	})
}

// StaticFile 注册静态文件路由
func (r *Router) StaticFile(path, filepath string) {
	r.server.StaticFile(r.prefix+path, filepath)

	// 收集路由信息
	r.routes = append(r.routes, RouteInfo{
		Method:  "GET",
		Path:    r.prefix + path,
		Handler: "StaticFile(" + filepath + ")",
	})
}

// StaticFS 注册静态文件系统路由
func (r *Router) StaticFS(path string, fs *app.FS) {
	r.server.StaticFS(r.prefix+path, fs)
}

// Use 使用中间件
func (r *Router) Use(handlers ...HandlerFunc) {
	h := make([]app.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		h[i] = func(ctx context.Context, c *app.RequestContext) {
			handler(ctx, NewRequestContext(c))
		}
	}
	r.server.Use(h...)
}
