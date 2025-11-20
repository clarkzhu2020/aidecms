package commands

import (
	"fmt"

	"github.com/chenyusolar/aidecms/database/migrations"
	"github.com/chenyusolar/aidecms/internal/app/models"
	"github.com/chenyusolar/aidecms/internal/app/services"
	"github.com/chenyusolar/aidecms/pkg/database"
	"github.com/spf13/cobra"
)

var cmsInitCmd = &cobra.Command{
	Use:   "cms:init",
	Short: "Initialize CMS tables and default data",
	Long:  "Create all CMS related tables and insert default roles and permissions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing CMS...")

		// 获取数据库连接
		db := database.GetDB()
		if db == nil {
			fmt.Println("Error: Database connection not available")
			return
		}

		dbInstance := &database.Database{DB: db}

		// 运行迁移
		fmt.Println("Creating CMS tables...")
		if err := migrations.CreateCMSTables(dbInstance); err != nil {
			fmt.Printf("Error creating tables: %v\n", err)
			return
		}
		fmt.Println("✓ CMS tables created successfully")

		// 创建默认角色和权限
		fmt.Println("\nCreating default roles and permissions...")
		if err := createDefaultRolesAndPermissions(); err != nil {
			fmt.Printf("Error creating roles and permissions: %v\n", err)
			return
		}
		fmt.Println("✓ Default roles and permissions created successfully")

		fmt.Println("\n✓ CMS initialization completed successfully!")
		fmt.Println("\nDefault roles created:")
		fmt.Println("  - super_admin: Full system access")
		fmt.Println("  - admin: Administrative access")
		fmt.Println("  - editor: Content management")
		fmt.Println("  - author: Create and manage own posts")
		fmt.Println("  - user: Basic user access")
	},
}

func createDefaultRolesAndPermissions() error {
	permService := services.NewPermissionService()

	// 定义角色
	roles := []models.Role{
		{
			Name:        "super_admin",
			DisplayName: "Super Administrator",
			Description: "Full system access with all permissions",
			IsSystem:    true,
		},
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Administrative access to manage content and users",
			IsSystem:    true,
		},
		{
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Can manage all content",
			IsSystem:    true,
		},
		{
			Name:        "author",
			DisplayName: "Author",
			Description: "Can create and manage own posts",
			IsSystem:    true,
		},
		{
			Name:        "user",
			DisplayName: "User",
			Description: "Basic user with read access",
			IsSystem:    true,
		},
	}

	// 定义权限
	permissions := []models.Permission{
		// 文章权限
		{Name: "post.create", DisplayName: "Create Post", Resource: "post", Action: "create", IsSystem: true},
		{Name: "post.read", DisplayName: "Read Post", Resource: "post", Action: "read", IsSystem: true},
		{Name: "post.update", DisplayName: "Update Post", Resource: "post", Action: "update", IsSystem: true},
		{Name: "post.delete", DisplayName: "Delete Post", Resource: "post", Action: "delete", IsSystem: true},
		{Name: "post.publish", DisplayName: "Publish Post", Resource: "post", Action: "publish", IsSystem: true},

		// 分类权限
		{Name: "category.create", DisplayName: "Create Category", Resource: "category", Action: "create", IsSystem: true},
		{Name: "category.read", DisplayName: "Read Category", Resource: "category", Action: "read", IsSystem: true},
		{Name: "category.update", DisplayName: "Update Category", Resource: "category", Action: "update", IsSystem: true},
		{Name: "category.delete", DisplayName: "Delete Category", Resource: "category", Action: "delete", IsSystem: true},

		// 标签权限
		{Name: "tag.create", DisplayName: "Create Tag", Resource: "tag", Action: "create", IsSystem: true},
		{Name: "tag.read", DisplayName: "Read Tag", Resource: "tag", Action: "read", IsSystem: true},
		{Name: "tag.update", DisplayName: "Update Tag", Resource: "tag", Action: "update", IsSystem: true},
		{Name: "tag.delete", DisplayName: "Delete Tag", Resource: "tag", Action: "delete", IsSystem: true},

		// 媒体权限
		{Name: "media.upload", DisplayName: "Upload Media", Resource: "media", Action: "upload", IsSystem: true},
		{Name: "media.read", DisplayName: "Read Media", Resource: "media", Action: "read", IsSystem: true},
		{Name: "media.update", DisplayName: "Update Media", Resource: "media", Action: "update", IsSystem: true},
		{Name: "media.delete", DisplayName: "Delete Media", Resource: "media", Action: "delete", IsSystem: true},

		// 用户管理权限
		{Name: "user.create", DisplayName: "Create User", Resource: "user", Action: "create", IsSystem: true},
		{Name: "user.read", DisplayName: "Read User", Resource: "user", Action: "read", IsSystem: true},
		{Name: "user.update", DisplayName: "Update User", Resource: "user", Action: "update", IsSystem: true},
		{Name: "user.delete", DisplayName: "Delete User", Resource: "user", Action: "delete", IsSystem: true},

		// 角色权限管理
		{Name: "role.manage", DisplayName: "Manage Roles", Resource: "role", Action: "manage", IsSystem: true},
		{Name: "permission.manage", DisplayName: "Manage Permissions", Resource: "permission", Action: "manage", IsSystem: true},
	}

	// 创建角色
	for i := range roles {
		if err := permService.CreateRole(&roles[i]); err != nil {
			// 忽略重复键错误
			fmt.Printf("Note: Role %s already exists\n", roles[i].Name)
		}
	}

	// 创建权限
	for i := range permissions {
		if err := permService.CreatePermission(&permissions[i]); err != nil {
			fmt.Printf("Note: Permission %s already exists\n", permissions[i].Name)
		}
	}

	// 为角色分配权限
	// Super Admin - 所有权限
	superAdmin, _ := permService.GetRoleByName("super_admin")
	if superAdmin != nil {
		var allPermIDs []uint
		for _, perm := range permissions {
			p, _ := permService.GetPermissionByName(perm.Name)
			if p != nil {
				allPermIDs = append(allPermIDs, p.ID)
			}
		}
		permService.SyncRolePermissions(superAdmin.ID, allPermIDs)
	}

	// Admin - 内容和用户管理权限
	admin, _ := permService.GetRoleByName("admin")
	if admin != nil {
		adminPerms := []string{
			"post.create", "post.read", "post.update", "post.delete", "post.publish",
			"category.create", "category.read", "category.update", "category.delete",
			"tag.create", "tag.read", "tag.update", "tag.delete",
			"media.upload", "media.read", "media.update", "media.delete",
			"user.read", "user.update",
		}
		var permIDs []uint
		for _, permName := range adminPerms {
			p, _ := permService.GetPermissionByName(permName)
			if p != nil {
				permIDs = append(permIDs, p.ID)
			}
		}
		permService.SyncRolePermissions(admin.ID, permIDs)
	}

	// Editor - 内容管理权限
	editor, _ := permService.GetRoleByName("editor")
	if editor != nil {
		editorPerms := []string{
			"post.create", "post.read", "post.update", "post.delete", "post.publish",
			"category.create", "category.read", "category.update",
			"tag.create", "tag.read", "tag.update",
			"media.upload", "media.read", "media.update",
		}
		var permIDs []uint
		for _, permName := range editorPerms {
			p, _ := permService.GetPermissionByName(permName)
			if p != nil {
				permIDs = append(permIDs, p.ID)
			}
		}
		permService.SyncRolePermissions(editor.ID, permIDs)
	}

	// Author - 创建和管理自己的内容
	author, _ := permService.GetRoleByName("author")
	if author != nil {
		authorPerms := []string{
			"post.create", "post.read", "post.update",
			"category.read", "tag.read", "tag.create",
			"media.upload", "media.read",
		}
		var permIDs []uint
		for _, permName := range authorPerms {
			p, _ := permService.GetPermissionByName(permName)
			if p != nil {
				permIDs = append(permIDs, p.ID)
			}
		}
		permService.SyncRolePermissions(author.ID, permIDs)
	}

	// User - 基本读取权限
	user, _ := permService.GetRoleByName("user")
	if user != nil {
		userPerms := []string{
			"post.read", "category.read", "tag.read", "media.read",
		}
		var permIDs []uint
		for _, permName := range userPerms {
			p, _ := permService.GetPermissionByName(permName)
			if p != nil {
				permIDs = append(permIDs, p.ID)
			}
		}
		permService.SyncRolePermissions(user.ID, permIDs)
	}

	return nil
}

// func init() {
// 	rootCmd.AddCommand(cmsInitCmd)
// }
