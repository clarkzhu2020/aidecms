# 队列系统

AideCMS 提供了强大的队列系统，支持邮件发送、任务处理等异步操作。

## 配置

队列配置在 `.env` 文件中:

```
QUEUE_CONNECTION=redis  # 或 database
QUEUE_RETRY_AFTER=90    # 失败后重试间隔(秒)
QUEUE_TRIES=3           # 最大重试次数
```

## 基本使用

### 创建任务

```go
import "internal/app/jobs"

// 创建邮件任务
jobs.NewEmailJob("welcome", "user@example.com", map[string]string{
    "name": "John",
}).Dispatch()
```

### 任务类示例

`internal/app/jobs/email_job.go`:

```go
package jobs

type EmailJob struct {
    Template  string
    Recipient string
    Variables map[string]string
}

func (j *EmailJob) Handle() error {
    return mail.To(j.Recipient).
        Template(j.Template, j.Variables).
        Send()
}
```

## 命令行操作

### 启动队列处理器

```bash
go run cmd/artisan/main.go queue:work
```

### 暂停队列

```bash
go run cmd/artisan/main.go queue:pause
# 或定时暂停(30分钟后)
go run cmd/artisan/main.go queue:pause 30m
```

### 恢复队列

```bash
go run cmd/artisan/main.go queue:resume
# 或定时恢复(1小时后)
go run cmd/artisan/main.go queue:resume 1h
```

### 查看队列状态

```bash
go run cmd/artisan/main.go queue:status
```

输出示例:
```
Email Queue Status:
ID        Subject              Status     Created        Retries  Priority
a1b2c3d4  Welcome email        pending    2023-01-01     0        3
e5f6g7h8  Password reset       failed     2023-01-02     2        1
```

### 队列统计

```bash
go run cmd/artisan/main.go queue:stats
```

输出示例:
```
Queue Statistics:
Total jobs:    15
Pending:       5
Sent:          8
Failed:        2
Avg send time: 2.5s
Current rate:  2s per email
```

## 高级功能

### 优先级

任务默认优先级为3(中等)，可设置1-5(1最高):

```go
job := jobs.NewEmailJob(...)
job.Priority = 1 // 最高优先级
job.Dispatch()
```

或使用命令行:

```bash
go run cmd/artisan/main.go queue:priority a1b2c3d4 1
```

### 速率限制

队列会自动根据发送成功率调整速率:

- 成功率 > 90%: 加快发送
- 成功率 < 70%: 减慢发送

可在 `.env` 中配置:

```
QUEUE_MIN_INTERVAL=1s    # 最快1秒1封
QUEUE_MAX_INTERVAL=30s   # 最慢30秒1封
```

### 失败处理

任务失败后会自动重试(最多3次)，失败任务会保留7天。

手动重试所有失败任务:

```bash
go run cmd/artisan/main.go queue:retry
```

清理已发送的旧任务(7天前):

```bash
go run cmd/artisan/main.go queue:clean