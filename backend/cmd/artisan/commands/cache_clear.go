package commands

import (
	"fmt"
	"os"
)

func CacheClear(args []string) {
	cacheDir := "storage/framework/cache"

	// 删除缓存目录
	if err := os.RemoveAll(cacheDir); err != nil {
		fmt.Printf("Failed to clear cache: %v\n", err)
		return
	}

	// 重新创建缓存目录
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		fmt.Printf("Failed to recreate cache directory: %v\n", err)
		return
	}

	fmt.Println("Application cache cleared successfully")
}
