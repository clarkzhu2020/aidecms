package controllers

import (
	"context"

	"github.com/clarkzhu2020/aidecms/internal/app/services"
	"github.com/clarkzhu2020/aidecms/pkg/framework"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	App         *framework.Application
	UserService *services.UserService
}

func NewUserController(app *framework.Application) *UserController {
	return &UserController{
		App:         app,
		UserService: services.NewUserService(),
	}
}

// 用户注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// 用户登录请求
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// 更新用户资料请求
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

// Register 用户注册
func (c *UserController) Register(ctx context.Context, reqCtx *framework.RequestContext) {
	var req RegisterRequest
	if err := reqCtx.BindJSON(&req); err != nil {
		reqCtx.JSON(400, map[string]interface{}{
			"error": "无效的请求数据",
		})
		return
	}

	// 调用用户服务注册用户
	user, err := c.UserService.Register(req.Username, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		reqCtx.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// 生成JWT令牌
	token, err := c.UserService.GenerateJWT(user.ID)
	if err != nil {
		reqCtx.JSON(500, map[string]interface{}{
			"error": "生成令牌失败",
		})
		return
	}

	// 返回用户信息和令牌
	reqCtx.JSON(201, map[string]interface{}{
		"user":  user.ToProfile(),
		"token": token,
	})
}

// Login 用户登录
func (c *UserController) Login(ctx context.Context, reqCtx *framework.RequestContext) {
	var req LoginRequest
	if err := reqCtx.BindJSON(&req); err != nil {
		reqCtx.JSON(400, map[string]interface{}{
			"error": "无效的请求数据",
		})
		return
	}

	// 调用用户服务登录
	user, token, err := c.UserService.Login(req.UsernameOrEmail, req.Password)
	if err != nil {
		reqCtx.JSON(401, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// 返回用户信息和令牌
	reqCtx.JSON(200, map[string]interface{}{
		"user":  user.ToProfile(),
		"token": token,
	})
}

// Profile 获取用户资料
func (c *UserController) Profile(ctx context.Context, reqCtx *framework.RequestContext) {
	// 从上下文中获取用户ID
	userID, exists := reqCtx.Get("user_id")
	if !exists {
		reqCtx.JSON(401, map[string]interface{}{
			"error": "未授权",
		})
		return
	}

	// 获取用户信息
	user, err := c.UserService.GetUserByID(userID.(uint))
	if err != nil {
		reqCtx.JSON(404, map[string]interface{}{
			"error": "用户不存在",
		})
		return
	}

	// 返回用户资料
	reqCtx.JSON(200, map[string]interface{}{
		"user": user.ToProfile(),
	})
}

// UpdateProfile 更新用户资料
func (c *UserController) UpdateProfile(ctx context.Context, reqCtx *framework.RequestContext) {
	// 从上下文中获取用户ID
	userID, exists := reqCtx.Get("user_id")
	if !exists {
		reqCtx.JSON(401, map[string]interface{}{
			"error": "未授权",
		})
		return
	}

	var req UpdateProfileRequest
	if err := reqCtx.BindJSON(&req); err != nil {
		reqCtx.JSON(400, map[string]interface{}{
			"error": "无效的请求数据",
		})
		return
	}

	// 更新用户资料
	user, err := c.UserService.UpdateProfile(userID.(uint), req.FirstName, req.LastName, req.Avatar)
	if err != nil {
		reqCtx.JSON(500, map[string]interface{}{
			"error": "更新资料失败",
		})
		return
	}

	// 返回更新后的用户资料
	reqCtx.JSON(200, map[string]interface{}{
		"user": user.ToProfile(),
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
