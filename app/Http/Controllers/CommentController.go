package controllers

import (
	"context"
	"strconv"
	"strings"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/clarkzhu2020/aidecms/pkg/response"
	"github.com/clarkzhu2020/aidecms/pkg/validator"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// CommentController 评论控制器
type CommentController struct{}

// NewCommentController 创建评论控制器
func NewCommentController() *CommentController {
	return &CommentController{}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID      uint   `json:"post_id" validate:"required"`
	Content     string `json:"content" validate:"required,min=5,max=1000"`
	ParentID    uint   `json:"parent_id"`
	AuthorName  string `json:"author_name" validate:"omitempty,min=2,max=100"`
	AuthorEmail string `json:"author_email" validate:"omitempty,email"`
	AuthorURL   string `json:"author_url" validate:"omitempty,url"`
	Rating      int    `json:"rating" validate:"omitempty,min=1,max=5"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" validate:"omitempty,min=5,max=1000"`
	Status  string `json:"status" validate:"omitempty,oneof=pending approved spam trash"`
}

// Create 创建评论
// @Summary      创建评论
// @Description  为文章创建新评论，支持嵌套回复
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        comment body CreateCommentRequest true "评论信息"
// @Success      201 {object} response.Response{data=models.CommentSwagger}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /comments [post]
func (c *CommentController) Create(ctx context.Context, hCtx *app.RequestContext) {
	var req CreateCommentRequest
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

	// 检查文章是否存在
	var post models.Post
	if err := db.First(&post, req.PostID).Error; err != nil {
		response.NotFound(hCtx, "Post not found")
		return
	}

	// 检查父评论是否存在（如果是回复）
	if req.ParentID > 0 {
		var parentComment models.Comment
		if err := db.First(&parentComment, req.ParentID).Error; err != nil {
			response.NotFound(hCtx, "Parent comment not found")
			return
		}
	}

	// 获取用户信息（如果已登录）
	var userID uint
	if uid, exists := hCtx.Get("user_id"); exists {
		userID = uid.(uint)
	}

	// 获取客户端信息
	clientIP := string(hCtx.ClientIP())
	userAgent := string(hCtx.UserAgent())

	comment := &models.Comment{
		PostID:      req.PostID,
		UserID:      userID,
		ParentID:    req.ParentID,
		Content:     req.Content,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		AuthorURL:   req.AuthorURL,
		AuthorIP:    clientIP,
		UserAgent:   userAgent,
		Status:      "pending", // 默认待审核
		Rating:      req.Rating,
		IsAnonymous: req.IsAnonymous,
	}

	// 简单的垃圾评论检测
	if isSpam(comment) {
		comment.Status = "spam"
	}

	if err := db.Create(comment).Error; err != nil {
		response.ServerError(hCtx, "Failed to create comment")
		return
	}

	// 预加载关联数据
	db.Preload("User").Preload("Parent").First(comment, comment.ID)

	response.Created(hCtx, comment, "Comment created successfully")
}

// List 获取评论列表
// @Summary      获取评论列表
// @Description  分页获取评论列表，支持按文章、状态筛选
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        post_id query int false "文章ID"
// @Param        status query string false "状态" Enums(pending, approved, spam, trash)
// @Param        page query int false "页码" default(1)
// @Param        per_page query int false "每页数量" default(20)
// @Param        tree query bool false "是否返回树形结构" default(false)
// @Success      200 {object} response.Response{data=[]models.CommentSwagger}
// @Failure      500 {object} response.Response
// @Router       /comments [get]
func (c *CommentController) List(ctx context.Context, hCtx *app.RequestContext) {
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
	postID := string(hCtx.Query("post_id"))
	status := string(hCtx.Query("status"))
	tree := string(hCtx.Query("tree")) == "true"

	var comments []models.Comment
	query := db.Model(&models.Comment{})

	if postID != "" {
		query = query.Where("post_id = ?", postID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// 默认只显示已批准的评论
		query = query.Where("status = ?", "approved")
	}

	// 计算总数
	var total int64
	query.Count(&total)

	if tree {
		// 返回树形结构（只取顶级评论）
		query = query.Where("parent_id = 0")
		offset := (page - 1) * perPage

		if err := query.Offset(offset).Limit(perPage).
			Order("created_at DESC").
			Preload("User").
			Preload("Children", "status = ?", "approved").
			Find(&comments).Error; err != nil {
			response.ServerError(hCtx, "Failed to fetch comments")
			return
		}

		// 递归加载子评论
		for i := range comments {
			loadCommentChildren(db, &comments[i])
		}
	} else {
		// 返回扁平列表
		offset := (page - 1) * perPage

		if err := query.Offset(offset).Limit(perPage).
			Order("created_at DESC").
			Preload("User").
			Preload("Parent").
			Find(&comments).Error; err != nil {
			response.ServerError(hCtx, "Failed to fetch comments")
			return
		}
	}

	// 分页响应
	response.SuccessWithMeta(hCtx, comments, response.NewMeta(page, perPage, total), "Comments fetched successfully")
}

// loadCommentChildren 递归加载子评论
func loadCommentChildren(db *gorm.DB, comment *models.Comment) {
	if len(comment.Children) > 0 {
		for i := range comment.Children {
			db.Where("parent_id = ? AND status = ?", comment.Children[i].ID, "approved").
				Order("created_at ASC").
				Preload("User").
				Find(&comment.Children[i].Children)
			loadCommentChildren(db, &comment.Children[i])
		}
	}
}

// Get 获取单个评论
// @Summary      获取评论详情
// @Description  根据ID获取评论详细信息
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id path int true "评论ID"
// @Success      200 {object} response.Response{data=models.CommentSwagger}
// @Failure      404 {object} response.Response
// @Router       /comments/{id} [get]
func (c *CommentController) Get(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var comment models.Comment
	db := database.GetDB()

	if err := db.Preload("User").Preload("Parent").Preload("Children").
		First(&comment, id).Error; err != nil {
		response.NotFound(hCtx, "Comment not found")
		return
	}

	response.Success(hCtx, comment, "Comment fetched successfully")
}

// Update 更新评论
// @Summary      更新评论
// @Description  更新已存在的评论
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id path int true "评论ID"
// @Param        comment body UpdateCommentRequest true "评论信息"
// @Success      200 {object} response.Response{data=models.CommentSwagger}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/comments/{id} [put]
func (c *CommentController) Update(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var req UpdateCommentRequest
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
	var comment models.Comment
	if err := db.First(&comment, id).Error; err != nil {
		response.NotFound(hCtx, "Comment not found")
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := db.Model(&comment).Updates(updates).Error; err != nil {
		response.ServerError(hCtx, "Failed to update comment")
		return
	}

	// 重新加载评论
	db.Preload("User").Preload("Parent").First(&comment, id)

	response.Success(hCtx, comment, "Comment updated successfully")
}

// Delete 删除评论
// @Summary      删除评论
// @Description  删除指定的评论（软删除）
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id path int true "评论ID"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/comments/{id} [delete]
func (c *CommentController) Delete(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var comment models.Comment
	if err := db.First(&comment, id).Error; err != nil {
		response.NotFound(hCtx, "Comment not found")
		return
	}

	// 软删除或标记为trash
	if err := db.Model(&comment).Update("status", "trash").Error; err != nil {
		response.ServerError(hCtx, "Failed to delete comment")
		return
	}

	response.Success(hCtx, nil, "Comment deleted successfully")
}

// Approve 批准评论
// @Summary      批准评论
// @Description  将待审核评论标记为已批准
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id path int true "评论ID"
// @Success      200 {object} response.Response{data=models.CommentSwagger}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/comments/{id}/approve [post]
func (c *CommentController) Approve(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var comment models.Comment
	if err := db.First(&comment, id).Error; err != nil {
		response.NotFound(hCtx, "Comment not found")
		return
	}

	comment.Approve()
	if err := db.Save(&comment).Error; err != nil {
		response.ServerError(hCtx, "Failed to approve comment")
		return
	}

	response.Success(hCtx, comment, "Comment approved successfully")
}

// MarkAsSpam 标记为垃圾评论
// @Summary      标记垃圾评论
// @Description  将评论标记为垃圾评论
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id path int true "评论ID"
// @Success      200 {object} response.Response{data=models.CommentSwagger}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/comments/{id}/spam [post]
func (c *CommentController) MarkAsSpam(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var comment models.Comment
	if err := db.First(&comment, id).Error; err != nil {
		response.NotFound(hCtx, "Comment not found")
		return
	}

	comment.MarkAsSpam()
	if err := db.Save(&comment).Error; err != nil {
		response.ServerError(hCtx, "Failed to mark comment as spam")
		return
	}

	response.Success(hCtx, comment, "Comment marked as spam")
}

// isSpam 简单的垃圾评论检测
func isSpam(comment *models.Comment) bool {
	content := strings.ToLower(comment.Content)

	// 垃圾关键词列表
	spamKeywords := []string{
		"viagra", "casino", "lottery", "prize", "winner",
		"click here", "buy now", "limited time",
		"http://", "https://", // 包含过多链接
	}

	// 检查垃圾关键词
	for _, keyword := range spamKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}

	// 检查是否包含过多链接
	linkCount := strings.Count(content, "http://") + strings.Count(content, "https://")
	if linkCount > 2 {
		return true
	}

	// 检查是否全是大写（超过50%）
	upperCount := 0
	totalLetters := 0
	for _, char := range comment.Content {
		if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
			totalLetters++
			if char >= 'A' && char <= 'Z' {
				upperCount++
			}
		}
	}
	if totalLetters > 0 && float64(upperCount)/float64(totalLetters) > 0.5 {
		return true
	}

	return false
}
