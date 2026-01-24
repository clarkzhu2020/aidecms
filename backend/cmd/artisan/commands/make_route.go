package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MakeRoute(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: make:route <name> <controller>")
		return
	}

	name := args[0]
	controller := args[1]
	routeFile := filepath.Join("routes", "api.go")

	// 检查路由文件是否存在
	if _, err := os.Stat(routeFile); os.IsNotExist(err) {
		fmt.Printf("Route file not found: %s\n", routeFile)
		return
	}

	// 读取路由文件内容
	content, err := os.ReadFile(routeFile)
	if err != nil {
		fmt.Printf("Failed to read route file: %v\n", err)
		return
	}

	// 添加路由注册代码
	newContent := strings.Replace(string(content),
		"// [Artisan Routes]",
		fmt.Sprintf("\t%sRouter := r.Group(\"/%s\")\n\t%sRouter.GET(\"\", %s.Index)\n\t// [Artisan Routes]",
			name, strings.ToLower(name), name, controller),
		-1)

	// 写入更新后的内容
	if err := os.WriteFile(routeFile, []byte(newContent), 0644); err != nil {
		fmt.Printf("Failed to write route file: %v\n", err)
		return
	}

	fmt.Printf("Route added for %s in %s\n", name, routeFile)
}
