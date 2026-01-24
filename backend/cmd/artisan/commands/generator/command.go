package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Command struct {
	RootDir     string
	TemplateDir string
}

func NewCommand() *Command {
	return &Command{
		RootDir:     "app",
		TemplateDir: "stubs",
	}
}

func (c *Command) Handle(args []string) {
	if len(args) < 1 {
		fmt.Println("Not enough arguments (missing: name)")
		return
	}

	name := args[0]
	parts := strings.Split(name, "/")
	structName := parts[len(parts)-1]
	structName = strings.Title(strings.Replace(structName, "_", "", -1))

	data := map[string]interface{}{
		"Name":      structName,
		"Namespace": strings.Join(parts[:len(parts)-1], "\\"),
	}

	var templateFile string
	switch os.Args[2] {
	case "make:command":
		templateFile = "command.stub"
	case "make:controller":
		templateFile = "controller.stub"
	case "make:model":
		templateFile = "model.stub"
	default:
		fmt.Println("Unsupported make command")
		return
	}

	if err := c.generateFile(templateFile, name, data); err != nil {
		fmt.Printf("Error generating file: %v\n", err)
	}
}

func (c *Command) generateFile(templateFile, name string, data map[string]interface{}) error {
	// Read template
	tplPath := filepath.Join(c.TemplateDir, templateFile)
	tplContent, err := os.ReadFile(tplPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Parse template
	tpl, err := template.New("").Parse(string(tplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create directory if needed
	outputPath := filepath.Join(c.RootDir, name+".go")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("Created %s\n", outputPath)
	return nil
}
