package upload

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ImageProcessor 图片处理器
type ImageProcessor struct {
	Storage Storage
}

// NewImageProcessor 创建图片处理器
func NewImageProcessor(storage Storage) *ImageProcessor {
	return &ImageProcessor{
		Storage: storage,
	}
}

// ThumbnailSize 缩略图尺寸配置
type ThumbnailSize struct {
	Name   string // 尺寸名称，如 "small", "medium", "large"
	Width  int    // 宽度
	Height int    // 高度
}

// ImageResult 图片处理结果
type ImageResult struct {
	Original   string            `json:"original"`   // 原图路径
	Thumbnails map[string]string `json:"thumbnails"` // 缩略图路径映射
}

// ProcessImage 处理图片（生成缩略图）
func (p *ImageProcessor) ProcessImage(sourcePath string, sizes []ThumbnailSize) (*ImageResult, error) {
	// 读取原图
	fullPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, sourcePath)
	src, err := imaging.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	result := &ImageResult{
		Original:   sourcePath,
		Thumbnails: make(map[string]string),
	}

	// 生成各个尺寸的缩略图
	for _, size := range sizes {
		thumbnailPath, err := p.createThumbnail(src, sourcePath, size)
		if err != nil {
			return nil, fmt.Errorf("failed to create thumbnail %s: %w", size.Name, err)
		}
		result.Thumbnails[size.Name] = thumbnailPath
	}

	return result, nil
}

// createThumbnail 创建单个缩略图
func (p *ImageProcessor) createThumbnail(src image.Image, originalPath string, size ThumbnailSize) (string, error) {
	// 调整图片大小（保持宽高比）
	thumb := imaging.Fit(src, size.Width, size.Height, imaging.Lanczos)

	// 生成缩略图路径
	ext := filepath.Ext(originalPath)
	nameWithoutExt := strings.TrimSuffix(originalPath, ext)
	thumbPath := fmt.Sprintf("%s_%s%s", nameWithoutExt, size.Name, ext)

	// 保存缩略图
	fullPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, thumbPath)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 保存图片
	if err := imaging.Save(thumb, fullPath); err != nil {
		return "", fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return thumbPath, nil
}

// Resize 调整图片大小
func (p *ImageProcessor) Resize(sourcePath string, width, height int) (string, error) {
	fullPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, sourcePath)
	src, err := imaging.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	// 调整大小
	resized := imaging.Resize(src, width, height, imaging.Lanczos)

	// 生成新路径
	ext := filepath.Ext(sourcePath)
	nameWithoutExt := strings.TrimSuffix(sourcePath, ext)
	newPath := fmt.Sprintf("%s_%dx%d%s", nameWithoutExt, width, height, ext)

	fullNewPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, newPath)

	// 保存图片
	if err := imaging.Save(resized, fullNewPath); err != nil {
		return "", fmt.Errorf("failed to save resized image: %w", err)
	}

	return newPath, nil
}

// Crop 裁剪图片
func (p *ImageProcessor) Crop(sourcePath string, x, y, width, height int) (string, error) {
	fullPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, sourcePath)
	src, err := imaging.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	// 裁剪
	cropped := imaging.Crop(src, image.Rect(x, y, x+width, y+height))

	// 生成新路径
	ext := filepath.Ext(sourcePath)
	nameWithoutExt := strings.TrimSuffix(sourcePath, ext)
	newPath := fmt.Sprintf("%s_cropped%s", nameWithoutExt, ext)

	fullNewPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, newPath)

	// 保存图片
	if err := imaging.Save(cropped, fullNewPath); err != nil {
		return "", fmt.Errorf("failed to save cropped image: %w", err)
	}

	return newPath, nil
}

// Compress 压缩图片
func (p *ImageProcessor) Compress(sourcePath string, quality int) error {
	fullPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, sourcePath)

	// 打开原图
	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := strings.ToLower(filepath.Ext(sourcePath))

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		img, _, err = image.Decode(file)
	}

	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// 创建临时文件
	tmpPath := fullPath + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// 重新编码并压缩
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: quality})
	case ".png":
		err = png.Encode(tmpFile, img)
	default:
		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: quality})
	}

	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to encode image: %w", err)
	}

	// 替换原文件
	if err := os.Rename(tmpPath, fullPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to replace original file: %w", err)
	}

	return nil
}

// GetImageInfo 获取图片信息
func (p *ImageProcessor) GetImageInfo(path string) (map[string]interface{}, error) {
	fullPath := filepath.Join(p.Storage.(*LocalStorage).BasePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image config: %w", err)
	}

	return map[string]interface{}{
		"width":  config.Width,
		"height": config.Height,
		"format": format,
	}, nil
}

// IsImage 检查文件是否为图片
func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}
