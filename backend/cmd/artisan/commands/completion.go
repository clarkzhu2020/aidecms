package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateCompletion(args []string) {
	if len(args) < 1 {
		fmt.Println("Shell type is required (bash or zsh)")
		return
	}

	shell := args[0]
	switch shell {
	case "bash":
		generateBashCompletion()
	case "zsh":
		generateZshCompletion()
	default:
		fmt.Println("Unsupported shell type. Use bash or zsh")
	}
}

func generateBashCompletion() {
	completion := `_artisan_completion() {
    local cur prev words cword
    _init_completion || return

    case ${prev} in
        artisan)
            COMPREPLY=($(compgen -W "make:controller make:model migrate cache:clear" -- "${cur}"))
            ;;
        *)
            COMPREPLY=()
            ;;
    esac
}

complete -F _artisan_completion artisan`

	saveCompletion("bash", completion)
}

func generateZshCompletion() {
	completion := `#compdef artisan

_artisan() {
    local state
    _arguments \
        '1: :->command' \
        '*: :->args'
    
    case $state in
        command)
            _values "artisan command" \
                "make:controller[Create controller]" \
                "make:model[Create model]" \
                "migrate[Run migrations]" \
                "cache:clear[Clear cache]"
            ;;
    esac
}

compdef _artisan artisan`

	saveCompletion("zsh", completion)
}

func saveCompletion(shell, content string) {
	dir := filepath.Join("storage", "framework", "completion")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create completion directory: %v\n", err)
		return
	}

	filePath := filepath.Join(dir, "artisan."+shell)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		fmt.Printf("Failed to save completion file: %v\n", err)
		return
	}

	fmt.Printf("%s completion file generated at: %s\n", shell, filePath)
	fmt.Printf("Add this to your ~/.%src file:\nsource %s\n", shell, filePath)
}
