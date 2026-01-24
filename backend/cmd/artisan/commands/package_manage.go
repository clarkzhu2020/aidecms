package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func PackageAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("Package name and version are required")
		return
	}

	pkg := args[0]
	version := ""
	if len(args) > 1 {
		version = args[1]
	}

	// 添加到go.mod
	cmdArgs := []string{"get", pkg}
	if version != "" {
		cmdArgs = append(cmdArgs, version)
	}

	cmd := exec.Command("go", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to add package: %v\n", err)
		fmt.Println(string(output))
		return
	}

	// 记录到packages.json
	recordPackage(pkg, version)
	fmt.Printf("Successfully added package: %s %s\n", pkg, version)
}

func recordPackage(name, version string) {
	filePath := filepath.Join("storage", "framework", "packages.json")

	var packages []Package
	if data, err := os.ReadFile(filePath); err == nil {
		json.Unmarshal(data, &packages)
	}

	packages = append(packages, Package{
		Name:    name,
		Version: version,
	})

	data, _ := json.MarshalIndent(packages, "", "  ")
	os.WriteFile(filePath, data, 0644)
}
