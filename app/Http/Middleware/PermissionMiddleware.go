package middleware

import (
	"context"
	"fmt"
	"strconv"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/internal/app/services"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/clarkzhu2020/aidecms/pkg/framework"
	"github.com/cloudwego/hertz/pkg/app"
)

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(permissionName string) framework.HandlerFunc {
	return func(ctx context.Context, reqCtx *framework.RequestContext) {
		// 从上下文获取用户ID（假设JWT中间件已经设置）
		userIDStr, exists := reqCtx.Get("user_id")
		if !exists {
			reqCtx.JSON(401, map[string]interface{}{
				"success": false,
				"error":   "Unauthorized",
				"message": "Authentication required",
			})
			reqCtx.Abort()
			return
		}

		userID, err := strconv.ParseUint(fmt.Sprint(userIDStr), 10, 32)
		if err != nil {
			reqCtx.JSON(401, map[string]interface{}{
				"success": false,
				"error":   "Invalid user ID",
			})
			reqCtx.Abort()
			return
		}

		// 检查权限
		permService := services.NewPermissionService()
		hasPermission, err := permService.CheckUserPermission(uint(userID), permissionName)
		if err != nil {
			reqCtx.JSON(500, map[string]interface{}{
				"success": false,
				"error":   "Failed to check permission",
				"message": err.Error(),
			})
			reqCtx.Abort()
			return
		}

		if !hasPermission {
			reqCtx.JSON(403, map[string]interface{}{
				"success": false,
				"error":   "Forbidden",
				"message": "You don't have permission to access this resource",
			})
			reqCtx.Abort()
			return
		}

		reqCtx.Next(ctx)
	}
}

// ResourcePermissionMiddleware 资源权限检查中间件
func ResourcePermissionMiddleware(resource, action string) framework.HandlerFunc {
	return func(ctx context.Context, reqCtx *framework.RequestContext) {
		userIDStr, exists := reqCtx.Get("user_id")
		if !exists {
			reqCtx.JSON(401, map[string]interface{}{
				"success": false,
				"error":   "Unauthorized",
				"message": "Authentication required",
			})
			reqCtx.Abort()
			return
		}

		userID, err := strconv.ParseUint(fmt.Sprint(userIDStr), 10, 32)
		if err != nil {
			reqCtx.JSON(401, map[string]interface{}{
				"success": false,
				"error":   "Invalid user ID",
			})
			reqCtx.Abort()
			return
		}

		permService := services.NewPermissionService()
		hasPermission, err := permService.CheckUserResourcePermission(uint(userID), resource, action)
		if err != nil {
			reqCtx.JSON(500, map[string]interface{}{
				"success": false,
				"error":   "Failed to check permission",
				"message": err.Error(),
			})
			reqCtx.Abort()
			return
		}

		if !hasPermission {
			reqCtx.JSON(403, map[string]interface{}{
				"success": false,
				"error":   "Forbidden",
				"message": fmt.Sprintf("You don't have permission to %s %s", action, resource),
			})
			reqCtx.Abort()
			return
		}

		reqCtx.Next(ctx)
	}
}

// RoleMiddleware 角色检查中间件
func RoleMiddleware(roleName string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, map[string]interface{}{
				"success": false,
				"error":   "Unauthorized",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(fmt.Sprint(userIDStr), 10, 32)
		if err != nil {
			c.JSON(401, map[string]interface{}{
				"success": false,
				"error":   "Invalid user ID",
			})
			c.Abort()
			return
		}

		// 检查用户是否有此角色
		var user models.User
		db := database.GetDB()
		if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
			c.JSON(401, map[string]interface{}{
				"success": false,
				"error":   "User not found",
			})
			c.Abort()
			return
		}

		hasRole := user.HasRole(roleName)

		if !hasRole {
			c.JSON(403, map[string]interface{}{
				"success": false,
				"error":   "Forbidden",
				"message": fmt.Sprintf("You must have %s role to access this resource", roleName),
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
