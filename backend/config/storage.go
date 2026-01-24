package config

import (
	"fmt"
	"os"

	"github.com/clarkzhu2020/aidecms/pkg/upload"
)

// StorageDriver 存储驱动类型
type StorageDriver string

const (
	DriverLocal StorageDriver = "local" // 本地文件存储
	DriverOSS   StorageDriver = "oss"   // 阿里云OSS
	DriverS3    StorageDriver = "s3"    // AWS S3
)

// GetStorageDriver 获取存储驱动配置
func GetStorageDriver() StorageDriver {
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "" {
		return DriverLocal // 默认使用本地存储
	}
	return StorageDriver(driver)
}

// GetStorage 获取存储实例
func GetStorage() (upload.Storage, error) {
	driver := GetStorageDriver()

	switch driver {
	case DriverLocal:
		return getLocalStorage()
	case DriverOSS:
		return getOSSStorage()
	case DriverS3:
		return getS3Storage()
	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", driver)
	}
}

// getLocalStorage 获取本地存储实例
func getLocalStorage() (upload.Storage, error) {
	basePath := os.Getenv("LOCAL_STORAGE_PATH")
	if basePath == "" {
		basePath = "./storage/uploads" // 默认路径
	}

	baseURL := os.Getenv("LOCAL_STORAGE_URL")
	if baseURL == "" {
		baseURL = "/uploads" // 默认URL前缀
	}

	return upload.NewLocalStorage(basePath, baseURL), nil
}

// getOSSStorage 获取阿里云OSS存储实例
func getOSSStorage() (upload.Storage, error) {
	config := &upload.OSSConfig{
		Endpoint:        os.Getenv("OSS_ENDPOINT"),
		AccessKeyID:     os.Getenv("OSS_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
		BucketName:      os.Getenv("OSS_BUCKET_NAME"),
		BaseURL:         os.Getenv("OSS_BASE_URL"),
	}

	// 验证必需配置
	if config.Endpoint == "" || config.AccessKeyID == "" ||
		config.AccessKeySecret == "" || config.BucketName == "" {
		return nil, fmt.Errorf("OSS configuration is incomplete")
	}

	return upload.NewOSSStorage(config)
}

// getS3Storage 获取AWS S3存储实例
func getS3Storage() (upload.Storage, error) {
	config := &upload.S3Config{
		Region:          os.Getenv("S3_REGION"),
		AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("S3_BUCKET_NAME"),
		BaseURL:         os.Getenv("S3_BASE_URL"),
	}

	// 验证必需配置
	if config.Region == "" || config.AccessKeyID == "" ||
		config.SecretAccessKey == "" || config.BucketName == "" {
		return nil, fmt.Errorf("S3 configuration is incomplete")
	}

	return upload.NewS3Storage(config)
}
