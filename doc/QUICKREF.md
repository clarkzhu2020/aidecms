# AideCMS 快速参考手册

## Artisan 命令速查

### 任务调度 (Schedule)
```bash
# 启动调度器守护进程
artisan schedule:work

# 运行一次所有到期任务
artisan schedule:run

# 列出所有已注册任务
artisan schedule:list
```

### 队列系统 (Queue)
```bash
# 启动队列 worker (5个并发)
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

### 事件系统 (Event)
```bash
# 运行事件系统演示
artisan event:test

# 列出注册的事件和监听器
artisan event:list

# 显示事件统计信息
artisan event:stats
```

### 限流系统 (Rate Limit)
```bash
# 运行限流演示
artisan ratelimit demo
```

### 健康检查 (Health Check)
```bash
# 运行健康检查演示
artisan health demo
```

---

## 代码示例速查

### 1. 任务调度

#### 基础用法
```go
import "github.com/chenyusolar/aidecms/pkg/schedule"

// 每分钟执行
schedule.EveryMinute().Do(func() {
    log.Println("Task executed")
})

// 每小时执行
schedule.Hourly().Do(hourlyTask)

// 每天特定时间执行
schedule.Daily().At("02:30").Do(backupTask)

// 每周特定时间执行
schedule.WeeklyOn(time.Monday).At("09:00").Do(weeklyReport)

// 每月特定日期执行
schedule.MonthlyOn(1).At("00:00").Do(monthlyTask)
```

#### 高级用法
```go
// 使用 Cron 表达式
schedule.Cron("*/5 * * * *").Do(task) // 每5分钟

// 自定义任务名称
schedule.EveryMinute().Name("my_task").Do(task)

// 手动触发任务
schedule.RunTask("my_task")

// 获取任务统计
stats := schedule.GetStats()
```

---

### 2. 队列系统

#### 基础用法
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

// 推送任务
queue.Push("emails", &EmailJob{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Welcome to our service!",
})

// 延迟任务
queue.DelayUntil("notifications", job, time.Now().Add(1*time.Hour))

// 启动 workers
ctx := context.Background()
queue.Work(ctx, "emails", 5) // 5个并发 workers
```

#### Redis 配置
```go
import "github.com/chenyusolar/aidecms/pkg/queue"

// 使用 Redis 驱动
redisQueue := queue.NewRedisDriver(redisClient, "myapp")
queue.UseDriver("default", redisQueue)

// 推送任务到 Redis 队列
queue.Push("default", job)
```

---

### 3. 事件系统

#### 基础用法
```go
import "github.com/chenyusolar/aidecms/pkg/event"

// 注册监听器 (同步)
event.Listen("user.registered", func(e event.Event) error {
    user := e.(*event.UserRegistered)
    log.Printf("New user: %s", user.Email)
    return sendWelcomeEmail(user)
})

// 注册异步监听器
event.ListenAsync("order.created", func(e event.Event) error {
    order := e.(*event.OrderCreated)
    return processOrder(order)
}, event.WithPriority(5))

// 触发事件
event.Dispatch(event.NewUserRegistered(
    userID,
    "user@example.com",
    "John Doe",
))
```

#### 自定义事件
```go
// 定义自定义事件
type ProductCreatedEvent struct {
    event.BaseEvent
    ProductID string
    Name      string
    Price     float64
}

func (e *ProductCreatedEvent) EventName() string {
    return "product.created"
}

// 触发自定义事件
event.Dispatch(&ProductCreatedEvent{
    ProductID: "prod-123",
    Name:      "New Product",
    Price:     99.99,
})

// 监听自定义事件
event.Listen("product.created", func(e event.Event) error {
    product := e.(*ProductCreatedEvent)
    return updateInventory(product)
})
```

---

### 4. 限流系统

#### Token Bucket (令牌桶)
```go
import "github.com/chenyusolar/aidecms/pkg/ratelimit"

// 创建限流器: 每秒100个请求，突发200
limiter := ratelimit.NewTokenBucket(100, 200)

// 检查是否允许
if limiter.Allow(userID) {
    // 处理请求
} else {
    // 返回 429 Too Many Requests
}
```

#### Sliding Window (滑动窗口)
```go
// 每分钟1000个请求
limiter := ratelimit.NewSlidingWindow(1000, 1*time.Minute)

if limiter.Allow(userID) {
    // 处理请求
}

// 获取统计
stats := limiter.GetStats(userID)
fmt.Printf("Used: %d/%d\n", stats["requests"], stats["limit"])
```

#### Fixed Window (固定窗口)
```go
// 每小时5000个请求
limiter := ratelimit.NewFixedWindow(5000, 1*time.Hour)

if limiter.Allow(userID) {
    // 处理请求
}
```

#### Hertz 中间件
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
```

---

### 5. 健康检查

#### 基础用法
```go
import "github.com/chenyusolar/aidecms/pkg/health"

// 创建健康检查器
hc := health.NewHealthChecker(5 * time.Second)

// 注册内置检查器
hc.Register(health.NewDatabaseChecker(db))
hc.Register(health.NewRedisChecker(redisClient))
hc.Register(health.NewMemoryChecker(70.0, 90.0))

// 执行检查
ctx := context.Background()
results := hc.Check(ctx)
status := hc.GetStatus(ctx)
```

#### 自定义检查器
```go
// 简单检查器
checker := health.NewSimpleChecker("my_service", func(ctx context.Context) error {
    // 检查逻辑
    if !serviceAvailable() {
        return errors.New("service unavailable")
    }
    return nil
})

hc.Register(checker)

// 带详细信息的检查器
checker := health.NewSimpleChecker("api", func(ctx context.Context) error {
    return checkAPI()
}).WithDetails(func(ctx context.Context) map[string]interface{} {
    return map[string]interface{}{
        "version": "1.0.0",
        "uptime":  getUptime(),
    }
})

// 可降级检查器 (慢响应视为降级)
checker := health.NewDegradableChecker("database", func(ctx context.Context) error {
    return db.Ping()
}, 100*time.Millisecond) // 超过100ms视为降级
```

#### Hertz 端点集成
```go
import "github.com/chenyusolar/aidecms/pkg/framework"

// 完整健康检查
h.GET("/health", framework.HealthEndpoint(hc))

// 摘要信息
h.GET("/health/summary", framework.HealthSummaryEndpoint(hc))

// Kubernetes readiness probe
h.GET("/health/ready", framework.ReadinessEndpoint(hc))

// Kubernetes liveness probe
h.GET("/health/live", framework.LivenessEndpoint())
```

#### Kubernetes 配置
```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: clarkgo
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
```

---

## HTTP API 端点

### 健康检查端点
```
GET /health                 - 完整健康检查
GET /health/summary         - 健康检查摘要
GET /health/ready           - 就绪检查 (K8s readiness)
GET /health/live            - 存活检查 (K8s liveness)
GET /health?name=database   - 单个检查器详情
GET /health?pretty=true     - 格式化输出
```

### 响应格式

#### /health
```json
{
  "status": "healthy",
  "checks": {
    "database": {
      "name": "database",
      "status": "healthy",
      "message": "Database is healthy",
      "timestamp": "2024-01-01T00:00:00Z",
      "duration": "10ms",
      "details": {
        "open_connections": 5,
        "in_use": 2,
        "idle": 3
      }
    }
  }
}
```

#### /health/summary
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "total_checks": 3,
  "healthy_count": 3,
  "degraded_count": 0,
  "unhealthy_count": 0
}
```

---

## 常用配置

### 环境变量
```bash
# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_NAME=aidecms
DB_USER=root
DB_PASSWORD=password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 应用配置
APP_ENV=production
APP_DEBUG=false
APP_PORT=8888
```

### 配置文件示例
```go
// config/app.go
type Config struct {
    Database struct {
        Host     string
        Port     int
        Name     string
        User     string
        Password string
    }
    
    Redis struct {
        Host     string
        Port     int
        Password string
        DB       int
    }
    
    Queue struct {
        Driver     string // "memory" or "redis"
        Workers    int
        RetryTimes int
    }
    
    RateLimit struct {
        Algorithm string // "token_bucket", "sliding_window", "fixed_window"
        Rate      int
        Capacity  int
    }
}
```

---

## 故障排查

### 常见问题

#### 1. 任务调度不执行
- 检查 `schedule:work` 是否运行
- 验证 Cron 表达式格式
- 查看日志: `storage/logs/artisan.log`

#### 2. 队列任务失败
- 检查 worker 是否启动: `queue:worker`
- 查看队列状态: `queue:status`
- 检查死信队列: `queue:stats`

#### 3. 限流不生效
- 验证中间件注册顺序
- 检查 key 生成逻辑
- 查看限流统计

#### 4. 健康检查超时
- 调整检查器超时时间
- 优化检查逻辑
- 使用缓存减少检查频率

---

## 性能优化建议

### 1. 任务调度
- 避免在调度任务中执行耗时操作
- 长时间任务推送到队列
- 合理设置任务执行间隔

### 2. 队列系统
- 生产环境使用 Redis 驱动
- Worker 数量 = CPU cores * 2
- 配置合理的超时和重试

### 3. 事件系统
- 重要操作使用同步监听器
- 耗时操作使用异步监听器
- 监听器保持无状态

### 4. 限流系统
- 根据场景选择合适算法
- Token Bucket 适合大多数场景
- 配置监控和告警

### 5. 健康检查
- 检查器应快速返回 (< 5s)
- 使用缓存避免频繁检查
- 设置合理降级阈值

---

## 监控指标

### 推荐监控项
- 任务执行成功率
- 队列深度和延迟
- 事件处理时间
- 限流触发次数
- 健康检查失败次数

### Prometheus 集成 (待实现)
```go
// 任务调度指标
schedule_task_executions_total
schedule_task_failures_total
schedule_task_duration_seconds

// 队列指标
queue_depth
queue_processing_duration_seconds
queue_dlq_size

// 限流指标
ratelimit_requests_total
ratelimit_denied_total

// 健康检查指标
health_check_status
health_check_duration_seconds
```

---

## 更多资源

- 完整文档: `doc/PHASE5_SUMMARY.md`
- API 文档: `godoc -http=:6060`
- 源码: https://github.com/chenyusolar/aidecms
- 问题反馈: GitHub Issues

---

**AideCMS CMS Framework - 快速、可靠、易用** 🚀
