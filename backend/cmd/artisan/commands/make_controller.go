package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func MakeController(args []string) {
	if len(args) < 1 {
		fmt.Println("Controller name is required")
		return
	}

	name := args[0]
	dir := filepath.Join("app", "Http", "Controllers")
	filePath := filepath.Join(dir, name+".go")

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	tmpl := `package controllers

import (
	"github.com/cloudwego/hertz/pkg/app"
)

type {{.Name}}Controller struct {
	// Add dependencies here
}

func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{}
}

// Add your controller methods here
func (c *{{.Name}}Controller) Index(ctx *app.RequestContext) {
	ctx.String(200, "Hello from {{.Name}}Controller")
}
`

	data := struct{ Name string }{Name: name}
	t, err := template.New("controller").Parse(tmpl)
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

	fmt.Printf("Controller created: %s\n", filePath)
}
