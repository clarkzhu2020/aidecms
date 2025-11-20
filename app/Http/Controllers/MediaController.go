package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chenyusolar/aidecms/config"
	"github.com/chenyusolar/aidecms/internal/app/models"
	"github.com/chenyusolar/aidecms/pkg/database"
	"github.com/chenyusolar/aidecms/pkg/upload"
	"github.com/cloudwego/hertz/pkg/app"
)

// MediaController 媒体控制器
type MediaController struct {
	uploader       *upload.Uploader
	imageProcessor *upload.ImageProcessor
}

// NewMediaController 创建媒体控制器
func NewMediaController() *MediaController {
	// 从配置获取存储实例（支持 local/oss/s3）
	storage, err := config.GetStorage()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize storage: %v", err))
	}

	// 初始化上传器
	uploader := upload.NewUploader(&upload.UploadConfig{
		MaxSize: 10 * 1024 * 1024, // 10MB
		AllowedExts: []string{
			".jpg", ".jpeg", ".png", ".gif", ".webp", // 图片
			".pdf", ".doc", ".docx", ".xls", ".xlsx", // 文档
			".zip", ".rar", // 压缩包
		},
		Storage: storage,
	})

	// 初始化图片处理器
	imageProcessor := upload.NewImageProcessor(storage)

	return &MediaController{
		uploader:       uploader,
		imageProcessor: imageProcessor,
	}
}

// Upload 上传文件
// @Summary      上传媒体文件
// @Description  上传一个或多个文件，支持图片、文档、压缩包等
// @Tags         Media
// @Accept       multipart/form-data
// @Produce      json
// @Param        files formData file true "文件(可多选)"
// @Success      201 {object} response.Response{data=[]models.MediaSwagger}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/media/upload [post]
func (c *MediaController) Upload(ctx context.Context, hCtx *app.RequestContext) {
	// 获取上传的文件
	form, err := hCtx.MultipartForm()
	if err != nil {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Failed to parse form data",
			"message": err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "No files uploaded",
		})
		return
	}

	// 上传文件
	results, err := c.uploader.UploadMultiple(files)
	if err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to upload files",
			"message": err.Error(),
		})
		return
	}

	db := database.GetDB()
	mediaRecords := make([]*models.Media, 0, len(results))

	// 保存到数据库
	for _, result := range results {
		media := &models.Media{
			FileName:     result.FileName,
			OriginalName: result.OriginalName,
			FilePath:     result.Path,
			FileURL:      result.URL,
			FileSize:     result.Size,
			MimeType:     result.MimeType,
			Extension:    result.Extension,
			Hash:         result.Hash,
			FileType:     models.GetFileType(result.MimeType),
		}

		// 如果是图片，处理缩略图
		if upload.IsImage(result.OriginalName) {
			thumbnailSizes := []upload.ThumbnailSize{
				{Name: "small", Width: 150, Height: 150},
				{Name: "medium", Width: 300, Height: 300},
				{Name: "large", Width: 800, Height: 800},
			}

			imageResult, err := c.imageProcessor.ProcessImage(result.Path, thumbnailSizes)
			if err == nil {
				thumbnailsJSON, _ := json.Marshal(imageResult.Thumbnails)
				media.Thumbnails = string(thumbnailsJSON)

				// 获取图片尺寸
				info, err := c.imageProcessor.GetImageInfo(result.Path)
				if err == nil {
					if width, ok := info["width"].(int); ok {
						media.Width = width
					}
					if height, ok := info["height"].(int); ok {
						media.Height = height
					}
				}
			}
		}

		if err := db.Create(media).Error; err != nil {
			hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   "Failed to save media record",
				"message": err.Error(),
			})
			return
		}

		mediaRecords = append(mediaRecords, media)
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    mediaRecords,
		"message": fmt.Sprintf("Successfully uploaded %d file(s)", len(mediaRecords)),
	})
}

// List 获取媒体列表
func (c *MediaController) List(ctx context.Context, hCtx *app.RequestContext) {
	db := database.GetDB()

	// 分页参数
	page, _ := strconv.Atoi(hCtx.Query("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(hCtx.Query("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// 文件类型过滤
	fileType := hCtx.Query("file_type")

	var media []models.Media
	query := db.Model(&models.Media{})

	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&media).Error; err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch media",
			"message": err.Error(),
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    media,
		"meta": map[string]interface{}{
			"current_page": page,
			"per_page":     perPage,
			"total":        total,
			"total_pages":  (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// Get 获取单个媒体详情
func (c *MediaController) Get(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var media models.Media

	if err := db.First(&media, id).Error; err != nil {
		hCtx.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Media not found",
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    media,
	})
}

// Update 更新媒体信息
func (c *MediaController) Update(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Alt         string `json:"alt"`
	}

	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request data",
		})
		return
	}

	db := database.GetDB()
	var media models.Media

	if err := db.First(&media, id).Error; err != nil {
		hCtx.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Media not found",
		})
		return
	}

	// 更新字段
	media.Title = req.Title
	media.Description = req.Description
	media.Alt = req.Alt

	if err := db.Save(&media).Error; err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to update media",
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    media,
		"message": "Media updated successfully",
	})
}

// Delete 删除媒体
func (c *MediaController) Delete(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")

	db := database.GetDB()
	var media models.Media

	if err := db.First(&media, id).Error; err != nil {
		hCtx.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Media not found",
		})
		return
	}

	// 删除物理文件
	if err := c.uploader.Delete(media.FilePath); err != nil {
		// 记录错误但继续删除数据库记录
		fmt.Printf("Failed to delete file: %v\n", err)
	}

	// 删除数据库记录
	if err := db.Delete(&media).Error; err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to delete media",
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Media deleted successfully",
	})
}
