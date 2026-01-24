package commands

import (
	"fmt"
	"os"
	"os/exec"
)

func TestRun(args []string) {
	// 默认运行所有测试
	testArgs := []string{"test", "./..."}
	if len(args) > 0 {
		// 支持指定测试路径
		testArgs = append([]string{"test"}, args...)
	}

	cmd := exec.Command("go", testArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running tests...")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
