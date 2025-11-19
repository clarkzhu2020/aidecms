package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Storage 存储接口
type Storage interface {
	// Save 保存文件
	Save(file io.Reader, path string) error
	// Delete 删除文件
	Delete(path string) error
	// Exists 检查文件是否存在
	Exists(path string) bool
	// URL 获取文件访问URL
	URL(path string) string
	// Size 获取文件大小
	Size(path string) (int64, error)
}

// LocalStorage 本地存储
type LocalStorage struct {
	BasePath string // 存储根目录
	BaseURL  string // 访问URL前缀
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

// Save 保存文件到本地
func (s *LocalStorage) Save(file io.Reader, path string) error {
	fullPath := filepath.Join(s.BasePath, path)

	// 创建目录
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 复制数据
	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.BasePath, path)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(path string) bool {
	fullPath := filepath.Join(s.BasePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// URL 获取文件访问URL
func (s *LocalStorage) URL(path string) string {
	return s.BaseURL + "/" + path
}

// Size 获取文件大小
func (s *LocalStorage) Size(path string) (int64, error) {
	fullPath := filepath.Join(s.BasePath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
