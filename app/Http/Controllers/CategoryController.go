package controllers

import (
	"context"
	"strconv"

	"github.com/chenyusolar/aidecms/internal/app/models"
	"github.com/chenyusolar/aidecms/pkg/database"
	"github.com/chenyusolar/aidecms/pkg/response"
	"github.com/chenyusolar/aidecms/pkg/validator"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gosimple/slug"
)

// CategoryController 分类控制器
type CategoryController struct{}

// NewCategoryController 创建分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=100"`
	Description     string `json:"description"`
	ParentID        *uint  `json:"parent_id"`
	Image           string `json:"image"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

// Create 创建分类
// @Summary      创建分类
// @Description  创建一个新的文章分类
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category body CreateCategoryRequest true "分类信息"
// @Success      201 {object} response.Response{data=models.CategorySwagger}
// @Failure      400 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/categories [post]
func (c *CategoryController) Create(ctx context.Context, hCtx *app.RequestContext) {
	var req CreateCategoryRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	if err := validator.Validate(&req); err != nil {
		if valErr, ok := err.(*validator.ValidationError); ok {
			response.ValidationError(hCtx, valErr.Errors)
			return
		}
		response.BadRequest(hCtx, err.Error())
		return
	}

	category := &models.Category{
		Name:            req.Name,
		Slug:            slug.Make(req.Name),
		Description:     req.Description,
		ParentID:        req.ParentID,
		Image:           req.Image,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}

	db := database.GetDB()
	if err := db.Create(category).Error; err != nil {
		response.ServerError(hCtx, "Failed to create category")
		return
	}

	response.Created(hCtx, category, "Category created successfully")
}

// List 获取分类列表
// @Summary      获取分类列表
// @Description  获取所有分类，支持树形结构
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        tree query bool false "是否返回树形结构" default(false)
// @Success      200 {object} response.Response{data=[]models.CategorySwagger}
// @Failure      500 {object} response.Response
// @Router       /categories [get]
func (c *CategoryController) List(ctx context.Context, hCtx *app.RequestContext) {
	db := database.GetDB()

	var categories []models.Category

	// 是否获取树形结构
	tree := string(hCtx.Query("tree"))

	if tree == "true" {
		// 获取顶级分类及其子分类
		if err := db.Where("parent_id IS NULL").
			Preload("Children").
			Order("sort ASC, created_at DESC").
			Find(&categories).Error; err != nil {
			response.ServerError(hCtx, "Failed to fetch categories")
			return
		}
	} else {
		// 获取所有分类（扁平列表）
		if err := db.Preload("Parent").
			Order("sort ASC, created_at DESC").
			Find(&categories).Error; err != nil {
			response.ServerError(hCtx, "Failed to fetch categories")
			return
		}
	}

	response.Success(hCtx, categories, "")
}

// Get 获取单个分类
func (c *CategoryController) Get(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var category models.Category

	if err := db.Preload("Parent").
		Preload("Children").
		First(&category, id).Error; err != nil {
		response.NotFound(hCtx, "Category not found")
		return
	}

	response.Success(hCtx, category, "")
}

// Update 更新分类
func (c *CategoryController) Update(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var req CreateCategoryRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	db := database.GetDB()
	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		response.NotFound(hCtx, "Category not found")
		return
	}

	category.Name = req.Name
	category.Slug = slug.Make(req.Name)
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.Image = req.Image
	category.MetaTitle = req.MetaTitle
	category.MetaDescription = req.MetaDescription

	if err := db.Save(&category).Error; err != nil {
		response.ServerError(hCtx, "Failed to update category")
		return
	}

	response.Success(hCtx, category, "Category updated successfully")
}

// Delete 删除分类
func (c *CategoryController) Delete(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		response.NotFound(hCtx, "Category not found")
		return
	}

	// 检查是否有子分类
	var childCount int64
	db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		response.BadRequest(hCtx, "Cannot delete category with subcategories")
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		response.ServerError(hCtx, "Failed to delete category")
		return
	}

	response.Success(hCtx, nil, "Category deleted successfully")
}

// TagController 标签控制器
type TagController struct{}

// NewTagController 创建标签控制器
func NewTagController() *TagController {
	return &TagController{}
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

// Create 创建标签
// @Summary      创建标签
// @Description  创建一个新的文章标签
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Param        tag body CreateTagRequest true "标签信息"
// @Success      201 {object} response.Response{data=models.TagSwagger}
// @Failure      400 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/tags [post]
func (c *TagController) Create(ctx context.Context, hCtx *app.RequestContext) {
	var req CreateTagRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	if err := validator.Validate(&req); err != nil {
		if valErr, ok := err.(*validator.ValidationError); ok {
			response.ValidationError(hCtx, valErr.Errors)
			return
		}
		response.BadRequest(hCtx, err.Error())
		return
	}

	tag := &models.Tag{
		Name: req.Name,
		Slug: slug.Make(req.Name),
	}

	db := database.GetDB()
	if err := db.Create(tag).Error; err != nil {
		response.ServerError(hCtx, "Failed to create tag")
		return
	}

	response.Created(hCtx, tag, "Tag created successfully")
}

// List 获取标签列表
// @Summary      获取标签列表
// @Description  分页获取所有标签
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1)
// @Param        per_page query int false "每页数量" default(50)
// @Success      200 {object} response.Response{data=[]models.TagSwagger}
// @Failure      500 {object} response.Response
// @Router       /tags [get]
func (c *TagController) List(ctx context.Context, hCtx *app.RequestContext) {
	db := database.GetDB()

	page, _ := strconv.Atoi(string(hCtx.Query("page")))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(string(hCtx.Query("per_page")))
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}

	var tags []models.Tag
	var total int64

	query := db.Model(&models.Tag{})
	query.Count(&total)

	offset := (page - 1) * perPage
	if err := query.Order("count DESC, name ASC").
		Offset(offset).
		Limit(perPage).
		Find(&tags).Error; err != nil {
		response.ServerError(hCtx, "Failed to fetch tags")
		return
	}

	meta := response.NewMeta(page, perPage, total)
	response.SuccessWithMeta(hCtx, tags, meta, "")
}

// Get 获取单个标签
func (c *TagController) Get(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var tag models.Tag

	if err := db.First(&tag, id).Error; err != nil {
		response.NotFound(hCtx, "Tag not found")
		return
	}

	response.Success(hCtx, tag, "")
}

// Update 更新标签
func (c *TagController) Update(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var req CreateTagRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	db := database.GetDB()
	var tag models.Tag

	if err := db.First(&tag, id).Error; err != nil {
		response.NotFound(hCtx, "Tag not found")
		return
	}

	tag.Name = req.Name
	tag.Slug = slug.Make(req.Name)

	if err := db.Save(&tag).Error; err != nil {
		response.ServerError(hCtx, "Failed to update tag")
		return
	}

	response.Success(hCtx, tag, "Tag updated successfully")
}

// Delete 删除标签
func (c *TagController) Delete(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var tag models.Tag

	if err := db.First(&tag, id).Error; err != nil {
		response.NotFound(hCtx, "Tag not found")
		return
	}

	if err := db.Delete(&tag).Error; err != nil {
		response.ServerError(hCtx, "Failed to delete tag")
		return
	}

	response.Success(hCtx, nil, "Tag deleted successfully")
}
