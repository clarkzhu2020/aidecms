package middleware

import (
	"context"
	"strings"

	"github.com/chenyusolar/aidecms/internal/app/services"
	"github.com/chenyusolar/aidecms/pkg/framework"
)

// JWTMiddleware JWT认证中间件
func JWTMiddleware() framework.HandlerFunc {
	userService := services.NewUserService()

	return func(ctx context.Context, reqCtx *framework.RequestContext) {
		// 获取Authorization头
		authHeader := reqCtx.GetHeader("Authorization")
		if authHeader == "" {
			reqCtx.JSON(401, map[string]interface{}{
				"error": "未提供授权令牌",
			})
			reqCtx.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			reqCtx.JSON(401, map[string]interface{}{
				"error": "无效的授权格式，应为'Bearer {token}'",
			})
			reqCtx.Abort()
			return
		}

		// 提取令牌
		tokenString := tokenParts[1]

		// 解析令牌
		userID, err := userService.ParseJWT(tokenString)
		if err != nil {
			reqCtx.JSON(401, map[string]interface{}{
				"error": "无效的令牌: " + err.Error(),
			})
			reqCtx.Abort()
			return
		}

		// 将用户ID添加到上下文
		reqCtx.Set("user_id", userID)

		// 继续处理请求
		reqCtx.Next(ctx)
	}
}
