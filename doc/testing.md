# 测试指南

AideCMS 提供了完整的测试支持，包括单元测试、集成测试和 API 测试。

## 单元测试

### 控制器测试

`internal/app/http/controllers/user_controller_test.go`:

```go
package controllers_test

import (
    "testing"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/stretchr/testify/assert"
    "internal/app/http/controllers"
    "internal/app/models"
)

func TestUserController_Index(t *testing.T) {
    // 准备测试数据
    models.DB.Create(&models.User{Name: "Test User"})
    
    // 创建请求上下文
    hCtx := app.NewContext(0)
    
    // 调用控制器方法
    c := controllers.UserController{}
    c.Index(context.Background(), hCtx)
    
    // 验证响应
    assert.Equal(t, 200, hCtx.Response.StatusCode())
    
    var users []models.User
    err := json.Unmarshal(hCtx.Response.Body(), &users)
    assert.NoError(t, err)
    assert.Greater(t, len(users), 0)
}
```

### 服务测试

`internal/app/services/mail_service_test.go`:

```go
package services_test

import (
    "testing"
    "internal/app/services"
)

func TestSendEmail(t *testing.T) {
    // 使用测试驱动
    services.MailDriver = "log"
    
    err := services.Mail.To("test@example.com").
        Subject("Test").
        Body("Test email").
        Send()
        
    assert.NoError(t, err)
}
```

## 集成测试

### 数据库测试

`internal/app/models/user_test.go`:

```go
package models_test

import (
    "testing"
    "internal/app/models"
)

func TestUserCreation(t *testing.T) {
    user := models.User{
        Name:  "Test",
        Email: "test@example.com",
    }
    
    err := models.DB.Create(&user).Error
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

## API 测试

### 使用 httpexpect

`tests/api/users_test.go`:

```go
package api_test

import (
    "testing"
    "github.com/cloudwego/hertz/pkg/common/test/assert"
    "github.com/cloudwego/hertz/pkg/common/ut"
    "internal/app/http"
)

func TestUserAPI(t *testing.T) {
    h := http.NewServer()
    defer h.Close()
    
    // 获取用户列表
    resp := ut.PerformRequest(h, "GET", "/api/users", nil)
    assert.DeepEqual(t, 200, resp.StatusCode())
    
    // 创建用户
    user := map[string]interface{}{
        "name":     "Test",
        "email":    "test@example.com",
        "password": "password",
    }
    resp = ut.PerformRequest(h, "POST", "/api/users", user)
    assert.DeepEqual(t, 201, resp.StatusCode())
}
```

## 测试覆盖率

生成测试覆盖率报告:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 持续集成

`.github/workflows/test.yml` 示例:

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.18'
      - run: go test -v -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v1