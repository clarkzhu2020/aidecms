# AideCMS 开发框架

AideCMS 是一个基于 CloudWeGo Hertz 的高性能 Go Web 开发框架，提供了优雅的语法和丰富的功能。

## 目录结构

```
clarkgo/
├── cmd/               # 命令行入口
│   └── artisan/       # Artisan 命令行工具
├── config/            # 配置文件
├── doc/               # 文档
├── internal/          # 内部应用代码
│   └── app/           # 应用核心
├── pkg/               # 可复用的包
├── storage/           # 存储目录
│   ├── cache/         # 缓存文件
│   ├── logs/          # 日志文件
│   └── queue/         # 队列数据
└── go.mod             # Go 模块定义
```

## 安装

1. 安装 Go (1.18+)
2. 克隆项目:
   ```bash
   git clone https://github.com/your-repo/clarkgo.git
   ```
3. 安装依赖:
   ```bash
   cd clarkgo
   go mod tidy
   ```

## 配置

复制 `.env.example` 为 `.env` 并修改配置:

```bash
cp .env.example .env
```

主要配置项:

```
APP_NAME=AideCMS
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost:8888

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=aidecms
DB_USERNAME=root
DB_PASSWORD=

QUEUE_CONNECTION=redis
```

## Artisan 命令行工具

框架提供了强大的命令行工具:

```bash
go run cmd/artisan/main.go [command]
```

常用命令:

| 命令 | 描述 |
|------|------|
| `serve` | 启动开发服务器 |
| `make:controller` | 创建新控制器 |
| `make:model` | 创建新模型 |
| `queue:work` | 处理队列任务 |
| `queue:pause` | 暂停队列处理 |
| `queue:resume` | 恢复队列处理 |
| `queue:status` | 查看队列状态 |

## 路由

路由定义在 `internal/app/http/routes.go`:

```go
package http

import "github.com/cloudwego/hertz/pkg/app/server"

func RegisterRoutes(h *server.Hertz) {
    h.GET("/", indexHandler)
    h.POST("/users", userHandler.Create)
}
```

## 控制器

控制器位于 `internal/app/http/controllers`:

```go
package controllers

type UserController struct{}

func (c *UserController) Create(ctx context.Context, hCtx *app.RequestContext) {
    // 控制器逻辑
}
```

## 模型

模型位于 `internal/app/models`:

```go
package models

type User struct {
    gorm.Model
    Name  string
    Email string `gorm:"unique"`
}