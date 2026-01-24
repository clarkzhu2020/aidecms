package services

import (
	"errors"
	"fmt"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService 创建权限服务实例
func NewPermissionService() *PermissionService {
	return &PermissionService{
		db: database.GetDB(),
	}
}

// AssignRoleToUser 为用户分配角色
func (s *PermissionService) AssignRoleToUser(userID, roleID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// 添加角色
	if err := s.db.Model(&user).Association("Roles").Append(&role); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RemoveRoleFromUser 移除用户的角色
func (s *PermissionService) RemoveRoleFromUser(userID, roleID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	if err := s.db.Model(&user).Association("Roles").Delete(&role); err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	return nil
}

// AssignPermissionToRole 为角色分配权限
func (s *PermissionService) AssignPermissionToRole(roleID, permissionID uint) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	if err := s.db.Model(&role).Association("Permissions").Append(&permission); err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}

	return nil
}

// RemovePermissionFromRole 移除角色的权限
func (s *PermissionService) RemovePermissionFromRole(roleID, permissionID uint) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	if err := s.db.Model(&role).Association("Permissions").Delete(&permission); err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}

	return nil
}

// GetUserPermissions 获取用户的所有权限
func (s *PermissionService) GetUserPermissions(userID uint) ([]models.Permission, error) {
	var user models.User
	if err := s.db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 收集所有权限（去重）
	permMap := make(map[uint]models.Permission)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permMap[perm.ID] = perm
		}
	}

	permissions := make([]models.Permission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// CheckUserPermission 检查用户是否有指定权限
func (s *PermissionService) CheckUserPermission(userID uint, permissionName string) (bool, error) {
	var user models.User
	if err := s.db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return user.HasPermission(permissionName), nil
}

// CheckUserResourcePermission 检查用户是否有资源的指定操作权限
func (s *PermissionService) CheckUserResourcePermission(userID uint, resource, action string) (bool, error) {
	var user models.User
	if err := s.db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return user.HasResourcePermission(resource, action), nil
}

// CreateRole 创建角色
func (s *PermissionService) CreateRole(role *models.Role) error {
	return s.db.Create(role).Error
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(permission *models.Permission) error {
	return s.db.Create(permission).Error
}

// GetRoleByName 根据名称获取角色
func (s *PermissionService) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	if err := s.db.Where("name = ?", name).Preload("Permissions").First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetPermissionByName 根据名称获取权限
func (s *PermissionService) GetPermissionByName(name string) (*models.Permission, error) {
	var permission models.Permission
	if err := s.db.Where("name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// SyncUserRoles 同步用户角色（替换所有角色）
func (s *PermissionService) SyncUserRoles(userID uint, roleIDs []uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var roles []models.Role
	if err := s.db.Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to fetch roles: %w", err)
	}

	if err := s.db.Model(&user).Association("Roles").Replace(roles); err != nil {
		return fmt.Errorf("failed to sync roles: %w", err)
	}

	return nil
}

// SyncRolePermissions 同步角色权限（替换所有权限）
func (s *PermissionService) SyncRolePermissions(roleID uint, permissionIDs []uint) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	var permissions []models.Permission
	if err := s.db.Find(&permissions, permissionIDs).Error; err != nil {
		return fmt.Errorf("failed to fetch permissions: %w", err)
	}

	if err := s.db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
		return fmt.Errorf("failed to sync permissions: %w", err)
	}

	return nil
}
