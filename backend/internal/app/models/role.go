package models

import (
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	gorm.Model
	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	DisplayName string `gorm:"size:100" json:"display_name"`
	Description string `gorm:"type:text" json:"description"`
	IsSystem    bool   `gorm:"default:false" json:"is_system"` // 是否为系统角色（不可删除）

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	gorm.Model
	Name        string `gorm:"size:100;uniqueIndex;not null" json:"name"`
	DisplayName string `gorm:"size:100" json:"display_name"`
	Description string `gorm:"type:text" json:"description"`
	Resource    string `gorm:"size:50" json:"resource"` // 资源类型: post, user, media等
	Action      string `gorm:"size:20" json:"action"`   // 操作: create, read, update, delete
	IsSystem    bool   `gorm:"default:false" json:"is_system"`

	// 关联
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// HasPermission 检查角色是否有指定权限
func (r *Role) HasPermission(permissionName string) bool {
	for _, perm := range r.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}
	return false
}

// HasResourcePermission 检查角色是否有资源的指定操作权限
func (r *Role) HasResourcePermission(resource, action string) bool {
	for _, perm := range r.Permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}
	return false
}
