package upload

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Storage AWS S3存储
type S3Storage struct {
	client     *s3.S3
	bucketName string
	region     string
	baseURL    string // CloudFront CDN域名或S3域名
}

// S3Config AWS S3配置
type S3Config struct {
	Region          string // AWS Region (如: us-east-1, ap-southeast-1)
	AccessKeyID     string // AWS AccessKey ID
	SecretAccessKey string // AWS SecretAccessKey
	BucketName      string // Bucket名称
	BaseURL         string // 访问域名 (CloudFront CDN或S3域名)
}

// NewS3Storage 创建AWS S3存储实例
func NewS3Storage(config *S3Config) (*S3Storage, error) {
	// 创建AWS会话
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// 创建S3客户端
	client := s3.New(sess)

	return &S3Storage{
		client:     client,
		bucketName: config.BucketName,
		region:     config.Region,
		baseURL:    config.BaseURL,
	}, nil
}

// Save 上传文件到S3
func (s *S3Storage) Save(file io.Reader, path string) error {
	// 读取文件内容
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 上传到S3
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"), // 公开读
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}

// Delete 删除S3文件
func (s *S3Storage) Delete(path string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (s *S3Storage) Exists(path string) bool {
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	return err == nil
}

// URL 获取文件访问URL
func (s *S3Storage) URL(path string) string {
	if s.baseURL != "" {
		return s.baseURL + "/" + path
	}
	// 如果没有配置BaseURL，使用S3默认域名
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, path)
}

// Size 获取文件大小
func (s *S3Storage) Size(path string) (int64, error) {
	result, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get object meta: %w", err)
	}

	if result.ContentLength == nil {
		return 0, fmt.Errorf("content-length not found")
	}

	return *result.ContentLength, nil
}

// SignURL 生成签名URL（用于私有文件访问）
func (s *S3Storage) SignURL(path string, expireSeconds int64) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})

	signedURL, err := req.Presign(time.Duration(expireSeconds) * time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return signedURL, nil
}
