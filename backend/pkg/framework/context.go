package framework

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// RequestContext 请求上下文
type RequestContext struct {
	*app.RequestContext
}

// NewRequestContext 创建一个新的请求上下文
func NewRequestContext(ctx *app.RequestContext) *RequestContext {
	return &RequestContext{
		RequestContext: ctx,
	}
}

// JSON 返回JSON响应
func (c *RequestContext) JSON(code int, obj interface{}) {
	c.RequestContext.JSON(code, obj)
}

// String 返回字符串响应
func (c *RequestContext) String(code int, format string, values ...interface{}) {
	c.RequestContext.String(code, format, values...)
}

// HTML 返回HTML响应
func (c *RequestContext) HTML(code int, html string, data interface{}) {
	c.RequestContext.HTML(code, html, data)
}

// File 返回文件响应
func (c *RequestContext) File(filepath string) {
	c.RequestContext.File(filepath)
}

// Redirect 重定向
func (c *RequestContext) Redirect(code int, location string) {
	c.RequestContext.Redirect(code, []byte(location))
}

// GetParam 获取路由参数
func (c *RequestContext) GetParam(key string) string {
	return c.RequestContext.Param(key)
}

// GetQuery 获取查询参数
func (c *RequestContext) GetQuery(key string) string {
	return c.RequestContext.Query(key)
}

// GetForm 获取表单参数
func (c *RequestContext) GetForm(key string) string {
	return c.RequestContext.PostForm(key)
}

// GetHeader 获取请求头
func (c *RequestContext) GetHeader(key string) string {
	return string(c.RequestContext.GetHeader(key))
}

// SetHeader 设置响应头
func (c *RequestContext) SetHeader(key, value string) {
	c.RequestContext.Header(key, value)
}

// GetCookie 获取Cookie
func (c *RequestContext) GetCookie(key string) string {
	return string(c.RequestContext.Cookie(key))
}

// SetCookie 设置Cookie
func (c *RequestContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.RequestContext.SetCookie(name, value, maxAge, path, domain, 0, secure, httpOnly)
}

// BindJSON 绑定JSON请求体
func (c *RequestContext) BindJSON(obj interface{}) error {
	return c.RequestContext.BindJSON(obj)
}

// BindForm 绑定表单请求体
func (c *RequestContext) BindForm(obj interface{}) error {
	return c.RequestContext.BindForm(obj)
}

// ClientIP 获取客户端IP
func (c *RequestContext) ClientIP() string {
	return c.RequestContext.ClientIP()
}

// UserAgent 获取用户代理
func (c *RequestContext) UserAgent() string {
	return string(c.RequestContext.UserAgent())
}

// Method 获取请求方法
func (c *RequestContext) Method() string {
	return string(c.RequestContext.Method())
}

// Path 获取请求路径
func (c *RequestContext) Path() string {
	return string(c.RequestContext.Path())
}

// ContentType 获取内容类型
func (c *RequestContext) ContentType() string {
	return string(c.RequestContext.ContentType())
}

// Status 设置状态码
func (c *RequestContext) Status(code int) {
	c.RequestContext.Status(code)
}

// Abort 中止请求处理
func (c *RequestContext) Abort() {
	c.RequestContext.Abort()
}

// AbortWithStatus 中止请求处理并设置状态码
func (c *RequestContext) AbortWithStatus(code int) {
	c.RequestContext.AbortWithStatus(code)
}

// Next 继续处理请求
func (c *RequestContext) Next(ctx context.Context) {
	c.RequestContext.Next(ctx)
}

// Set 设置上下文值
func (c *RequestContext) Set(key string, value interface{}) {
	c.RequestContext.Set(key, value)
}

// Get 获取上下文值
func (c *RequestContext) Get(key string) (interface{}, bool) {
	return c.RequestContext.Get(key)
}
