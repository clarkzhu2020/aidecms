package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func PackageInstallAll(args []string) {
	filePath := filepath.Join("storage", "framework", "packages.json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("No packages to install - packages.json not found")
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read packages.json: %v\n", err)
		return
	}

	var packages []Package
	if err := json.Unmarshal(data, &packages); err != nil {
		fmt.Printf("Failed to parse packages.json: %v\n", err)
		return
	}

	for _, pkg := range packages {
		fmt.Printf("Installing package: %s %s\n", pkg.Name, pkg.Version)

		cmdArgs := []string{"get", pkg.Name}
		if pkg.Version != "" {
			cmdArgs = append(cmdArgs, pkg.Version)
		}

		cmd := exec.Command("go", cmdArgs...)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("Failed to install %s: %v\n", pkg.Name, err)
			fmt.Println(string(output))
		} else {
			fmt.Printf("Successfully installed %s\n", pkg.Name)
		}
	}
}
