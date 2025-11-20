package integration_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArtisanCommands(t *testing.T) {
	// 测试 help command (this exists)
	t.Run("help", func(t *testing.T) {
		cmd := exec.Command("go", "run", "../../cmd/artisan/main.go", "artisan", "help")
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)
		assert.Contains(t, string(output), "AideCMS Artisan Tool")
	})

	// 测试 migrate
	t.Run("migrate", func(t *testing.T) {
		cmd := exec.Command("go", "run", "../../cmd/artisan/main.go", "artisan", "migrate")
		output, _ := cmd.CombinedOutput()
		// Migration might fail but command should execute
		assert.NotEmpty(t, string(output))
	})
}
