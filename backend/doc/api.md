# API 开发

AideCMS 基于 Hertz 框架提供了强大的 API 开发支持。

## 路由定义

在 `internal/app/http/routes.go` 中定义 API 路由:

```go
package http

import "github.com/cloudwego/hertz/pkg/app/server"

func RegisterAPIRoutes(h *server.Hertz) {
    api := h.Group("/api")
    {
        api.GET("/users", userController.Index)
        api.POST("/users", userController.Store)
        api.GET("/users/:id", userController.Show)
    }
}
```

## 控制器示例

`internal/app/http/controllers/user_controller.go`:

```go
package controllers

import (
    "github.com/cloudwego/hertz/pkg/app"
    "internal/app/models"
)

type UserController struct{}

func (c *UserController) Index(ctx context.Context, hCtx *app.RequestContext) {
    var users []models.User
    models.DB.Find(&users)
    
    hCtx.JSON(200, users)
}

func (c *UserController) Show(ctx context.Context, hCtx *app.RequestContext) {
    id := hCtx.Param("id")
    var user models.User
    
    if err := models.DB.First(&user, id).Error; err != nil {
        hCtx.JSON(404, map[string]string{"error": "User not found"})
        return
    }
    
    hCtx.JSON(200, user)
}
```

## 请求验证

使用内置验证器:

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func (c *UserController) Store(ctx context.Context, hCtx *app.RequestContext) {
    var req CreateUserRequest
    if err := hCtx.BindAndValidate(&req); err != nil {
        hCtx.JSON(400, map[string]string{"error": err.Error()})
        return
    }
    
    // 处理请求...
}
```

## 响应格式

统一响应格式:

```go
func successResponse(hCtx *app.RequestContext, data interface{}) {
    hCtx.JSON(200, map[string]interface{}{
        "success": true,
        "data":    data,
    })
}

func errorResponse(hCtx *app.RequestContext, code int, message string) {
    hCtx.JSON(code, map[string]interface{}{
        "success": false,
        "error":   message,
    })
}
```

## 中间件

创建中间件:

```go
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, hCtx *app.RequestContext) {
        token := hCtx.GetHeader("Authorization")
        if token == "" {
            hCtx.AbortWithStatusJSON(401, map[string]string{"error": "Unauthorized"})
            return
        }
        
        // 验证token...
        hCtx.Next(ctx)
    }
}
```

注册中间件:

```go
api.Use(middleware.AuthMiddleware())
```

## 文档生成

使用 Swagger:

1. 安装 swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. 添加注释到控制器:
```go
// @Summary 获取用户列表
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /api/users [get]
func (c *UserController) Index(ctx context.Context, hCtx *app.RequestContext) {
    // ...
}
```

3. 生成文档:
```bash
swag init -g internal/app/http/routes.go
```

4. 访问文档:
```
http://localhost:8888/swagger/index.html