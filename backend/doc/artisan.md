# AideCMS Artisan 命令行工具

## 基本使用

```bash
go run . artisan [command] [options] [arguments]
```

## 命令参考

### 数据库命令
```bash
# 运行迁移
go run . artisan migrate

# 创建新迁移文件
go run . artisan make:migration create_users_table

# 清空并重新运行所有迁移
go run . artisan migrate:fresh
```

### 缓存命令
```bash
# 清空应用缓存
go run . artisan cache:clear
```

### 队列命令
```bash
# 启动队列工作进程
go run . artisan queue:work

# 重试失败任务
go run . artisan queue:retry
```

### 计划任务
```bash
# 运行计划任务
go run . artisan schedule:run

# 列出所有计划任务
go run . artisan schedule:list
```

### 事件系统
```bash
# 创建新事件
go run . artisan make:event OrderShipped

# 创建事件监听器
go run . artisan make:listener SendShipmentNotification