package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator 验证器
type Validator struct {
	validate *validator.Validate
}

// NewValidator 创建验证器实例
func NewValidator() *Validator {
	v := validator.New()

	// 注册自定义验证规则
	v.RegisterValidation("slug", validateSlug)
	v.RegisterValidation("username", validateUsername)

	return &Validator{
		validate: v,
	}
}

// Validate 验证结构体
func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return NewValidationError(validationErrors)
		}
		return err
	}
	return nil
}

// ValidateVar 验证单个变量
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// ValidationError 验证错误
type ValidationError struct {
	Errors map[string][]string `json:"errors"`
}

// NewValidationError 创建验证错误
func NewValidationError(errs validator.ValidationErrors) *ValidationError {
	errors := make(map[string][]string)

	for _, err := range errs {
		field := strings.ToLower(err.Field())
		message := formatErrorMessage(err)

		if _, exists := errors[field]; !exists {
			errors[field] = []string{}
		}
		errors[field] = append(errors[field], message)
	}

	return &ValidationError{
		Errors: errors,
	}
}

// Error 实现error接口
func (e *ValidationError) Error() string {
	var messages []string
	for field, errs := range e.Errors {
		for _, err := range errs {
			messages = append(messages, fmt.Sprintf("%s: %s", field, err))
		}
	}
	return strings.Join(messages, "; ")
}

// formatErrorMessage 格式化错误消息
func formatErrorMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "slug":
		return fmt.Sprintf("%s must be a valid slug (lowercase letters, numbers, hyphens only)", field)
	case "username":
		return fmt.Sprintf("%s must be a valid username (letters, numbers, underscore only)", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "numeric":
		return fmt.Sprintf("%s must contain only numbers", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// 自定义验证规则

// validateSlug 验证slug格式
func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	// slug只能包含小写字母、数字和连字符
	for _, char := range slug {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return false
		}
	}
	return len(slug) > 0
}

// validateUsername 验证用户名格式
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// 用户名只能包含字母、数字和下划线
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return len(username) >= 3 && len(username) <= 20
}

// 全局验证器实例
var defaultValidator = NewValidator()

// Validate 使用全局验证器验证
func Validate(data interface{}) error {
	return defaultValidator.Validate(data)
}

// ValidateVar 使用全局验证器验证变量
func ValidateVar(field interface{}, tag string) error {
	return defaultValidator.ValidateVar(field, tag)
}
