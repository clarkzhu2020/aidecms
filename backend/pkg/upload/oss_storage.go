package upload

import (
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSStorage 阿里云OSS存储
type OSSStorage struct {
	client     *oss.Client
	bucketName string
	baseURL    string // CDN域名或Bucket域名
}

// OSSConfig 阿里云OSS配置
type OSSConfig struct {
	Endpoint        string // OSS Endpoint (如: oss-cn-hangzhou.aliyuncs.com)
	AccessKeyID     string // AccessKey ID
	AccessKeySecret string // AccessKey Secret
	BucketName      string // Bucket名称
	BaseURL         string // 访问域名 (CDN域名或Bucket域名)
}

// NewOSSStorage 创建阿里云OSS存储实例
func NewOSSStorage(config *OSSConfig) (*OSSStorage, error) {
	// 创建OSS客户端
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %w", err)
	}

	return &OSSStorage{
		client:     client,
		bucketName: config.BucketName,
		baseURL:    config.BaseURL,
	}, nil
}

// Save 上传文件到OSS
func (s *OSSStorage) Save(file io.Reader, path string) error {
	// 获取Bucket
	bucket, err := s.client.Bucket(s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %w", err)
	}

	// 上传文件
	if err := bucket.PutObject(path, file); err != nil {
		return fmt.Errorf("failed to upload file to OSS: %w", err)
	}

	return nil
}

// Delete 删除OSS文件
func (s *OSSStorage) Delete(path string) error {
	bucket, err := s.client.Bucket(s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %w", err)
	}

	if err := bucket.DeleteObject(path); err != nil {
		return fmt.Errorf("failed to delete file from OSS: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (s *OSSStorage) Exists(path string) bool {
	bucket, err := s.client.Bucket(s.bucketName)
	if err != nil {
		return false
	}

	exists, err := bucket.IsObjectExist(path)
	if err != nil {
		return false
	}

	return exists
}

// URL 获取文件访问URL
func (s *OSSStorage) URL(path string) string {
	if s.baseURL != "" {
		return s.baseURL + "/" + path
	}
	// 如果没有配置BaseURL，使用Bucket默认域名
	return fmt.Sprintf("https://%s.%s/%s", s.bucketName, s.client.Config.Endpoint, path)
}

// Size 获取文件大小
func (s *OSSStorage) Size(path string) (int64, error) {
	bucket, err := s.client.Bucket(s.bucketName)
	if err != nil {
		return 0, fmt.Errorf("failed to get bucket: %w", err)
	}

	// 获取对象元信息
	props, err := bucket.GetObjectDetailedMeta(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get object meta: %w", err)
	}

	// 获取Content-Length
	contentLength := props.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("content-length not found")
	}

	var size int64
	fmt.Sscanf(contentLength, "%d", &size)
	return size, nil
}

// SignURL 生成签名URL（用于私有文件访问）
func (s *OSSStorage) SignURL(path string, expireSeconds int64) (string, error) {
	bucket, err := s.client.Bucket(s.bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket: %w", err)
	}

	signedURL, err := bucket.SignURL(path, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return signedURL, nil
}
