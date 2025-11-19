package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/clarkgo/clarkgo/internal/app/models"
	"github.com/clarkgo/clarkgo/pkg/database"
	"github.com/clarkgo/clarkgo/pkg/response"
	"github.com/clarkgo/clarkgo/pkg/validator"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gosimple/slug"
)

// PostController 文章控制器
type PostController struct{}

// NewPostController 创建文章控制器
func NewPostController() *PostController {
	return &PostController{}
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title           string `json:"title" validate:"required,min=3,max=200"`
	Content         string `json:"content" validate:"required,min=10"`
	Excerpt         string `json:"excerpt"`
	FeaturedImage   string `json:"featured_image"`
	Status          string `json:"status" validate:"required,oneof=draft published archived"`
	CategoryID      uint   `json:"category_id"`
	Tags            []uint `json:"tags"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title           string `json:"title" validate:"omitempty,min=3,max=200"`
	Content         string `json:"content" validate:"omitempty,min=10"`
	Excerpt         string `json:"excerpt"`
	FeaturedImage   string `json:"featured_image"`
	Status          string `json:"status" validate:"omitempty,oneof=draft published archived"`
	CategoryID      uint   `json:"category_id"`
	Tags            []uint `json:"tags"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
}

// Create 创建文章
// @Summary      创建文章
// @Description  创建一篇新文章
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        post body CreatePostRequest true "文章信息"
// @Success      201 {object} response.Response{data=models.PostSwagger}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/posts [post]
func (c *PostController) Create(ctx context.Context, hCtx *app.RequestContext) {
	var req CreatePostRequest
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

	// 获取当前用户ID（从JWT中间件设置）
	userID, _ := hCtx.Get("user_id")

	// 生成slug
	postSlug := slug.Make(req.Title)

	post := &models.Post{
		Title:           req.Title,
		Slug:            postSlug,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		FeaturedImage:   req.FeaturedImage,
		Status:          req.Status,
		AuthorID:        userID.(uint),
		CategoryID:      req.CategoryID,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		MetaKeywords:    req.MetaKeywords,
	}

	// 如果状态为已发布，设置发布时间
	if req.Status == "published" {
		post.Publish()
	}

	db := database.GetDB()

	// 开始事务
	tx := db.Begin()

	// 创建文章
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		response.ServerError(hCtx, "Failed to create post")
		return
	}

	// 关联标签
	if len(req.Tags) > 0 {
		var tags []models.Tag
		if err := tx.Find(&tags, req.Tags).Error; err != nil {
			tx.Rollback()
			response.ServerError(hCtx, "Failed to fetch tags")
			return
		}
		if err := tx.Model(post).Association("Tags").Append(tags); err != nil {
			tx.Rollback()
			response.ServerError(hCtx, "Failed to associate tags")
			return
		}
	}

	tx.Commit()

	// 预加载关联数据
	db.Preload("Author").Preload("Category").Preload("Tags").First(post, post.ID)

	response.Created(hCtx, post, "Post created successfully")
}

// List 获取文章列表
// @Summary      获取文章列表
// @Description  分页获取文章列表，支持按状态、分类、作者筛选
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1)
// @Param        per_page query int false "每页数量" default(20)
// @Param        status query string false "状态" Enums(draft, published, archived)
// @Param        category_id query int false "分类ID"
// @Param        author_id query int false "作者ID"
// @Success      200 {object} response.Response{data=[]models.PostSwagger}
// @Failure      500 {object} response.Response
// @Router       /posts [get]
func (c *PostController) List(ctx context.Context, hCtx *app.RequestContext) {
	db := database.GetDB()

	// 分页参数
	page, _ := strconv.Atoi(string(hCtx.Query("page")))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(string(hCtx.Query("per_page")))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// 过滤参数
	status := string(hCtx.Query("status"))
	categoryID := string(hCtx.Query("category_id"))
	authorID := string(hCtx.Query("author_id"))

	var posts []models.Post
	query := db.Model(&models.Post{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if authorID != "" {
		query = query.Where("author_id = ?", authorID)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * perPage
	if err := query.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&posts).Error; err != nil {
		response.ServerError(hCtx, "Failed to fetch posts")
		return
	}

	meta := response.NewMeta(page, perPage, total)
	response.SuccessWithMeta(hCtx, posts, meta, "")
}

// Get 获取单篇文章
// @Summary      获取文章详情
// @Description  根据ID获取文章详细信息
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id path int true "文章ID"
// @Success      200 {object} response.Response{data=models.PostSwagger}
// @Failure      404 {object} response.Response
// @Router       /posts/{id} [get]
func (c *PostController) Get(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var post models.Post

	if err := db.Preload("Author").
		Preload("Category").
		Preload("Tags").
		First(&post, id).Error; err != nil {
		response.NotFound(hCtx, "Post not found")
		return
	}

	// 增加浏览次数
	db.Model(&post).UpdateColumn("view_count", post.ViewCount+1)
	post.ViewCount++

	response.Success(hCtx, post, "")
}

// Update 更新文章
// @Summary      更新文章
// @Description  更新已存在的文章
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id path int true "文章ID"
// @Param        post body UpdatePostRequest true "文章信息"
// @Success      200 {object} response.Response{data=models.PostSwagger}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      422 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/posts/{id} [put]
func (c *PostController) Update(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var req UpdatePostRequest
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
	var post models.Post

	if err := db.First(&post, id).Error; err != nil {
		response.NotFound(hCtx, "Post not found")
		return
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
		post.Slug = slug.Make(req.Title)
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	post.Excerpt = req.Excerpt
	post.FeaturedImage = req.FeaturedImage
	if req.Status != "" {
		post.Status = req.Status
		// 如果从草稿变为发布，设置发布时间
		if req.Status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}
	if req.CategoryID != 0 {
		post.CategoryID = req.CategoryID
	}
	post.MetaTitle = req.MetaTitle
	post.MetaDescription = req.MetaDescription
	post.MetaKeywords = req.MetaKeywords

	tx := db.Begin()

	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		response.ServerError(hCtx, "Failed to update post")
		return
	}

	// 更新标签
	if req.Tags != nil {
		var tags []models.Tag
		if err := tx.Find(&tags, req.Tags).Error; err != nil {
			tx.Rollback()
			response.ServerError(hCtx, "Failed to fetch tags")
			return
		}
		if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			response.ServerError(hCtx, "Failed to update tags")
			return
		}
	}

	tx.Commit()

	// 重新加载
	db.Preload("Author").Preload("Category").Preload("Tags").First(&post, post.ID)

	response.Success(hCtx, post, "Post updated successfully")
}

// Delete 删除文章
// @Summary      删除文章
// @Description  软删除指定的文章
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id path int true "文章ID"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/posts/{id} [delete]
func (c *PostController) Delete(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var post models.Post

	if err := db.First(&post, id).Error; err != nil {
		response.NotFound(hCtx, "Post not found")
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		response.ServerError(hCtx, "Failed to delete post")
		return
	}

	response.Success(hCtx, nil, "Post deleted successfully")
}

// Publish 发布文章
// @Summary      发布文章
// @Description  将草稿文章发布上线
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id path int true "文章ID"
// @Success      200 {object} response.Response{data=models.PostSwagger}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/posts/{id}/publish [post]
func (c *PostController) Publish(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var post models.Post

	if err := db.First(&post, id).Error; err != nil {
		response.NotFound(hCtx, "Post not found")
		return
	}

	post.Publish()

	if err := db.Save(&post).Error; err != nil {
		response.ServerError(hCtx, "Failed to publish post")
		return
	}

	response.Success(hCtx, post, "Post published successfully")
}
