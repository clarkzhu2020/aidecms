package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func MakeListener(args []string) {
	if len(args) < 1 {
		fmt.Println("Listener name is required")
		return
	}

	name := args[0]
	dir := filepath.Join("app", "Listeners")
	filePath := filepath.Join(dir, name+".go")

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	tmpl := `package listeners

import "github.com/chenyusolar/aidecms/app/Events"

type {{.Name}} struct {}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}

func (l *{{.Name}}) Handle(event *Events.ExampleEvent) error {
	// Add your listener logic here
	return nil
}`

	data := struct{ Name string }{Name: name}
	t, err := template.New("listener").Parse(tmpl)
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

	fmt.Printf("Listener created: %s\n", filePath)
}
