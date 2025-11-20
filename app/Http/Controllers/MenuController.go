package controllers

import (
	"context"

	"github.com/chenyusolar/aidecms/internal/app/models"
	"github.com/chenyusolar/aidecms/pkg/database"
	"github.com/chenyusolar/aidecms/pkg/response"
	"github.com/chenyusolar/aidecms/pkg/validator"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// MenuController 菜单控制器
type MenuController struct{}

// NewMenuController 创建菜单控制器
func NewMenuController() *MenuController {
	return &MenuController{}
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Title       string `json:"title" validate:"required,min=2,max=200"`
	URL         string `json:"url" validate:"required,max=500"`
	Icon        string `json:"icon"`
	Target      string `json:"target" validate:"omitempty,oneof=_self _blank"`
	Position    string `json:"position" validate:"required,oneof=header footer sidebar mobile"`
	ParentID    uint   `json:"parent_id"`
	Sort        int    `json:"sort"`
	IsActive    bool   `json:"is_active"`
	CSSClass    string `json:"css_class"`
	Description string `json:"description"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Title       string `json:"title" validate:"omitempty,min=2,max=200"`
	URL         string `json:"url" validate:"omitempty,max=500"`
	Icon        string `json:"icon"`
	Target      string `json:"target" validate:"omitempty,oneof=_self _blank"`
	Position    string `json:"position" validate:"omitempty,oneof=header footer sidebar mobile"`
	ParentID    uint   `json:"parent_id"`
	Sort        int    `json:"sort"`
	IsActive    *bool  `json:"is_active"`
	CSSClass    string `json:"css_class"`
	Description string `json:"description"`
}

// Create 创建菜单
// @Summary      创建菜单
// @Description  创建一个新的菜单项
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        menu body CreateMenuRequest true "菜单信息"
// @Success      201 {object} response.Response{data=models.MenuSwagger}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/menus [post]
func (c *MenuController) Create(ctx context.Context, hCtx *app.RequestContext) {
	var req CreateMenuRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 验证请求
	if err := validator.Validate(&req); err != nil {
		if valErr, ok := err.(*validator.ValidationError); ok {
			response.ValidationError(hCtx, valErr.Errors)
			return
		}
		response.BadRequest(hCtx, err.Error())
		return
	}

	// 设置默认值
	if req.Target == "" {
		req.Target = "_self"
	}
	if !req.IsActive {
		req.IsActive = true
	}

	menu := &models.Menu{
		Name:        req.Name,
		Title:       req.Title,
		URL:         req.URL,
		Icon:        req.Icon,
		Target:      req.Target,
		Position:    req.Position,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
		IsActive:    req.IsActive,
		CSSClass:    req.CSSClass,
		Description: req.Description,
	}

	db := database.GetDB()
	if err := db.Create(menu).Error; err != nil {
		response.ServerError(hCtx, "Failed to create menu")
		return
	}

	response.Created(hCtx, menu, "Menu created successfully")
}

// List 获取菜单列表
// @Summary      获取菜单列表
// @Description  获取菜单列表，支持按位置筛选和树形结构
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        position query string false "位置" Enums(header, footer, sidebar, mobile)
// @Param        tree query bool false "是否返回树形结构" default(false)
// @Success      200 {object} response.Response{data=[]models.MenuSwagger}
// @Failure      500 {object} response.Response
// @Router       /menus [get]
func (c *MenuController) List(ctx context.Context, hCtx *app.RequestContext) {
	db := database.GetDB()

	position := string(hCtx.Query("position"))
	tree := string(hCtx.Query("tree")) == "true"

	var menus []models.Menu
	query := db.Model(&models.Menu{})

	if position != "" {
		query = query.Where("position = ?", position)
	}

	// 只查询激活的菜单
	query = query.Where("is_active = ?", true)

	if tree {
		// 返回树形结构
		query = query.Where("parent_id = 0").Order("sort ASC, id ASC")
		if err := query.Preload("Children", "is_active = ?", true).Find(&menus).Error; err != nil {
			response.ServerError(hCtx, "Failed to fetch menus")
			return
		}

		// 递归加载子菜单
		for i := range menus {
			loadChildren(db, &menus[i])
		}
	} else {
		// 返回扁平列表
		query = query.Order("position ASC, sort ASC, id ASC")
		if err := query.Find(&menus).Error; err != nil {
			response.ServerError(hCtx, "Failed to fetch menus")
			return
		}
	}

	response.Success(hCtx, menus, "Menus fetched successfully")
}

// loadChildren 递归加载子菜单
func loadChildren(db *gorm.DB, menu *models.Menu) {
	if len(menu.Children) > 0 {
		for i := range menu.Children {
			db.Where("parent_id = ? AND is_active = ?", menu.Children[i].ID, true).
				Order("sort ASC, id ASC").
				Find(&menu.Children[i].Children)
			loadChildren(db, &menu.Children[i])
		}
	}
}

// Get 获取单个菜单
// @Summary      获取菜单详情
// @Description  根据ID获取菜单详细信息
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id path int true "菜单ID"
// @Success      200 {object} response.Response{data=models.MenuSwagger}
// @Failure      404 {object} response.Response
// @Router       /menus/{id} [get]
func (c *MenuController) Get(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var menu models.Menu
	db := database.GetDB()

	if err := db.Preload("Children").Preload("Parent").First(&menu, id).Error; err != nil {
		response.NotFound(hCtx, "Menu not found")
		return
	}

	response.Success(hCtx, menu, "Menu fetched successfully")
}

// Update 更新菜单
// @Summary      更新菜单
// @Description  更新已存在的菜单
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id path int true "菜单ID"
// @Param        menu body UpdateMenuRequest true "菜单信息"
// @Success      200 {object} response.Response{data=models.MenuSwagger}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      422 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/menus/{id} [put]
func (c *MenuController) Update(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var req UpdateMenuRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 验证请求
	if err := validator.Validate(&req); err != nil {
		if valErr, ok := err.(*validator.ValidationError); ok {
			response.ValidationError(hCtx, valErr.Errors)
			return
		}
		response.BadRequest(hCtx, err.Error())
		return
	}

	db := database.GetDB()
	var menu models.Menu
	if err := db.First(&menu, id).Error; err != nil {
		response.NotFound(hCtx, "Menu not found")
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Target != "" {
		updates["target"] = req.Target
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.ParentID > 0 {
		updates["parent_id"] = req.ParentID
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.CSSClass != "" {
		updates["css_class"] = req.CSSClass
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := db.Model(&menu).Updates(updates).Error; err != nil {
		response.ServerError(hCtx, "Failed to update menu")
		return
	}

	// 重新加载菜单
	db.Preload("Children").Preload("Parent").First(&menu, id)

	response.Success(hCtx, menu, "Menu updated successfully")
}

// Delete 删除菜单
// @Summary      删除菜单
// @Description  删除指定的菜单（软删除）
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id path int true "菜单ID"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/menus/{id} [delete]
func (c *MenuController) Delete(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var menu models.Menu
	if err := db.First(&menu, id).Error; err != nil {
		response.NotFound(hCtx, "Menu not found")
		return
	}

	// 检查是否有子菜单
	var childCount int64
	db.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		response.BadRequest(hCtx, "Cannot delete menu with children")
		return
	}

	if err := db.Delete(&menu).Error; err != nil {
		response.ServerError(hCtx, "Failed to delete menu")
		return
	}

	response.Success(hCtx, nil, "Menu deleted successfully")
}

// Reorder 重新排序菜单
// @Summary      重新排序菜单
// @Description  批量更新菜单排序
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        orders body []map[string]interface{} true "排序数据 [{\"id\": 1, \"sort\": 1}]"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/menus/reorder [post]
func (c *MenuController) Reorder(ctx context.Context, hCtx *app.RequestContext) {
	var orders []map[string]interface{}
	if err := hCtx.BindJSON(&orders); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	db := database.GetDB()
	tx := db.Begin()

	for _, order := range orders {
		id, ok1 := order["id"]
		sort, ok2 := order["sort"]

		if !ok1 || !ok2 {
			tx.Rollback()
			response.BadRequest(hCtx, "Invalid order data")
			return
		}

		if err := tx.Model(&models.Menu{}).Where("id = ?", id).Update("sort", sort).Error; err != nil {
			tx.Rollback()
			response.ServerError(hCtx, "Failed to update menu order")
			return
		}
	}

	tx.Commit()
	response.Success(hCtx, nil, "Menu reordered successfully")
}
