# AideCMS - 企业级 Go CMS 平台框架

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Hertz](https://img.shields.io/badge/Hertz-CloudWeGo-blue)](https://github.com/cloudwego/hertz)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

AideCMS 是一个基于 CloudWeGo Hertz 框架开发的企业级 CMS 平台框架，提供完整的任务调度、队列系统、事件驱动、限流保护和健康监控等核心功能。

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

### 🌐 Web3 区块链集成
- ✅ 多链支持（Bitcoin、Ethereum、BSC、Solana）
- ✅ 地址余额查询
- ✅ 交易信息查询
- ✅ 区块高度查询
- ✅ 钱包信息查询
- ✅ 地址格式验证
- ✅ Gas 价格查询（EVM 链）
- ✅ SPL Token 支持（Solana）

### 💱 加密货币交易所集成
- ✅ 多交易所支持（Coinbase、KuCoin、Hyperliquid）
- ✅ 中心化交易所（CEX）和去中心化交易所（DEX）
- ✅ 账户余额查询
- ✅ 交易对价格查询
- ✅ 持仓信息查询（Hyperliquid）
- ✅ 资金费率查询（Hyperliquid）
- ✅ 跨交易所价格比较
- ✅ 安全的 API 签名认证
- ✅ HTTP API 和 CLI 命令

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
git clone https://github.com/zhuclark2020/aidecms.git
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
    "github.com/chenyusolar/aidecms/pkg/framework"
)

func main() {
    app := framework.NewApplication()
    
    // 注册路由
    app.RegisterRoutes(func(router *framework.Router) {
        router.GET("/", func(ctx context.Context, c *framework.RequestContext) {
            c.JSON(200, map[string]interface{}{
                "message": "Welcome to AideCMS!",
            })
        })
    })
    
    // 启动服务器
    app.Run(":8888")
}
```

### 2. 任务调度

```go
import "github.com/chenyusolar/aidecms/pkg/schedule"

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
import "github.com/chenyusolar/aidecms/pkg/queue"

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
import "github.com/chenyusolar/aidecms/pkg/event"

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
    "github.com/chenyusolar/aidecms/pkg/ratelimit"
    "github.com/chenyusolar/aidecms/pkg/framework"
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
import "github.com/chenyusolar/aidecms/pkg/health"

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

### 7. Web3 区块链集成

```go
import (
    "github.com/chenyusolar/aidecms/pkg/web3"
    "context"
    "time"
)

// 初始化 Web3 客户端
web3.InitializeClients()
manager := web3.GetManager()

// 查询 Ethereum 余额
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

address := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"
balance, err := manager.GetBalance(ctx, web3.Ethereum, address)
if err != nil {
    panic(err)
}
log.Printf("ETH Balance: %s wei", balance)

// 查询交易信息
tx, err := manager.GetTransaction(ctx, web3.Ethereum, txHash)
if err != nil {
    panic(err)
}
log.Printf("Transaction: %+v", tx)

// 多链余额查询
addresses := web3.MultiChainAddress{
    Bitcoin:  "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    Ethereum: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
    BSC:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
    Solana:   "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK",
}
balances, _ := addresses.GetAllBalances(ctx)
for chain, balance := range balances {
    log.Printf("%s: %s", chain, balance)
}
```

### 8. 加密货币交易所集成

```go
import (
    "github.com/chenyusolar/aidecms/pkg/web3"
    "context"
)

// 初始化交易所客户端
web3.InitializeClients()
manager := web3.GetExchangeManager()

ctx := context.Background()

// 查询 Coinbase 余额
balance, err := manager.GetBalance(ctx, web3.ExchangeCoinbase, "BTC")
if err != nil {
    panic(err)
}
log.Printf("BTC Balance: %s", balance)

// 查询交易对价格
price, err := manager.GetPrice(ctx, web3.ExchangeCoinbase, "BTC-USD")
if err != nil {
    panic(err)
}
log.Printf("BTC Price: $%s", price)

// 跨交易所价格比较
prices, _ := web3.GetAllExchangePrices(ctx, "BTC-USD")
for exchange, price := range prices {
    log.Printf("%s: $%s", exchange, price)
}
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
│   ├── validator/             # 验证器
│   └── web3/                  # Web3 区块链集成 ⭐
│       ├── web3.go            # 核心接口
│       ├── ethereum.go        # Ethereum/BSC 客户端
│       ├── bitcoin.go         # Bitcoin 客户端
│       ├── solana.go          # Solana 客户端
│       ├── coinbase.go        # Coinbase 交易所 ⭐
│       ├── kucoin.go          # KuCoin 交易所 ⭐
│       ├── hyperliquid.go     # Hyperliquid DEX ⭐
│       ├── exchange.go        # 交易所管理器 ⭐
│       ├── config.go          # 配置管理
│       ├── token.go           # 代币相关
│       └── web3_test.go       # 测试
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

AideCMS 提供了强大的 Artisan 命令行工具，用于开发和管理。

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

### Web3 区块链命令
```bash
# 初始化 Web3 客户端
artisan web3 init

# 查看支持的链
artisan web3 chains

# 查询地址余额
artisan web3 balance ethereum 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0
artisan web3 balance bitcoin 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
artisan web3 balance solana 7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK

# 查询交易信息
artisan web3 transaction ethereum 0x1234...

# 获取最新区块
artisan web3 block ethereum

# 验证地址格式
artisan web3 validate ethereum 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0

# 获取钱包信息
artisan web3 wallet ethereum 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0
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
MAIL_FROM_NAME=AideCMS

# AI 配置（可选）
AI_PROVIDER=openai
AI_API_KEY=sk-xxx
AI_MODEL=gpt-3.5-turbo
```

### 配置读取

```go
import "github.com/chenyusolar/aidecms/pkg/config"

// 获取字符串配置
host := config.GetEnv("DB_HOST", "localhost")

// 获取整数配置
port := config.GetEnvInt("DB_PORT", 3306)

// 获取布尔配置
debug := config.GetEnvBool("APP_DEBUG", false)
```


## 🤖 AI 大模型与对话功能（BiruAI）

AideCMS 集成 CloudWeGo Eino 框架，内置强大的 AI 大模型能力，支持 OpenAI、Anthropic、豆包、通义千问、ChatGLM 等主流模型，统一接口，支持多轮对话、流式输出、嵌入向量、灵活配置。

### 主要特性
- 多模型支持：OpenAI、Anthropic、豆包、通义千问、ChatGLM
- 统一 API：/api/ai/chat、/completion、/embedding、/conversation
- 多轮对话与上下文管理
- 流式输出（SSE）与嵌入向量生成
- 命令行一键配置与测试
- 灵活配置与热切换

### 快速上手
```bash
# 配置 OpenAI
go run cmd/artisan/main.go ai:setup openai sk-your-api-key gpt-4
# 配置千问
go run cmd/artisan/main.go ai:setup qianwen your-api-key qwen-max
# 测试连接
go run cmd/artisan/main.go ai:test
# 聊天
go run cmd/artisan/main.go ai:chat "你好，介绍下AideCMS"
```

### 主要 API 路由
- POST `/api/ai/chat`         # 聊天/对话
- POST `/api/ai/completion`   # 文本补全
- POST `/api/ai/embedding`    # 向量嵌入
- POST `/api/ai/conversation` # 多轮对话
- GET  `/api/ai/models`       # 支持模型列表

#### 聊天 API 示例
```bash
curl -X POST http://localhost:8888/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下 Go 语言", "model": "qianwen"}'
```

#### 嵌入向量 API 示例
```bash
curl -X POST http://localhost:8888/api/ai/embedding \
  -H "Content-Type: application/json" \
  -d '{"input": ["Hello world", "你好世界"], "model": "openai"}'
```

### 代码集成示例
```go
import "github.com/clarkzhu2020/aidecms/pkg/ai"
config := &ai.Config{Provider: "openai", APIKey: "sk-xxx", Model: "gpt-4"}
client, _ := ai.NewClient(config)
resp, _ := client.CreateCompletion(context.Background(), "请介绍Go语言")
fmt.Println(resp)
```

### 配置与管理
- 命令行管理：`ai:config list|show|delete|default`
- 配置文件：`config/ai/openai.json`、`config/ai/qianwen.json` 等
- 支持环境变量存储敏感信息

### 支持模型一览
| 提供商     | 典型模型           | 能力         |
|------------|--------------------|--------------|
| OpenAI     | gpt-4, gpt-3.5     | 聊天/补全/嵌入 |
| Anthropic  | claude-3-opus等    | 聊天/补全     |
| 豆包       | ep-xxx             | 聊天/补全/嵌入 |
| 千问       | qwen-max等         | 聊天/补全/嵌入 |
| ChatGLM    | glm-4等            | 聊天/补全     |

### 最佳实践
- 推荐使用流式输出提升体验
- 合理设置超时与重试
- 监控 API 调用与日志
- 敏感信息用环境变量管理

详细用法见 [AI 集成文档](doc/ai.md)

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
import "github.com/chenyusolar/aidecms/pkg/queue"

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
    "github.com/chenyusolar/aidecms/pkg/framework"
    "github.com/chenyusolar/aidecms/pkg/ratelimit"
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
import "github.com/chenyusolar/aidecms/pkg/health"

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
import "github.com/chenyusolar/aidecms/pkg/framework"

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

### 6. Web3 区块链集成

支持 Bitcoin、Ethereum、BSC、Solana 等多条公链：

```go
import "github.com/chenyusolar/aidecms/pkg/web3"

// 初始化所有区块链客户端
if err := web3.InitializeClients(); err != nil {
    log.Fatal(err)
}

manager := web3.GetManager()
ctx := context.Background()

// ===== 查询地址余额 =====
// Ethereum
ethBalance, _ := manager.GetBalance(ctx, web3.Ethereum, "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
fmt.Printf("ETH Balance: %s wei\n", ethBalance)

// Bitcoin
btcBalance, _ := manager.GetBalance(ctx, web3.Bitcoin, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
fmt.Printf("BTC Balance: %s satoshi\n", btcBalance)

// BSC (Binance Smart Chain)
bscBalance, _ := manager.GetBalance(ctx, web3.BSC, "0x...")
fmt.Printf("BNB Balance: %s wei\n", bscBalance)

// Solana
solBalance, _ := manager.GetBalance(ctx, web3.Solana, "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK")
fmt.Printf("SOL Balance: %s lamports\n", solBalance)

// ===== 查询交易信息 =====
tx, err := manager.GetTransaction(ctx, web3.Ethereum, "0xabc...")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("From: %s\n", tx.From)
fmt.Printf("To: %s\n", tx.To)
fmt.Printf("Value: %s\n", tx.Value)
fmt.Printf("Status: %s\n", tx.Status)

// ===== 查询区块高度 =====
client, _ := manager.GetClient(web3.Ethereum)
blockNumber, _ := client.GetBlockNumber(ctx)
fmt.Printf("Latest Block: %d\n", blockNumber)

// ===== 查询钱包完整信息 =====
walletInfo, err := web3.GetWalletInfo(ctx, web3.Ethereum, "0x742d35...")
fmt.Printf("Balance: %s\n", walletInfo.Balance)
fmt.Printf("Nonce: %d\n", walletInfo.Nonce)
fmt.Printf("Code: %s\n", walletInfo.Code)

// ===== Gas 价格查询（EVM 链）=====
gasPrice, _ := client.GetGasPrice(ctx)
fmt.Printf("Gas Price: %s gwei\n", gasPrice)

// ===== 地址验证 =====
if err := web3.ValidateAddress(web3.Bitcoin, "1A1zP1eP..."); err == nil {
    fmt.Println("✅ Valid Bitcoin address")
}

// ===== 多链余额批量查询 =====
addresses := web3.MultiChainAddress{
    Bitcoin:  "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    Ethereum: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
    BSC:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
    Solana:   "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK",
}

balances, err := addresses.GetAllBalances(ctx)
for chain, balance := range balances {
    fmt.Printf("%s: %s\n", chain, balance)
}
```

**环境变量配置**：

```env
# Ethereum
WEB3_ETHEREUM_RPC=https://mainnet.infura.io/v3/YOUR-PROJECT-ID

# BSC (Binance Smart Chain)
WEB3_BSC_RPC=https://bsc-dataseed.binance.org/

# Bitcoin
WEB3_BITCOIN_RPC=https://bitcoin-mainnet.core.chainstack.com
WEB3_BITCOIN_API_KEY=your_api_key

# Solana
WEB3_SOLANA_RPC=https://api.mainnet-beta.solana.com
```

**HTTP API 端点**：

```bash
# 查询余额
curl http://localhost:8080/api/web3/ethereum/balance/0x742d35...

# 查询交易
curl http://localhost:8080/api/web3/ethereum/transaction/0xabc...

# 查询区块高度
curl http://localhost:8080/api/web3/ethereum/block-number

# 查询钱包信息
curl http://localhost:8080/api/web3/ethereum/wallet/0x742d35...

# 验证地址
curl http://localhost:8080/api/web3/bitcoin/validate/1A1zP1eP...

# 支持的链列表
curl http://localhost:8080/api/web3/chains
```

**CLI 命令**：

```bash
# 初始化 Web3 客户端
artisan web3 init

# 查询余额
artisan web3 balance ethereum 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0
artisan web3 balance bitcoin 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa

# 查询交易
artisan web3 transaction ethereum 0xabc123...

# 查询区块高度
artisan web3 block ethereum

# 查询钱包信息
artisan web3 wallet ethereum 0x742d35...

# 验证地址
artisan web3 validate bitcoin 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa

# 列出支持的链
artisan web3 chains
```

### 7. 加密货币交易所集成

支持中心化交易所（CEX）和去中心化交易所（DEX）：

#### 支持的交易所

| 交易所 | 类型 | 特点 |
|--------|------|------|
| **Coinbase** | CEX | 美国合规交易所，适合法币出入金 |
| **KuCoin** | CEX | 币种丰富，手续费低 |
| **Hyperliquid** | DEX | 永续合约，高杠杆（50x），非托管 |

#### 基本使用

```go
import "github.com/chenyusolar/aidecms/pkg/web3"

// 初始化交易所客户端
if err := web3.InitializeClients(); err != nil {
    log.Fatal(err)
}

manager := web3.GetExchangeManager()
ctx := context.Background()

// ===== 查询单个币种余额 =====
// Coinbase
balance, _ := manager.GetBalance(ctx, web3.Coinbase, "BTC")
fmt.Printf("Coinbase BTC: %s\n", balance)

// KuCoin
balance, _ = manager.GetBalance(ctx, web3.KuCoin, "ETH")
fmt.Printf("KuCoin ETH: %s\n", balance)

// Hyperliquid (DEX)
balance, _ = manager.GetBalance(ctx, web3.Hyperliquid, "USDC")
fmt.Printf("Hyperliquid USDC: %s\n", balance)

// ===== 查询所有余额 =====
balances, err := manager.GetBalances(ctx, web3.Coinbase)
if err != nil {
    log.Fatal(err)
}
for currency, amount := range balances {
    fmt.Printf("%s: %s\n", currency, amount)
}
// 输出:
// BTC: 0.5
// ETH: 10.0
// USD: 50000.0

// ===== 查询交易对价格 =====
// Coinbase
price, _ := manager.GetPrice(ctx, web3.Coinbase, "BTC-USD")
fmt.Printf("Coinbase BTC-USD: $%s\n", price)

// KuCoin
price, _ = manager.GetPrice(ctx, web3.KuCoin, "ETH-USDT")
fmt.Printf("KuCoin ETH-USDT: $%s\n", price)

// Hyperliquid
price, _ = manager.GetPrice(ctx, web3.Hyperliquid, "BTC-USD")
fmt.Printf("Hyperliquid BTC-USD: $%s\n", price)

// ===== 跨交易所价格比较 =====
prices, err := web3.GetAllExchangePrices(ctx, "BTC-USD")
for exchange, price := range prices {
    fmt.Printf("%s: $%s\n", exchange, price)
}
// 输出:
// coinbase: 45678.90
// kucoin: 45680.50
// hyperliquid: 45675.20

// ===== 跨交易所余额查询 =====
balances, err = web3.GetAllExchangeBalances(ctx, "BTC")
for exchange, balance := range balances {
    fmt.Printf("%s: %s BTC\n", exchange, balance)
}

// ===== 获取支持的交易所列表 =====
exchanges := manager.GetSupportedExchanges()
fmt.Printf("Supported exchanges: %v\n", exchanges)
```

#### Hyperliquid 高级功能（DEX）

Hyperliquid 作为去中心化永续合约交易所，提供额外功能：

```go
// 创建 Hyperliquid 客户端（需要以太坊私钥）
privateKey := os.Getenv("EXCHANGE_HYPERLIQUID_PRIVATE_KEY")
client, err := web3.NewHyperliquidClient(privateKey)
if err != nil {
    log.Fatal(err)
}

// ===== 查询持仓信息 =====
positions, err := client.GetPositions(ctx)
for _, pos := range positions {
    fmt.Printf("币种: %s\n", pos.Coin)
    fmt.Printf("  数量: %s\n", pos.Size)
    fmt.Printf("  开仓价: %s\n", pos.EntryPrice)
    fmt.Printf("  持仓价值: %s\n", pos.PositionValue)
    fmt.Printf("  未实现盈亏: %s\n", pos.UnrealizedPnl)
    fmt.Printf("  杠杆: %sx\n", pos.Leverage)
    fmt.Printf("  清算价: %s\n", pos.Liquidation)
}

// ===== 查询资金费率 =====
fundingRate, _ := client.GetFundingRate(ctx, "BTC")
fmt.Printf("BTC Funding Rate: %s\n", fundingRate)

// ===== 查询24小时交易量 =====
volume, _ := client.Get24HVolume(ctx, "BTC")
fmt.Printf("BTC 24H Volume: $%s\n", volume)

// ===== 查询订单簿 =====
orderBook, _ := client.GetOrderBook(ctx, "BTC")
fmt.Printf("订单簿: %+v\n", orderBook)

// ===== 下单交易 =====
// 限价做多
order := web3.OrderRequest{
    Coin:       "BTC",
    IsBuy:      true,
    Size:       0.1,
    LimitPrice: 45000.0,
    ReduceOnly: false,
}
orderID, err := client.PlaceOrder(ctx, order)
fmt.Printf("Order ID: %s\n", orderID)

// 市价平仓
closeOrder := web3.OrderRequest{
    Coin:       "BTC",
    IsBuy:      false,
    Size:       0.1,
    LimitPrice: 0,        // 0 表示市价
    ReduceOnly: true,     // 只减仓
}
client.PlaceOrder(ctx, closeOrder)

// 取消订单
err = client.CancelOrder(ctx, "BTC", 12345)
```

**环境变量配置**：

```env
# Coinbase
EXCHANGE_COINBASE_API_KEY=your_api_key
EXCHANGE_COINBASE_API_SECRET=your_api_secret

# KuCoin
EXCHANGE_KUCOIN_API_KEY=your_api_key
EXCHANGE_KUCOIN_API_SECRET=your_api_secret
EXCHANGE_KUCOIN_PASSPHRASE=your_passphrase

# Hyperliquid DEX (使用以太坊私钥)
EXCHANGE_HYPERLIQUID_PRIVATE_KEY=your_ethereum_private_key_without_0x
EXCHANGE_HYPERLIQUID_ADDRESS=your_ethereum_address
```

**HTTP API 端点**：

```bash
# 查询余额
curl http://localhost:8080/api/exchange/coinbase/balance/BTC
curl http://localhost:8080/api/exchange/hyperliquid/balance/USDC

# 查询所有余额
curl http://localhost:8080/api/exchange/coinbase/balances

# 查询价格
curl http://localhost:8080/api/exchange/coinbase/price/BTC-USD
curl http://localhost:8080/api/exchange/hyperliquid/price/BTC-USD

# 支持的交易所列表
curl http://localhost:8080/api/exchange/supported

# 跨交易所余额查询
curl http://localhost:8080/api/exchange/all/balance/BTC

# 跨交易所价格比较
curl http://localhost:8080/api/exchange/all/price/BTC-USD
```

**CLI 命令**：

```bash
# 列出支持的交易所
artisan exchange list

# 查询余额
artisan exchange balance coinbase BTC
artisan exchange balance hyperliquid USDC

# 查询所有余额
artisan exchange balances coinbase

# 查询价格
artisan exchange price coinbase BTC-USD
artisan exchange price hyperliquid BTC-USD

# 跨交易所价格比较
artisan exchange compare BTC-USD

# 跨交易所余额查询
artisan exchange balance-all BTC
```

**CEX vs DEX 对比**：

| 特性 | CEX (Coinbase/KuCoin) | DEX (Hyperliquid) |
|------|----------------------|-------------------|
| **托管** | 托管（交易所保管） | 非托管（用户自己保管） |
| **KYC** | 需要 | 不需要 |
| **交易类型** | 现货 | 永续合约 |
| **杠杆** | 3-10x | 最高 50x |
| **认证** | API Key + Secret | 以太坊私钥签名 |
| **手续费** | 0.1-0.5% | Maker -0.02%, Taker 0.05% |
| **适用场景** | 法币出入金、简单交易 | 合约交易、高杠杆、隐私 |

**安全提示**：

- ⚠️ **CEX**: 只授予只读权限（View），不要授予提现权限
- ⚠️ **DEX**: 使用专用钱包，不要使用主钱包的私钥
- ⚠️ 私钥和 API 密钥永远不要硬编码，使用环境变量
- ⚠️ 定期检查账户活动和 API 调用记录
- ⚠️ 启用 IP 白名单（如果交易所支持）

### 8. CMS 功能

```go
// 用户认证
import "github.com/chenyusolar/aidecms/app/Http/Middleware"

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
open http://localhost:6060/pkg/github.com/chenyusolar/aidecms/
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

如果发现 Bug 或有功能建议，请[创建 Issue](https://github.com/zhuclark2020/aidecms/issues)。

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

- 作者：Clark Zhu
- GitHub: [@chenyusolar](https://github.com/chenyusolar)
- 项目地址: [https://github.com/zhuclark2020/aidecms](https://github.com/zhuclark2020/aidecms)

---

**AideCMS - 让 Go Web 开发更简单、更高效！** 🚀

Made with ❤️ by AideCMS Team