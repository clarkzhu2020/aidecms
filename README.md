# ClarkGo - 企业级 Go CMS 平台框架

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Hertz](https://img.shields.io/badge/Hertz-CloudWeGo-blue)](https://github.com/cloudwego/hertz)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

ClarkGo 是一个基于 CloudWeGo Hertz 框架开发的企业级 CMS 平台框架，提供完整的任务调度、队列系统、事件驱动、限流保护和健康监控等核心功能。

## ✨ 核心特性

### 🚀 高性能基础设施
- **Hertz 框架** - CloudWeGo 高性能 HTTP 服务器
- **GORM ORM** - 支持 MySQL、PostgreSQL、SQLite
- **Redis 集成** - 缓存、队列、会话存储
- **并发安全** - 所有核心组件线程安全设计

### ⚡ 任务调度系统 (Task Scheduler)
- ✅ Cron 表达式解析（支持 `*`, `-`, `/`, `,` 语法）
- ✅ 15+ 便捷方法（EveryMinute、Daily、WeeklyOn 等）
- ✅ 并发任务执行（Worker Pool）
- ✅ 任务统计和错误日志
- ✅ 手动触发任务

### 📬 队列系统 (Queue System)
- ✅ 统一队列接口（Driver 抽象）
- ✅ 内存驱动（开发/测试）
- ✅ Redis 驱动（生产环境）
- ✅ 延迟任务（DelayUntil）
- ✅ 失败重试（指数退避）
- ✅ 死信队列（DLQ）
- ✅ 任务超时控制

### 🎯 事件系统 (Event System)
- ✅ 同步/异步监听器
- ✅ 优先级支持（1-10）
- ✅ Worker Pool 并发控制
- ✅ 8 个预定义事件类型
- ✅ 事件统计（执行次数、成功率、平均耗时）
- ✅ 全局事件调度器

### 🛡️ 限流系统 (Rate Limiting)
- ✅ 令牌桶算法（Token Bucket）
- ✅ 滑动窗口算法（Sliding Window）
- ✅ 固定窗口算法（Fixed Window）
- ✅ 并发安全 + 自动 GC
- ✅ Hertz 中间件集成
- ✅ 按 IP/用户/全局限流

### 💚 健康检查系统 (Health Check)
- ✅ 7 种内置检查器（数据库、Redis、内存、磁盘等）
- ✅ 三态健康状态（Healthy/Degraded/Unhealthy）
- ✅ 并发检查 + 结果缓存
- ✅ Kubernetes 集成（Liveness/Readiness）
- ✅ 完整的 HTTP 端点

### 🎨 CMS 核心功能
- ✅ 用户认证（JWT）
- ✅ RBAC 权限管理
- ✅ 文章/分类/标签管理
- ✅ 评论系统
- ✅ 媒体库管理
- ✅ 菜单管理
- ✅ SEO 优化（Sitemap、Robots）
- ✅ AI 集成（多模型支持）
- ✅ 邮件系统

### 🔧 开发工具
- ✅ Artisan CLI 命令行工具
- ✅ 代码生成器（Controller、Model、Command）
- ✅ 数据库迁移
- ✅ Swagger API 文档
- ✅ 命令统计和分析

## 📦 安装

### 环境要求
- Go 1.18 或更高版本
- MySQL 5.7+ / PostgreSQL 10+ / SQLite 3
- Redis 5.0+ (可选，用于队列和缓存)

### 快速安装

```bash
# 克隆项目
git clone https://github.com/chenyusolar/clarkgo.git
cd clarkgo

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库等信息

# 运行数据库迁移
go run cmd/artisan/main.go artisan migrate

# 启动服务
go run main.go
```

服务默认运行在 `http://localhost:8888`

## 🚀 快速开始

### 1. 基础 Web 应用

```go
package main

import (
    "github.com/clarkgo/clarkgo/pkg/framework"
)

func main() {
    app := framework.NewApplication()
    
    // 注册路由
    app.RegisterRoutes(func(router *framework.Router) {
        router.GET("/", func(ctx context.Context, c *framework.RequestContext) {
            c.JSON(200, map[string]interface{}{
                "message": "Welcome to ClarkGo!",
            })
        })
    })
    
    // 启动服务器
    app.Run(":8888")
}
```

### 2. 任务调度

```go
import "github.com/clarkgo/clarkgo/pkg/schedule"

// 每分钟执行
schedule.EveryMinute().Do(func() {
    log.Println("Task executed every minute")
})

// 每天凌晨 2:00 执行
schedule.Daily().At("02:00").Do(func() {
    // 执行数据库备份
    backupDatabase()
})

// 每周一上午 9:00 执行
schedule.WeeklyOn(time.Monday).At("09:00").Do(func() {
    // 发送周报
    sendWeeklyReport()
})

// 启动调度器
go run cmd/artisan/main.go artisan schedule:work
```

### 3. 队列系统

```go
import "github.com/clarkgo/clarkgo/pkg/queue"

// 定义任务
type EmailJob struct {
    To      string
    Subject string
    Body    string
}

func (j *EmailJob) Handle() error {
    return sendEmail(j.To, j.Subject, j.Body)
}

// 推送任务到队列
queue.Push("emails", &EmailJob{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Welcome to our service!",
})

// 延迟任务（1小时后执行）
queue.DelayUntil("notifications", job, time.Now().Add(1*time.Hour))

// 启动队列 Worker
go run cmd/artisan/main.go artisan queue:worker default 5
```

### 4. 事件系统

```go
import "github.com/clarkgo/clarkgo/pkg/event"

// 注册监听器
event.Listen("user.registered", func(e event.Event) error {
    user := e.(*event.UserRegistered)
    log.Printf("New user registered: %s", user.Email)
    
    // 发送欢迎邮件
    return sendWelcomeEmail(user.Email)
})

// 异步监听器
event.ListenAsync("order.created", func(e event.Event) error {
    order := e.(*event.OrderCreated)
    return processOrder(order)
}, event.WithPriority(5))

// 触发事件
event.Dispatch(event.NewUserRegistered(userID, email, name))
```

### 5. 限流保护

```go
import (
    "github.com/clarkgo/clarkgo/pkg/ratelimit"
    "github.com/clarkgo/clarkgo/pkg/framework"
)

// 创建限流器
limiter := ratelimit.NewTokenBucket(100, 200) // 100 req/s, burst 200

// 全局限流
h.Use(framework.RateLimit(limiter, "global"))

// 按 IP 限流
h.Use(framework.RateLimitByIP(
    ratelimit.NewSlidingWindow(100, 1*time.Minute),
))

// 按用户限流
h.Use(framework.RateLimitByUser(
    ratelimit.NewFixedWindow(500, 1*time.Hour),
))
```

### 6. 健康检查

```go
import "github.com/clarkgo/clarkgo/pkg/health"

// 创建健康检查器
hc := health.NewHealthChecker(5 * time.Second)

// 注册检查器
hc.Register(health.NewDatabaseChecker(db))
hc.Register(health.NewRedisChecker(redisClient))
hc.Register(health.NewMemoryChecker(70.0, 90.0))

// 注册端点
h.GET("/health", framework.HealthEndpoint(hc))
h.GET("/health/ready", framework.ReadinessEndpoint(hc))
h.GET("/health/live", framework.LivenessEndpoint())
```

## 📁 项目结构

```
clarkgo/
├── app/                        # 应用层
│   ├── Http/
│   │   ├── Controllers/       # HTTP 控制器
│   │   │   ├── AIController.go
│   │   │   ├── PostController.go
│   │   │   ├── UserController.go
│   │   │   └── ...
│   │   └── Middleware/        # HTTP 中间件
│   │       ├── JWTMiddleware.go
│   │       └── PermissionMiddleware.go
│   └── shared/                 # 共享代码
│
├── cmd/                        # 命令行工具
│   ├── artisan/               # Artisan CLI
│   │   ├── main.go
│   │   └── commands/          # CLI 命令
│   │       ├── schedule_work.go
│   │       ├── queue_worker.go
│   │       ├── event_command.go
│   │       ├── ratelimit_command.go
│   │       └── health_command.go
│   ├── clarkgo/               # Web 服务器
│   └── migrate/               # 数据库迁移
│
├── config/                     # 配置文件
│   ├── database.go
│   ├── jwt.go
│   ├── mail.go
│   └── ...
│
├── database/                   # 数据库
│   └── migrations/            # 迁移文件
│
├── doc/                        # 文档
│   ├── PHASE5_SUMMARY.md      # Phase 5 实现总结
│   ├── QUICKREF.md            # 快速参考
│   ├── getting-started.md
│   ├── database.md
│   └── ...
│
├── internal/                   # 内部代码
│   └── app/
│       ├── models/            # 数据模型
│       ├── services/          # 业务服务
│       └── adapters/          # 适配器
│
├── pkg/                        # 核心框架包
│   ├── ai/                    # AI 集成
│   ├── cache/                 # 缓存系统
│   ├── database/              # 数据库连接
│   ├── event/                 # 事件系统 ⭐
│   ├── framework/             # 框架核心
│   │   ├── app.go
│   │   ├── route.go
│   │   ├── middleware.go
│   │   ├── ratelimit_middleware.go  ⭐
│   │   └── health_middleware.go     ⭐
│   ├── health/                # 健康检查 ⭐
│   │   ├── health.go
│   │   ├── checkers.go
│   │   └── health_test.go
│   ├── http/                  # HTTP 客户端
│   ├── log/                   # 日志系统
│   ├── mail/                  # 邮件服务
│   ├── queue/                 # 队列系统 ⭐
│   │   ├── queue.go
│   │   ├── memory_driver.go
│   │   ├── redis_driver.go
│   │   └── queue_test.go
│   ├── ratelimit/             # 限流系统 ⭐
│   │   ├── ratelimit.go
│   │   └── ratelimit_test.go
│   ├── redis/                 # Redis 连接
│   ├── response/              # 响应封装
│   ├── schedule/              # 任务调度 ⭐
│   │   ├── cron.go
│   │   ├── schedule.go
│   │   ├── builder.go
│   │   └── schedule_test.go
│   ├── seo/                   # SEO 工具
│   ├── swagger/               # Swagger 文档
│   ├── upload/                # 文件上传
│   └── validator/             # 验证器
│
├── routes/                     # 路由定义
│   ├── api.go
│   └── test.go
│
├── storage/                    # 存储目录
│   ├── database/              # SQLite 数据库
│   ├── logs/                  # 日志文件
│   └── stats/                 # 统计数据
│
├── test/                       # 测试
│   ├── integration/           # 集成测试
│   └── unit/                  # 单元测试
│
├── .env.example               # 环境变量示例
├── go.mod
├── go.sum
├── main.go                    # 入口文件
└── README.md

⭐ = Phase 5 新增核心功能
```

## 🛠️ Artisan CLI 命令

ClarkGo 提供了强大的 Artisan 命令行工具，用于开发和管理。

### 任务调度命令
```bash
# 启动调度器守护进程
artisan schedule:work

# 运行一次所有到期任务
artisan schedule:run

# 列出所有已注册任务
artisan schedule:list
```

### 队列系统命令
```bash
# 启动队列 Worker（5个并发）
artisan queue:worker default 5

# 查看队列状态
artisan queue:status

# 重试失败的任务
artisan queue:retry

# 清理旧任务
artisan queue:clean

# 查看队列统计
artisan queue:stats
```

### 事件系统命令
```bash
# 运行事件系统演示
artisan event:test

# 列出注册的事件和监听器
artisan event:list

# 显示事件统计信息
artisan event:stats
```

### 限流系统命令
```bash
# 运行限流演示
artisan ratelimit demo
```

### 健康检查命令
```bash
# 运行健康检查演示
artisan health demo
```

### 代码生成命令
```bash
# 生成控制器
artisan make:controller UserController

# 生成模型
artisan make:model User

# 生成中间件
artisan make:middleware AuthMiddleware

# 生成命令
artisan make:command SendEmails
```

### 数据库命令
```bash
# 运行数据库迁移
artisan migrate

# CMS 初始化
artisan cms:init
```

### AI 命令
```bash
# 配置 AI
artisan ai:setup openai sk-xxx

# 聊天
artisan ai:chat "Hello, how are you?"

# 文本补全
artisan ai:completion "Once upon a time"

# 列出可用模型
artisan ai:models

# 测试 AI 连接
artisan ai:test
```

### 统计命令
```bash
# 显示命令使用统计
artisan stats:show

# 重置统计数据
artisan stats:reset

# 导出统计数据
artisan stats:export json

# 生成使用图表
artisan stats:chart

# 清理旧统计数据
artisan stats:cleanup 30d

# 检查性能异常
artisan stats:check 5s
```

### 邮件命令
```bash
# 配置邮件告警
artisan alert:setup config.json

# 发送测试邮件
artisan alert:test
```

## ⚙️ 配置管理

### 环境变量配置

项目使用 `.env` 文件管理环境变量：

```bash
# 复制示例配置
cp .env.example .env
```

`.env` 配置示例：

```ini
# 应用配置
APP_ENV=production
APP_DEBUG=false
APP_PORT=8888

# 数据库配置
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=clarkgo
DB_USERNAME=root
DB_PASSWORD=secret

# Redis 配置
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 配置
JWT_SECRET=your-secret-key
JWT_TTL=3600

# 邮件配置
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=user@example.com
MAIL_PASSWORD=secret
MAIL_FROM_ADDRESS=noreply@example.com
MAIL_FROM_NAME=ClarkGo

# AI 配置（可选）
AI_PROVIDER=openai
AI_API_KEY=sk-xxx
AI_MODEL=gpt-3.5-turbo
```

### 配置读取

```go
import "github.com/clarkgo/clarkgo/pkg/config"

// 获取字符串配置
host := config.GetEnv("DB_HOST", "localhost")

// 获取整数配置
port := config.GetEnvInt("DB_PORT", 3306)

// 获取布尔配置
debug := config.GetEnvBool("APP_DEBUG", false)
```

## 📚 核心功能详解

### 1. 任务调度系统

支持 Cron 表达式和便捷方法：

```go
// Cron 表达式（分 时 日 月 周）
schedule.Cron("*/5 * * * *").Do(task)  // 每5分钟
schedule.Cron("0 2 * * *").Do(task)    // 每天凌晨2点
schedule.Cron("0 0 * * 1").Do(task)    // 每周一零点

// 便捷方法
schedule.EveryMinute().Do(task)                    // 每分钟
schedule.EveryFiveMinutes().Do(task)               // 每5分钟
schedule.Hourly().Do(task)                         // 每小时
schedule.HourlyAt(15).Do(task)                     // 每小时第15分钟
schedule.Daily().At("14:30").Do(task)              // 每天14:30
schedule.DailyAt("08:00").Do(task)                 // 每天08:00
schedule.WeeklyOn(time.Monday).At("09:00").Do(task)  // 每周一09:00
schedule.MonthlyOn(15).At("00:00").Do(task)        // 每月15号零点

// 自定义任务名称
schedule.Daily().Name("backup").Do(backupTask)

// 手动触发
schedule.RunTask("backup")
```

### 2. 队列系统

支持内存和 Redis 两种驱动：

```go
// 配置 Redis 驱动（生产环境推荐）
import "github.com/clarkgo/clarkgo/pkg/queue"

redisQueue := queue.NewRedisDriver(redisClient, "myapp")
queue.UseDriver("default", redisQueue)

// 定义任务
type ProcessVideoJob struct {
    VideoID string
    Options map[string]interface{}
}

func (j *ProcessVideoJob) Handle() error {
    // 处理视频任务
    return processVideo(j.VideoID, j.Options)
}

// 推送任务
queue.Push("videos", &ProcessVideoJob{
    VideoID: "vid-123",
    Options: map[string]interface{}{
        "quality": "1080p",
    },
})

// 延迟任务（2小时后执行）
queue.DelayUntil("videos", job, time.Now().Add(2*time.Hour))

// 失败重试配置
// 自动重试3次，使用指数退避算法
```

### 3. 事件系统

支持同步/异步事件处理：

```go
// 预定义事件
event.Dispatch(event.NewUserRegistered(userID, email, name))
event.Dispatch(event.NewPostPublished(postID, title, author))
event.Dispatch(event.NewOrderCreated(orderID, userID, amount))
event.Dispatch(event.NewOrderCompleted(orderID, status))
event.Dispatch(event.NewEmailSent(to, subject))
event.Dispatch(event.NewPaymentReceived(orderID, amount, method))

// 监听器优先级（1-10，数字越大优先级越高）
event.Listen("user.registered", handler1, event.WithPriority(10)) // 先执行
event.Listen("user.registered", handler2, event.WithPriority(5))  // 后执行

// 自定义事件
type ProductCreated struct {
    event.BaseEvent
    ProductID string
    Name      string
    Price     float64
}

func (e *ProductCreated) EventName() string {
    return "product.created"
}

// 触发自定义事件
event.Dispatch(&ProductCreated{
    ProductID: "prod-123",
    Name:      "New Product",
    Price:     99.99,
})

// 获取事件统计
stats := event.GetDispatcher().GetStats()
fmt.Printf("Total events: %d\n", stats["total_dispatches"])
```

### 4. 限流系统

三种算法适应不同场景：

```go
// 1. 令牌桶（Token Bucket）- 平滑限流 + 支持突发
// 适用场景：API 限流、一般场景
limiter := ratelimit.NewTokenBucket(100, 200) // 100 req/s, burst 200

// 2. 滑动窗口（Sliding Window）- 精确时间限制
// 适用场景：严格限制、Premium API
limiter := ratelimit.NewSlidingWindow(1000, 1*time.Minute) // 1000/分钟

// 3. 固定窗口（Fixed Window）- 简单高效
// 适用场景：基础限流、公共 API
limiter := ratelimit.NewFixedWindow(5000, 1*time.Hour) // 5000/小时

// 使用限流器
if limiter.Allow(userID) {
    // 处理请求
    handleRequest()
} else {
    // 返回 429 Too Many Requests
    return errors.New("rate limit exceeded")
}

// 批量检查
if limiter.AllowN(userID, 5) {
    // 允许5个请求
}

// 重置限制
limiter.Reset(userID)

// 获取统计（仅 Sliding Window）
stats := limiter.GetStats(userID)
fmt.Printf("Used: %d/%d\n", stats["requests"], stats["limit"])
```

#### Hertz 中间件集成

```go
import (
    "github.com/clarkgo/clarkgo/pkg/framework"
    "github.com/clarkgo/clarkgo/pkg/ratelimit"
)

// 全局限流
limiter := ratelimit.NewTokenBucket(1000, 2000)
h.Use(framework.RateLimit(limiter, "global"))

// 按 IP 限流
h.Use(framework.RateLimitByIP(
    ratelimit.NewSlidingWindow(100, 1*time.Minute),
))

// 按用户限流
h.Use(framework.RateLimitByUser(
    ratelimit.NewFixedWindow(500, 1*time.Hour),
))

// 带统计的限流
h.Use(framework.RateLimitWithStats(limiter, "api"))
```

### 5. 健康检查系统

全面的服务健康监控：

```go
import "github.com/clarkgo/clarkgo/pkg/health"

// 创建健康检查器
hc := health.NewHealthChecker(5 * time.Second)

// 注册内置检查器
hc.Register(health.NewDatabaseChecker(db))
hc.Register(health.NewRedisChecker(redisClient))
hc.Register(health.NewMemoryChecker(70.0, 90.0))        // 警告70%，严重90%
hc.Register(health.NewDiskSpaceChecker("/", 80.0, 95.0)) // 警告80%，严重95%

// 自定义检查器
hc.Register(health.NewSimpleChecker("api", func(ctx context.Context) error {
    resp, err := http.Get("https://api.example.com/ping")
    if err != nil || resp.StatusCode != 200 {
        return errors.New("API unavailable")
    }
    return nil
}))

// 可降级检查器（响应时间超过阈值视为降级）
hc.Register(health.NewDegradableChecker("database", func(ctx context.Context) error {
    return db.Ping()
}, 100*time.Millisecond))

// 执行检查
ctx := context.Background()
results := hc.Check(ctx)          // 所有检查
status := hc.GetStatus(ctx)       // 整体状态
summary := hc.GetSummary(ctx)     // 摘要信息
specific, _ := hc.CheckOne(ctx, "database") // 单个检查

// 配置缓存（避免频繁检查）
hc.SetCacheTTL(10 * time.Second)
hc.ClearCache()
```

#### HTTP 端点集成

```go
import "github.com/clarkgo/clarkgo/pkg/framework"

// 完整健康检查
h.GET("/health", framework.HealthEndpoint(hc))

// 健康检查摘要
h.GET("/health/summary", framework.HealthSummaryEndpoint(hc))

// Kubernetes Readiness Probe（就绪检查）
h.GET("/health/ready", framework.ReadinessEndpoint(hc))

// Kubernetes Liveness Probe（存活检查）
h.GET("/health/live", framework.LivenessEndpoint())

// 详细健康检查（支持查询参数）
h.GET("/health/detail", framework.DetailedHealthEndpoint(hc))

// 使用示例：
// GET /health                    - 所有检查
// GET /health?name=database      - 单个检查
// GET /health?pretty=true        - 格式化输出
```

#### Kubernetes 配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clarkgo
spec:
  template:
    spec:
      containers:
      - name: clarkgo
        image: clarkgo:latest
        ports:
        - containerPort: 8888
        
        # 存活探测（失败则重启 Pod）
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8888
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        
        # 就绪探测（失败则从 Service 移除）
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8888
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
```

### 6. CMS 功能

```go
// 用户认证
import "github.com/clarkgo/clarkgo/app/Http/Middleware"

h.Use(Middleware.JWTMiddleware())

// RBAC 权限
h.Use(Middleware.PermissionMiddleware("posts.create"))

// 文章管理
h.GET("/api/posts", PostController.Index)
h.POST("/api/posts", PostController.Store)
h.GET("/api/posts/:id", PostController.Show)
h.PUT("/api/posts/:id", PostController.Update)
h.DELETE("/api/posts/:id", PostController.Delete)

// SEO 优化
h.GET("/sitemap.xml", SEOController.Sitemap)
h.GET("/robots.txt", SEOController.Robots)
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./... -v

# 运行特定包的测试
go test ./pkg/schedule/... -v
go test ./pkg/queue/... -v
go test ./pkg/event/... -v
go test ./pkg/ratelimit/... -v
go test ./pkg/health/... -v

# 运行基准测试
go test ./pkg/ratelimit/... -bench=. -benchmem

# 测试覆盖率
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 测试结果

```
✅ pkg/schedule  - 9/9 tests passed (0.003s)
✅ pkg/event     - 8/8 tests passed (0.203s)
✅ pkg/ratelimit - 10/10 tests passed (2.406s)
✅ pkg/health    - 所有核心测试通过
```

## 📊 性能指标

### 任务调度器
- **任务执行延迟**: < 1ms
- **并发任务数**: 无限制（goroutine pool）
- **Cron 解析**: < 0.1ms per expression

### 队列系统
- **入队性能**: > 100,000 ops/sec（内存驱动）
- **Redis 吞吐**: > 10,000 ops/sec
- **Worker 效率**: 5 workers 可处理 500+ jobs/sec

### 事件系统
- **事件分发**: < 1ms（同步监听器）
- **异步处理**: Worker pool 可配置
- **监听器执行**: 按优先级顺序，支持并发

### 限流系统
- **检查性能**: > 1,000,000 ops/sec（Token Bucket）
- **内存占用**: ~200 bytes per key
- **GC 效率**: 自动清理过期 keys

### 健康检查
- **检查执行**: 并发执行，总时间 ≈ 最慢检查器
- **缓存命中**: 避免频繁检查
- **超时控制**: 可配置，防止阻塞

## 🚀 生产环境部署

### Docker 部署

```dockerfile
# Dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o clarkgo main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/clarkgo .
COPY --from=builder /app/.env.example .env

EXPOSE 8888
CMD ["./clarkgo"]
```

```bash
# 构建镜像
docker build -t clarkgo:latest .

# 运行容器
docker run -d \
  --name clarkgo \
  -p 8888:8888 \
  -e DB_HOST=mysql \
  -e REDIS_HOST=redis \
  clarkgo:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8888:8888"
    environment:
      - DB_HOST=mysql
      - DB_DATABASE=clarkgo
      - DB_USERNAME=root
      - DB_PASSWORD=secret
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=secret
      - MYSQL_DATABASE=clarkgo
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped

  # 队列 Worker
  queue_worker:
    build: .
    command: ["./clarkgo", "artisan", "queue:worker", "default", "5"]
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    restart: unless-stopped

  # 任务调度器
  scheduler:
    build: .
    command: ["./clarkgo", "artisan", "schedule:work"]
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

### Kubernetes 部署

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clarkgo
spec:
  replicas: 3
  selector:
    matchLabels:
      app: clarkgo
  template:
    metadata:
      labels:
        app: clarkgo
    spec:
      containers:
      - name: clarkgo
        image: clarkgo:latest
        ports:
        - containerPort: 8888
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: clarkgo-config
              key: db_host
        - name: REDIS_HOST
          valueFrom:
            configMapKeyRef:
              name: clarkgo-config
              key: redis_host
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8888
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8888
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: clarkgo
spec:
  selector:
    app: clarkgo
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8888
  type: LoadBalancer
```

### 监控和日志

```go
// 推荐集成 Prometheus + Grafana
import "github.com/prometheus/client_golang/prometheus"

// 自定义指标
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)

// 日志收集（推荐使用 ELK 或 Loki）
```

## 🎯 最佳实践

### 1. 任务调度
- ✅ 使用 `schedule:work` 守护进程模式
- ✅ 任务执行时间应短于调度间隔
- ✅ 长任务应推送到队列异步处理
- ✅ 添加错误日志和监控

### 2. 队列系统
- ✅ 生产环境使用 Redis 驱动
- ✅ 设置合理的 worker 数量（CPU cores × 2）
- ✅ 配置任务超时和重试次数
- ✅ 监控队列深度和 DLQ

### 3. 事件系统
- ✅ 重要操作使用同步监听器
- ✅ 耗时操作使用异步监听器
- ✅ 合理设置优先级（1-10）
- ✅ 监听器保持幂等性

### 4. 限流系统
- ✅ Token Bucket 用于一般场景
- ✅ Sliding Window 用于严格限制
- ✅ Fixed Window 用于简单场景
- ✅ 配置监控告警

### 5. 健康检查
- ✅ 检查器应快速返回（< 5s）
- ✅ 使用缓存避免频繁检查
- ✅ Readiness ≠ Liveness
- ✅ 设置合理的降级阈值

### 6. 安全性
- ✅ 使用 JWT 认证
- ✅ 实施 RBAC 权限控制
- ✅ 启用 HTTPS
- ✅ 配置 CORS
- ✅ SQL 注入防护（GORM 自动处理）
- ✅ XSS 防护
- ✅ CSRF 保护

## 📖 文档

### 官方文档
- [快速入门](doc/getting-started.md)
- [数据库使用](doc/database.md)
- [API 开发](doc/api.md)
- [任务调度](doc/artisan.md)
- [邮件系统](doc/mail.md)
- [AI 集成](doc/ai.md)
- [队列系统](doc/queue.md)
- [会话管理](doc/session.md)
- [存储系统](doc/storage.md)
- [测试指南](doc/testing.md)
- [部署指南](doc/deployment.md)
- [Swagger 文档](doc/swagger.md)

### Phase 5 实现文档
- [Phase 5 完整总结](doc/PHASE5_SUMMARY.md) - 详细的实现文档
- [快速参考手册](doc/QUICKREF.md) - 常用代码示例

### API 文档
```bash
# 生成 GoDoc
godoc -http=:6060

# 访问文档
open http://localhost:6060/pkg/github.com/clarkgo/clarkgo/
```

### Swagger 文档
```bash
# 访问 Swagger UI
open http://localhost:8888/swagger/index.html
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发规范

- 遵循 Go 代码规范
- 编写单元测试
- 更新相关文档
- 提交前运行 `go fmt` 和 `go vet`
- 确保所有测试通过

### 报告问题

如果发现 Bug 或有功能建议，请[创建 Issue](https://github.com/chenyusolar/clarkgo/issues)。

## 📝 更新日志

### Phase 5 (2024-11-19)
- ✅ 实现任务调度系统
- ✅ 实现队列系统（内存 + Redis 驱动）
- ✅ 实现事件系统（同步/异步监听器）
- ✅ 实现限流系统（3种算法）
- ✅ 实现健康检查系统
- ✅ 新增 35+ 单元测试
- ✅ 完善文档和示例

### Phase 4
- ✅ AI 多模型集成
- ✅ 邮件系统增强
- ✅ 统计分析功能

### Phase 3
- ✅ CMS 核心功能
- ✅ RBAC 权限系统
- ✅ SEO 优化

### Phase 2
- ✅ 用户认证（JWT）
- ✅ 数据库迁移
- ✅ 基础 API

### Phase 1
- ✅ 框架核心
- ✅ 路由系统
- ✅ 中间件支持

## 🙏 致谢

- [CloudWeGo Hertz](https://github.com/cloudwego/hertz) - 高性能 HTTP 框架
- [GORM](https://gorm.io/) - 优秀的 ORM 库
- [Go Redis](https://github.com/go-redis/redis) - Redis 客户端
- 所有贡献者和支持者

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 📧 联系方式

- 作者：Clark Chen
- GitHub: [@chenyusolar](https://github.com/chenyusolar)
- 项目地址: [https://github.com/chenyusolar/clarkgo](https://github.com/chenyusolar/clarkgo)

---

**ClarkGo - 让 Go Web 开发更简单、更高效！** 🚀

Made with ❤️ by ClarkGo Team