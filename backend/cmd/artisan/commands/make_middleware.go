package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func MakeMiddleware(args []string) {
	if len(args) < 1 {
		fmt.Println("Middleware name is required")
		return
	}

	name := args[0]
	dir := filepath.Join("app", "Http", "Middleware")
	filePath := filepath.Join(dir, name+".go")

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	tmpl := `package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
)

func {{.Name}}() app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		// Add your middleware logic here
		ctx.Next()
	}
}
`

	data := struct{ Name string }{Name: name}
	t, err := template.New("middleware").Parse(tmpl)
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

	fmt.Printf("Middleware created: %s\n", filePath)
}
