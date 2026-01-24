package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username  string     `gorm:"size:100;not null;unique" json:"username"`
	Email     string     `gorm:"size:100;not null;unique" json:"email"`
	Password  string     `gorm:"size:100;not null" json:"-"`
	FirstName string     `gorm:"size:100" json:"first_name"`
	LastName  string     `gorm:"size:100" json:"last_name"`
	Avatar    string     `gorm:"size:255" json:"avatar"`
	LastLogin *time.Time `json:"last_login"`
	Status    string     `gorm:"size:20;default:'active'" json:"status"` // active, inactive, banned

	// 关联
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// UserProfile 用户资料（不包含敏感信息）
type UserProfile struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Avatar    string     `json:"avatar"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login"`
}

// ToProfile 将用户模型转换为用户资料
func (u *User) ToProfile() UserProfile {
	return UserProfile{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		LastLogin: u.LastLogin,
	}
}

// HasRole 检查用户是否有指定角色
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission 检查用户是否有指定权限
func (u *User) HasPermission(permissionName string) bool {
	for _, role := range u.Roles {
		if role.HasPermission(permissionName) {
			return true
		}
	}
	return false
}

// HasResourcePermission 检查用户是否有资源的指定操作权限
func (u *User) HasResourcePermission(resource, action string) bool {
	for _, role := range u.Roles {
		if role.HasResourcePermission(resource, action) {
			return true
		}
	}
	return false
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.HasRole("admin") || u.HasRole("super_admin")
}
