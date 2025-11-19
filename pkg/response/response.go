package response

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta 分页元数据
type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int64 `json:"total_pages"`
}

// Success 成功响应
func Success(ctx *app.RequestContext, data interface{}, message string) {
	ctx.JSON(200, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SuccessWithMeta 带分页信息的成功响应
func SuccessWithMeta(ctx *app.RequestContext, data interface{}, meta *Meta, message string) {
	ctx.JSON(200, Response{
		Success: true,
		Data:    data,
		Message: message,
		Meta:    meta,
	})
}

// Created 创建成功响应
func Created(ctx *app.RequestContext, data interface{}, message string) {
	if message == "" {
		message = "Resource created successfully"
	}
	ctx.JSON(201, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// NoContent 无内容响应
func NoContent(ctx *app.RequestContext) {
	ctx.Status(204)
}

// Error 错误响应
func Error(ctx *app.RequestContext, statusCode int, error string, message string) {
	ctx.JSON(statusCode, Response{
		Success: false,
		Error:   error,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(ctx *app.RequestContext, message string) {
	Error(ctx, 400, "Bad Request", message)
}

// Unauthorized 401错误
func Unauthorized(ctx *app.RequestContext, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	Error(ctx, 401, "Unauthorized", message)
}

// Forbidden 403错误
func Forbidden(ctx *app.RequestContext, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	Error(ctx, 403, "Forbidden", message)
}

// NotFound 404错误
func NotFound(ctx *app.RequestContext, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(ctx, 404, "Not Found", message)
}

// ValidationError 验证错误响应
func ValidationError(ctx *app.RequestContext, errors interface{}) {
	ctx.JSON(422, map[string]interface{}{
		"success": false,
		"error":   "Validation Failed",
		"message": "The given data was invalid",
		"errors":  errors,
	})
}

// ServerError 500错误
func ServerError(ctx *app.RequestContext, message string) {
	if message == "" {
		message = "Internal server error"
	}
	Error(ctx, 500, "Internal Server Error", message)
}

// NewMeta 创建分页元数据
func NewMeta(currentPage, perPage int, total int64) *Meta {
	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	if totalPages < 0 {
		totalPages = 0
	}

	return &Meta{
		CurrentPage: currentPage,
		PerPage:     perPage,
		Total:       total,
		TotalPages:  totalPages,
	}
}
