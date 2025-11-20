package unit_test

import (
	"testing"

	"github.com/chenyusolar/aidecms/pkg/framework"
	"github.com/stretchr/testify/assert"
)

func TestApplicationCreation(t *testing.T) {
	app := framework.NewApplication()
	assert.NotNil(t, app)
	assert.Equal(t, "AideCMS", app.AppName)
	assert.Equal(t, "1.0.0", app.AppVersion)
	assert.Equal(t, "development", app.Env)
	assert.True(t, app.Debug)
}

func TestApplicationBoot(t *testing.T) {
	app := framework.NewApplication().
		SetConfigPath("../../config").
		Boot()

	assert.NotNil(t, app.Server)
	assert.NotNil(t, app.Router)
	assert.NotNil(t, app.Config)
	assert.NotNil(t, app.Logger)
}

func TestApplicationRun(t *testing.T) {
	app := framework.NewApplication().
		SetConfigPath("../../config").
		Boot()

	// 测试运行不会panic
	assert.NotPanics(t, func() {
		go app.Run()
	})
}
