package commands

import (
	"fmt"
	"os"
	"os/exec"
)

func TestCoverage(args []string) {
	profile := "coverage.out"
	if len(args) > 0 {
		profile = args[0]
	}

	// 生成覆盖率数据
	cmd := exec.Command("go", "test", "-coverprofile="+profile, "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running tests with coverage...")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}

	// 生成HTML报告
	htmlCmd := exec.Command("go", "tool", "cover", "-html="+profile)
	htmlCmd.Stdout = os.Stdout
	htmlCmd.Stderr = os.Stderr

	fmt.Println("Generating coverage report...")
	if err := htmlCmd.Run(); err != nil {
		fmt.Printf("Failed to generate coverage report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Coverage report generated successfully")
}
