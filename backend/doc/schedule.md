# 任务调度系统

AideCMS 提供了强大的任务调度系统，支持 Cron 表达式和预定义调度。

## 功能特性

- ✅ Cron 表达式支持（分钟级精度）
- ✅ 预定义调度（每分钟、每小时、每天等）
- ✅ 流式 API（链式调用）
- ✅ 任务执行日志
- ✅ 并发安全
- ✅ 任务统计
- ✅ 手动触发任务

## 快速开始

### 1. 创建调度器

```go
import "github.com/chenyusolar/aidecms/pkg/schedule"

scheduler := schedule.NewScheduler()
scheduler.Start()
defer scheduler.Stop()
```

### 2. 注册任务

#### 使用预定义调度

```go
// 每分钟执行
scheduler.NewTask("cleanup-temp").
    EveryMinute().
    Description("Clean up temporary files").
    Do(func() error {
        // 清理逻辑
        return nil
    })

// 每天凌晨 2:00 执行
scheduler.NewTask("backup-db").
    DailyAt(2, 0).
    Description("Backup database").
    Do(func() error {
        // 备份逻辑
        return nil
    })

// 每周一上午 9:00 执行
scheduler.NewTask("weekly-report").
    WeeklyOn(1, 9, 0). // 1 = Monday
    Description("Generate weekly report").
    Do(func() error {
        // 报告生成逻辑
        return nil
    })
```

#### 使用 Cron 表达式

```go
scheduler.NewTask("custom-task").
    Cron("0 */2 * * *"). // 每2小时执行
    Description("Custom scheduled task").
    Do(func() error {
        // 自定义逻辑
        return nil
    })
```

## Cron 表达式格式

格式：`分钟 小时 日期 月份 星期`

```
字段         允许值          特殊字符
分钟         0-59           * , - /
小时         0-23           * , - /
日期         1-31           * , - /
月份         1-12           * , - /
星期         0-6 (0=周日)   * , - /
```

### 特殊字符说明

- `*` - 任意值
- `,` - 列表（例如：1,3,5）
- `-` - 范围（例如：1-5）
- `/` - 步长（例如：*/5 表示每5分钟）

### Cron 表达式示例

```bash
# 每分钟
* * * * *

# 每小时
0 * * * *

# 每天凌晨 2:00
0 2 * * *

# 每周一上午 9:00
0 9 * * 1

# 每月 1 号凌晨 3:00
0 3 1 * *

# 每 5 分钟
*/5 * * * *

# 每 2 小时
0 */2 * * *

# 工作日上午 9:00
0 9 * * 1-5

# 周末上午 10:00
0 10 * * 0,6
```

## 预定义调度方法

```go
// 时间间隔
EveryMinute()           // 每分钟
EveryFiveMinutes()      // 每 5 分钟
EveryTenMinutes()       // 每 10 分钟
EveryFifteenMinutes()   // 每 15 分钟
EveryThirtyMinutes()    // 每 30 分钟
Hourly()                // 每小时
HourlyAt(15)            // 每小时第 15 分钟

// 每天
Daily()                 // 每天午夜
DailyAt(8, 30)         // 每天 8:30

// 每周
Weekly()                // 每周日午夜
WeeklyOn(1, 9, 0)      // 每周一 9:00
Weekdays()             // 工作日午夜
Weekends()             // 周末午夜

// 每月
Monthly()               // 每月 1 号午夜
MonthlyOn(15, 14, 0)   // 每月 15 号 14:00

// 每年
Yearly()                // 每年 1 月 1 日午夜
```

## 命令行工具

### 启动调度器工作进程

```bash
go run cmd/artisan/main.go schedule:work
```

### 列出所有任务

```bash
go run cmd/artisan/main.go schedule:list
```

### 运行所有到期任务（一次性）

```bash
go run cmd/artisan/main.go schedule:run
```

## 任务管理

### 列出所有任务

```go
tasks := scheduler.ListTasks()
for _, task := range tasks {
    fmt.Printf("Task: %s, Next run: %s\n", task.Name, task.NextRunAt)
}
```

### 获取单个任务

```go
task, err := scheduler.GetTask("task-id")
if err != nil {
    log.Fatal(err)
}
```

### 手动运行任务

```go
err := scheduler.RunNow("task-id")
if err != nil {
    log.Fatal(err)
}
```

### 移除任务

```go
err := scheduler.RemoveTask("task-id")
if err != nil {
    log.Fatal(err)
}
```

## 任务日志

### 获取任务执行日志

```go
// 获取特定任务的最近 10 条日志
logs := scheduler.GetLogs("task-id", 10)
for _, log := range logs {
    fmt.Printf("Task: %s, Duration: %v, Success: %v\n",
        log.TaskName, log.Duration, log.Success)
    if !log.Success {
        fmt.Printf("Error: %s\n", log.Error)
    }
}

// 获取所有任务的最近 50 条日志
allLogs := scheduler.GetLogs("", 50)
```

## 统计信息

```go
stats := scheduler.GetStats()
fmt.Printf("Total tasks: %d\n", stats["total_tasks"])
fmt.Printf("Running tasks: %d\n", stats["running_tasks"])
fmt.Printf("Total runs: %d\n", stats["total_runs"])
fmt.Printf("Success rate: %.2f%%\n", stats["success_rate"])
```

## 实际应用示例

### 数据库备份

```go
scheduler.NewTask("backup-database").
    DailyAt(2, 0).
    Description("Daily database backup at 2 AM").
    Do(func() error {
        cmd := exec.Command("mysqldump", "-u", "root", "mydb")
        output, err := cmd.Output()
        if err != nil {
            return err
        }
        return os.WriteFile("/backup/db.sql", output, 0644)
    })
```

### 缓存清理

```go
scheduler.NewTask("cleanup-cache").
    EveryFiveMinutes().
    Description("Clean up expired cache entries").
    Do(func() error {
        cache := redis.GetCache()
        return cache.CleanupExpired()
    })
```

### 发送每日报告

```go
scheduler.NewTask("daily-report").
    DailyAt(9, 0).
    Description("Send daily report to admins").
    Do(func() error {
        report := generateDailyReport()
        return mail.Send("admin@example.com", "Daily Report", report)
    })
```

### 清理旧日志

```go
scheduler.NewTask("cleanup-logs").
    Weekly().
    Description("Clean up logs older than 30 days").
    Do(func() error {
        cutoff := time.Now().AddDate(0, 0, -30)
        return db.Where("created_at < ?", cutoff).Delete(&Log{}).Error
    })
```

## 最佳实践

### 1. 错误处理

```go
scheduler.NewTask("important-task").
    Hourly().
    Do(func() error {
        // 总是返回 error 以便记录失败
        if err := doSomething(); err != nil {
            // 记录详细错误
            log.Printf("Task failed: %v", err)
            return err
        }
        return nil
    })
```

### 2. 超时控制

```go
scheduler.NewTask("long-running-task").
    Daily().
    Do(func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
        defer cancel()
        
        return performLongOperation(ctx)
    })
```

### 3. 幂等性

确保任务可以安全地重复执行：

```go
scheduler.NewTask("send-notifications").
    EveryFiveMinutes().
    Do(func() error {
        // 只处理未发送的通知
        notifications := getUnsentNotifications()
        for _, notif := range notifications {
            if err := send(notif); err != nil {
                continue // 失败的保留下次处理
            }
            markAsSent(notif)
        }
        return nil
    })
```

### 4. 避免并发问题

如果任务需要独占资源，使用分布式锁：

```go
scheduler.NewTask("exclusive-task").
    Hourly().
    Do(func() error {
        lock := redis.Lock("task:exclusive", 5*time.Minute)
        if !lock.Acquire() {
            return fmt.Errorf("task already running")
        }
        defer lock.Release()
        
        // 执行独占任务
        return doExclusiveWork()
    })
```

## 监控和告警

### 监控任务执行

```go
// 定期检查任务状态
go func() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        stats := scheduler.GetStats()
        
        // 检查失败率
        if successRate := stats["success_rate"].(float64); successRate < 90 {
            sendAlert(fmt.Sprintf("Task success rate low: %.2f%%", successRate))
        }
        
        // 检查任务是否卡住
        tasks := scheduler.ListTasks()
        for _, task := range tasks {
            if task.IsRunning && time.Since(task.LastRunAt) > 30*time.Minute {
                sendAlert(fmt.Sprintf("Task %s stuck", task.Name))
            }
        }
    }
}()
```

## 注意事项

1. **时区**：所有时间使用服务器本地时区
2. **精度**：调度精度为分钟级，不支持秒级调度
3. **并发**：同一任务不会并发执行，如果上次还在运行则跳过
4. **持久化**：任务配置不持久化，重启后需要重新注册
5. **分布式**：当前实现不支持分布式调度，多实例会重复执行

## 下一步

- 查看 [队列系统](queue.md) 了解异步任务处理
- 查看 [事件系统](event.md) 了解事件驱动架构
