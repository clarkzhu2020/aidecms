package integration_test

import (
	"context"
	"testing"

	"github.com/clarkzhu2020/aidecms/pkg/framework"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	app := framework.NewApplication().
		SetConfigPath("../../config").
		Boot()

	assert.NotNil(t, app.DB)

	// 测试数据库连接
	if app.DB != nil {
		sqlDB, err := app.DB.DB.DB()
		assert.NoError(t, err)
		assert.NoError(t, sqlDB.Ping())
	}
}

func TestRedisConnection(t *testing.T) {
	app := framework.NewApplication().
		SetConfigPath("../../config").
		Boot()

	// 只有当Redis启用时才测试
	if app.Redis != nil {
		_, err := app.Redis.Ping(context.Background()).Result()
		assert.NoError(t, err)
	}
}
