package commands

import (
	"fmt"
	"os/exec"
)

func PackageInstall(args []string) {
	if len(args) < 1 {
		fmt.Println("Package name is required")
		return
	}

	pkg := args[0]
	fmt.Printf("Installing package: %s\n", pkg)

	// 使用go get安装包
	cmd := exec.Command("go", "get", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to install package: %v\n", err)
		fmt.Println(string(output))
		return
	}

	fmt.Printf("Successfully installed package: %s\n", pkg)
}
