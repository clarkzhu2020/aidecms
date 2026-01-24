package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func MakeTest(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: make:test <type> <name>")
		fmt.Println("Types: unit, integration")
		return
	}

	testType := args[0]
	name := args[1]

	var dir, tmpl string

	switch testType {
	case "unit":
		dir = filepath.Join("test", "unit")
		tmpl = `package unit_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test{{.Name}}(t *testing.T) {
	assert.True(t, true, "Sample test")
}`
	case "integration":
		dir = filepath.Join("test", "integration")
		tmpl = `package integration_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test{{.Name}}(t *testing.T) {
	assert.True(t, true, "Sample integration test")
}`
	default:
		fmt.Printf("Unknown test type: %s\n", testType)
		return
	}

	filePath := filepath.Join(dir, strings.ToLower(name)+"_test.go")

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	data := struct{ Name string }{Name: name}
	t, err := template.New("test").Parse(tmpl)
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

	fmt.Printf("Test created: %s\n", filePath)
}
