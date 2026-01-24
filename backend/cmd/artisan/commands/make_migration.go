package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func MakeMigration(args []string) {
	if len(args) < 1 {
		fmt.Println("Migration name is required")
		return
	}

	name := args[0]
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.go", timestamp, strings.ToLower(name))
	filePath := filepath.Join("database", "migrations", fileName)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	tmpl := `package migrations

import (
	"gorm.io/gorm"
)

func init() {
	migrationList = append(migrationList, &{{.Name}}{})
}

type {{.Name}} struct {}

func (m *{{.Name}}) Up(db *gorm.DB) error {
	// Migration logic here
	return nil
}

func (m *{{.Name}}) Down(db *gorm.DB) error {
	// Rollback logic here
	return nil
}`

	data := struct{ Name string }{Name: name}
	t, err := template.New("migration").Parse(tmpl)
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

	fmt.Printf("Migration created: %s\n", filePath)
}
