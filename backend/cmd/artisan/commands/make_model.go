package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func MakeModel(args []string) {
	if len(args) < 1 {
		fmt.Println("Model name is required")
		return
	}

	name := args[0]
	dir := filepath.Join("app", "Models")
	filePath := filepath.Join(dir, name+".go")

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	tmpl := `package models

import (
	"gorm.io/gorm"
)

type {{.Name}} struct {
	gorm.Model
	// Add your fields here
}

// Add your model methods here
`

	data := struct{ Name string }{Name: name}
	t, err := template.New("model").Parse(tmpl)
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		fmt.Printf("Failed to execute template: %v\n", err)
		return
	}

	fmt.Printf("Model created: %s\n", filePath)
}
