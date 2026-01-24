package upload

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Uploader 文件上传器
type Uploader struct {
	Storage      Storage
	MaxSize      int64    // 最大文件大小（字节）
	AllowedExts  []string // 允许的文件扩展名
	AllowedMimes []string // 允许的MIME类型
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize      int64
	AllowedExts  []string
	AllowedMimes []string
	Storage      Storage
}

// NewUploader 创建上传器
func NewUploader(config *UploadConfig) *Uploader {
	if config.MaxSize == 0 {
		config.MaxSize = 10 * 1024 * 1024 // 默认10MB
	}

	if len(config.AllowedExts) == 0 {
		config.AllowedExts = []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"}
	}

	return &Uploader{
		Storage:      config.Storage,
		MaxSize:      config.MaxSize,
		AllowedExts:  config.AllowedExts,
		AllowedMimes: config.AllowedMimes,
	}
}

// UploadResult 上传结果
type UploadResult struct {
	OriginalName string `json:"original_name"` // 原始文件名
	FileName     string `json:"file_name"`     // 存储文件名
	Path         string `json:"path"`          // 存储路径
	URL          string `json:"url"`           // 访问URL
	Size         int64  `json:"size"`          // 文件大小
	Extension    string `json:"extension"`     // 文件扩展名
	MimeType     string `json:"mime_type"`     // MIME类型
	Hash         string `json:"hash"`          // 文件哈希
}

// Upload 上传单个文件
func (u *Uploader) Upload(fileHeader *multipart.FileHeader) (*UploadResult, error) {
	// 验证文件
	if err := u.validate(fileHeader); err != nil {
		return nil, err
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 计算文件哈希
	hash, err := u.calculateHash(file)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	// 重置文件读取位置
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file: %w", err)
	}

	// 生成存储路径
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	storagePath := u.generatePath(ext)

	// 保存文件
	if err := u.Storage.Save(file, storagePath); err != nil {
		return nil, err
	}

	result := &UploadResult{
		OriginalName: fileHeader.Filename,
		FileName:     filepath.Base(storagePath),
		Path:         storagePath,
		URL:          u.Storage.URL(storagePath),
		Size:         fileHeader.Size,
		Extension:    ext,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Hash:         hash,
	}

	return result, nil
}

// UploadMultiple 上传多个文件
func (u *Uploader) UploadMultiple(files []*multipart.FileHeader) ([]*UploadResult, error) {
	results := make([]*UploadResult, 0, len(files))

	for _, fileHeader := range files {
		result, err := u.Upload(fileHeader)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

// validate 验证文件
func (u *Uploader) validate(fileHeader *multipart.FileHeader) error {
	// 检查文件大小
	if fileHeader.Size > u.MaxSize {
		return fmt.Errorf("file size exceeds maximum allowed size: %d bytes", u.MaxSize)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !u.isAllowedExtension(ext) {
		return fmt.Errorf("file extension %s is not allowed", ext)
	}

	// 检查MIME类型
	mimeType := fileHeader.Header.Get("Content-Type")
	if len(u.AllowedMimes) > 0 && !u.isAllowedMime(mimeType) {
		return fmt.Errorf("mime type %s is not allowed", mimeType)
	}

	return nil
}

// isAllowedExtension 检查扩展名是否允许
func (u *Uploader) isAllowedExtension(ext string) bool {
	for _, allowed := range u.AllowedExts {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}

// isAllowedMime 检查MIME类型是否允许
func (u *Uploader) isAllowedMime(mime string) bool {
	for _, allowed := range u.AllowedMimes {
		if strings.EqualFold(mime, allowed) {
			return true
		}
	}
	return false
}

// generatePath 生成存储路径
func (u *Uploader) generatePath(ext string) string {
	now := time.Now()
	// 按日期分目录: uploads/2024/01/02/uuid.ext
	return filepath.Join(
		"uploads",
		now.Format("2006"),
		now.Format("01"),
		now.Format("02"),
		uuid.New().String()+ext,
	)
}

// calculateHash 计算文件MD5哈希
func (u *Uploader) calculateHash(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Delete 删除文件
func (u *Uploader) Delete(path string) error {
	return u.Storage.Delete(path)
}
