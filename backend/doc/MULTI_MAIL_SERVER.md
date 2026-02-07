# 多邮件服务器支持文档

## 概述

AideCMS现在支持多邮件服务器配置，提供以下核心功能：

1. **自动故障切换** - 当某个邮件服务器不可用时，自动切换到其他服务器
2. **负载均衡** - 支持轮询、优先级、权重等多种负载均衡策略
3. **健康检查** - 实时监控服务器状态，自动恢复健康的服务器
4. **智能封禁检测** - 自动识别限流和封禁情况，自动切换服务器
5. **批量发送优化** - 群发邮件时自动分配到多个服务器，降低封禁风险

## 目录

- [配置说明](#配置说明)
- [负载均衡模式](#负载均衡模式)
- [API端点](#api端点)
- [使用示例](#使用示例)
- [故障切换机制](#故障切换机制)
- [健康检查](#健康检查)
- [最佳实践](#最佳实践)

## 配置说明

### 方式1: 单服务器配置（兼容旧版本）

```bash
MAIL_MAILER=smtp
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your@gmail.com
MAIL_PASSWORD=your-password
MAIL_ENCRYPTION=tls
MAIL_FROM_ADDRESS=your@gmail.com
MAIL_FROM_NAME="Your Company"
```

### 方式2: 多服务器配置（推荐）

在 `.env` 文件中配置多个邮件服务器：

```bash
# 基础配置
MAIL_DRIVER=smtp
MAIL_FAILOVER_ENABLED=true
MAIL_LOAD_BALANCE_MODE=round_robin  # round_robin, priority, weighted
MAIL_HEALTH_CHECK_INTERVAL=60

# 多服务器配置（JSON格式，注意是单行）
MAIL_SERVERS=[
  {
    "id": "gmail_primary",
    "host": "smtp.gmail.com",
    "port": 587,
    "username": "user1@gmail.com",
    "password": "password1",
    "encryption": "tls",
    "from_name": "Your Company",
    "from_email": "user1@gmail.com",
    "enabled": true,
    "priority": 1,
    "weight": 10
  },
  {
    "id": "outlook_secondary",
    "host": "smtp-mail.outlook.com",
    "port": 587,
    "username": "user2@outlook.com",
    "password": "password2",
    "encryption": "tls",
    "from_name": "Your Company",
    "from_email": "user2@outlook.com",
    "enabled": true,
    "priority": 2,
    "weight": 8
  },
  {
    "id": "custom_smtp",
    "host": "smtp.example.com",
    "port": 587,
    "username": "user3@example.com",
    "password": "password3",
    "encryption": "tls",
    "from_name": "Your Company",
    "from_email": "user3@example.com",
    "enabled": true,
    "priority": 3,
    "weight": 5
  }
]
```

### 配置参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| `id` | 服务器唯一标识 | `gmail_primary` |
| `host` | SMTP服务器地址 | `smtp.gmail.com` |
| `port` | SMTP端口 | `587` |
| `username` | SMTP用户名 | `user@gmail.com` |
| `password` | SMTP密码 | `your-password` |
| `encryption` | 加密方式 | `tls`, `ssl`, `none` |
| `from_name` | 发件人姓名 | `Your Company` |
| `from_email` | 发件人邮箱 | `user@gmail.com` |
| `enabled` | 是否启用 | `true`, `false` |
| `priority` | 优先级（数字越小优先级越高） | `1`, `2`, `3` |
| `weight` | 权重（用于加权负载均衡） | `10`, `8`, `5` |

### 全局配置参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `MAIL_FAILOVER_ENABLED` | 是否启用故障切换 | `true` |
| `MAIL_LOAD_BALANCE_MODE` | 负载均衡模式 | `round_robin` |
| `MAIL_HEALTH_CHECK_INTERVAL` | 健康检查间隔（秒） | `60` |

## 负载均衡模式

### 1. 轮询模式 (round_robin)

按顺序依次使用每个服务器，适用于服务器性能相近的场景。

**特点：**
- 简单公平
- 所有服务器平均分担负载
- 不会出现某些服务器过载的情况

**使用场景：**
- 所有服务器性能相同
- 需要平均分配发送量
- 群发邮件时的默认选择

### 2. 优先级模式 (priority)

优先使用优先级高的服务器，高优先级服务器不可用时才使用低优先级服务器。

**特点：**
- 优先使用主服务器
- 主服务器故障时自动切换到备用服务器
- 适合有主备关系的场景

**使用场景：**
- 有主备服务器部署
- 成本敏感（优先使用免费/低成本服务器）
- 需要保证主服务器的稳定性

### 3. 权重模式 (weighted)

根据服务器权重分配发送量，权重越高分配越多。

**特点：**
- 根据服务器性能分配负载
- 性能强的服务器处理更多邮件
- 灵活的负载分配

**使用场景：**
- 服务器性能差异较大
- 某些服务器有更高的发送配额
- 需要优化资源利用率

## API端点

### 1. 发送邮件 - POST /api/multi-mail/send

发送单封邮件，自动选择可用的服务器。

**请求体：**
```json
{
  "to": ["recipient@example.com"],
  "cc": ["cc@example.com"],
  "bcc": ["bcc@example.com"],
  "subject": "邮件主题",
  "body": "邮件内容",
  "html_body": "<h1>HTML内容</h1>",
  "attachments": [
    {
      "filename": "file.pdf",
      "content": "base64编码内容",
      "mime_type": "application/pdf"
    }
  ]
}
```

**响应示例：**
```json
{
  "success": true,
  "message": "Email sent successfully",
  "data": {
    "to": ["recipient@example.com"],
    "subject": "邮件主题"
  }
}
```

### 2. 批量发送邮件 - POST /api/multi-mail/send-bulk

批量发送邮件，自动负载均衡到多个服务器。

**请求体：**
```json
{
  "emails": [
    {
      "to": ["user1@example.com"],
      "subject": "主题1",
      "body": "内容1"
    },
    {
      "to": ["user2@example.com"],
      "subject": "主题2",
      "body": "内容2"
    }
  ]
}
```

**响应示例：**
```json
{
  "success": true,
  "message": "Bulk email sending completed. Success: 198, Failed: 2",
  "data": {
    "total": 200,
    "success_count": 198,
    "failure_count": 2,
    "server_stats": {
      "gmail_primary": 80,
      "outlook_secondary": 78,
      "custom_smtp": 40
    },
    "results": [
      {
        "index": 0,
        "success": true,
        "server_id": "gmail_primary",
        "to": ["user1@example.com"]
      }
    ]
  }
}
```

### 3. 获取服务器状态 - GET /api/multi-mail/servers

获取所有邮件服务器的当前状态。

**响应示例：**
```json
{
  "success": true,
  "message": "Server status retrieved successfully",
  "data": {
    "gmail_primary": {
      "host": "smtp.gmail.com",
      "port": 587,
      "is_healthy": true,
      "is_banned": false,
      "rate_limit_hit": false,
      "fail_count": 0,
      "total_sent": 1523,
      "total_failed": 5,
      "last_error": "",
      "enabled": true,
      "priority": 1,
      "weight": 10
    },
    "outlook_secondary": {
      "host": "smtp-mail.outlook.com",
      "port": 587,
      "is_healthy": true,
      "is_banned": false,
      "rate_limit_hit": false,
      "fail_count": 2,
      "total_sent": 890,
      "total_failed": 12,
      "last_error": "timeout",
      "enabled": true,
      "priority": 2,
      "weight": 8
    }
  }
}
```

### 4. 检查服务器健康 - GET /api/multi-mail/server/health

检查指定服务器的健康状态。

**参数：**
- `server_id` (必填) - 服务器ID

**请求示例：**
```bash
curl http://localhost:8888/api/multi-mail/server/health?server_id=gmail_primary
```

### 5. 封禁服务器 - POST /api/multi-mail/server/ban

手动封禁指定服务器。

**参数：**
- `server_id` (必填) - 服务器ID
- `reason` (可选) - 封禁原因

**请求示例：**
```bash
curl -X POST "http://localhost:8888/api/multi-mail/server/ban?server_id=gmail_primary&reason=manual_ban"
```

### 6. 解封服务器 - POST /api/multi-mail/server/unban

解封指定服务器。

**参数：**
- `server_id` (必填) - 服务器ID

**请求示例：**
```bash
curl -X POST "http://localhost:8888/api/multi-mail/server/unban?server_id=gmail_primary"
```

### 7. 获取邮件配置 - GET /api/multi-mail/config

获取当前邮件配置（敏感信息已隐藏）。

**响应示例：**
```json
{
  "success": true,
  "message": "Mail config retrieved successfully",
  "data": {
    "driver": "smtp",
    "failover_enabled": true,
    "load_balance_mode": "round_robin",
    "servers": [
      {
        "id": "gmail_primary",
        "host": "smtp.gmail.com",
        "port": 587,
        "encryption": "tls",
        "from_name": "Your Company",
        "from_email": "user@gmail.com",
        "enabled": true,
        "priority": 1,
        "weight": 10,
        "username": "us***@gmail.com"
      }
    ]
  }
}
```

### 8. 测试所有服务器 - GET /api/multi-mail/test

测试所有服务器的连接状态。

**响应示例：**
```json
{
  "success": true,
  "message": "Server connection test completed",
  "data": {
    "total_servers": 3,
    "results": [
      {
        "server_id": "gmail_primary",
        "host": "smtp.gmail.com",
        "port": 587,
        "status": {
          "is_healthy": true,
          "enabled": true
        }
      }
    ]
  }
}
```

## 使用示例

### 示例1: 基础发送

```bash
curl -X POST http://localhost:8888/api/multi-mail/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "测试邮件",
    "body": "这是一封测试邮件"
  }'
```

### 示例2: 批量发送（群发）

```bash
curl -X POST http://localhost:8888/api/multi-mail/send-bulk \
  -H "Content-Type: application/json" \
  -d '{
    "emails": [
      {
        "to": ["user1@example.com"],
        "subject": "活动通知",
        "body": "您有新的活动通知"
      },
      {
        "to": ["user2@example.com"],
        "subject": "活动通知",
        "body": "您有新的活动通知"
      },
      {
        "to": ["user3@example.com"],
        "subject": "活动通知",
        "body": "您有新的活动通知"
      }
    ]
  }'
```

### 示例3: 查看服务器状态

```bash
curl http://localhost:8888/api/multi-mail/servers
```

### 示例4: 手动封禁问题服务器

```bash
curl -X POST "http://localhost:8888/api/multi-mail/server/ban?server_id=gmail_primary&reason=sending_too_many_emails"
```

## 故障切换机制

### 自动故障切换流程

1. **发送失败检测**
   - 当发送邮件失败时，系统自动记录错误
   - 连续失败5次后标记服务器为不健康

2. **自动切换**
   - 如果启用了 `MAIL_FAILOVER_ENABLED=true`
   - 系统自动尝试下一个可用服务器
   - 按照负载均衡模式选择下一个服务器

3. **服务器恢复**
   - 健康检查器定期检测服务器状态
   - 服务器恢复健康后自动启用
   - 被解封的服务器重新加入可用队列

### 智能封禁检测

系统会自动检测以下情况并采取相应措施：

| 检测项 | 触发条件 | 处理方式 |
|---------|---------|---------|
| 限流 | 错误包含 "rate limit", "too many" | 标记为限流状态，暂停使用 |
| 封禁 | 错误包含 "blocked", "banned", "suspended" | 标记为封禁状态，自动切换 |
| 连续失败 | 连续失败5次以上 | 标记为不健康，暂停使用 |

### 自动恢复机制

- 限流服务器在一段时间后自动尝试恢复
- 健康检查通过的服务器自动解封
- 手动解封的服务器立即恢复使用

## 健康检查

### 健康检查机制

健康检查器会定期检查所有服务器的连接状态：

1. **连接测试** - 尝试建立SMTP连接
2. **认证验证** - 验证认证是否成功
3. **状态更新** - 更新服务器健康状态
4. **自动恢复** - 恢复健康的服务器

### 健康检查配置

```bash
MAIL_HEALTH_CHECK_INTERVAL=60  # 每60秒检查一次
```

### 健康检查API

```bash
# 检查单个服务器
GET /api/multi-mail/server/health?server_id=gmail_primary

# 检查所有服务器
GET /api/multi-mail/servers

# 测试连接
GET /api/multi-mail/test
```

## 最佳实践

### 1. 服务器选择

**推荐组合：**
- **Gmail** - 可靠性高，但有日发送限制
- **Outlook** - 免费版本限制较少
- **自定义SMTP** - 企业邮箱，适合大量发送
- **第三方邮件服务** - SendGrid, Mailgun, SES等

### 2. 负载均衡配置

**场景1: 均等分配**
```bash
MAIL_LOAD_BALANCE_MODE=round_robin
# 所有服务器权重相同
```

**场景2: 主备模式**
```bash
MAIL_LOAD_BALANCE_MODE=priority
# 主服务器 priority=1
# 备用服务器 priority=2,3...
```

**场景3: 性能加权**
```bash
MAIL_LOAD_BALANCE_MODE=weighted
# 高性能服务器 weight=10
# 中等性能服务器 weight=5
# 低性能服务器 weight=2
```

### 3. 发送策略

**小批量发送（<100封）**
- 使用单个高优先级服务器
- 快速完成发送

**中等批量（100-1000封）**
- 使用轮询模式
- 平均分配到所有服务器

**大批量（>1000封）**
- 使用权重模式
- 根据服务器性能分配
- 监控服务器状态，及时调整

### 4. 防封禁建议

1. **控制发送速度**
   - 每分钟不超过50封
   - 使用批量接口控制发送节奏

2. **多样化内容**
   - 避免完全相同的邮件内容
   - 添加个性化信息

3. **监控服务器状态**
   - 定期检查服务器状态
   - 及时处理封禁服务器

4. **使用发送队列**
   - 不要一次性发送大量邮件
   - 分批发送，间隔一段时间

### 5. 故障处理

1. **配置足够的备用服务器**
   - 至少配置2-3个备用服务器
   - 确保至少2个服务器始终可用

2. **启用故障切换**
   ```bash
   MAIL_FAILOVER_ENABLED=true
   ```

3. **设置合理的健康检查间隔**
   ```bash
   MAIL_HEALTH_CHECK_INTERVAL=60  # 1分钟
   ```

4. **监控和告警**
   - 监控服务器失败率
   - 设置告警阈值（如失败率>10%）
   - 及时处理问题服务器

### 6. 测试和验证

1. **测试发送**
   ```bash
   # 测试单个服务器
   curl http://localhost:8888/api/multi-mail/test
   
   # 测试所有服务器
   curl http://localhost:8888/api/multi-mail/servers
   ```

2. **小批量测试**
   - 先发送10-20封测试邮件
   - 检查服务器状态和发送结果
   - 确认所有服务器正常工作

3. **监控发送结果**
   - 查看批量发送的详细结果
   - 检查每个服务器的分配情况
   - 分析失败原因

## 故障排除

### 问题1: 所有服务器都不可用

**原因：**
- 所有服务器都被封禁
- 配置错误
- 网络问题

**解决方案：**
1. 检查服务器状态
2. 手动解封服务器
3. 验证配置是否正确
4. 检查网络连接

### 问题2: 邮件发送频繁失败

**原因：**
- 限流
- 服务器配置错误
- 认证失败

**解决方案：**
1. 降低发送频率
2. 检查服务器凭据
3. 查看服务器状态
4. 启用更多备用服务器

### 问题3: 某个服务器一直不被使用

**原因：**
- 优先级太低
- 权重为0
- 被封禁或禁用

**解决方案：**
1. 调整优先级
2. 增加权重
3. 检查服务器状态
4. 解封服务器

## 迁移指南

### 从单服务器迁移到多服务器

1. **备份现有配置**
   ```bash
   cp .env .env.backup
   ```

2. **更新配置**
   - 将现有单服务器配置添加到 `MAIL_SERVERS`
   - 设置 `MAIL_FAILOVER_ENABLED=true`
   - 选择合适的负载均衡模式

3. **测试新配置**
   ```bash
   # 测试服务器连接
   curl http://localhost:8888/api/multi-mail/test
   
   # 发送测试邮件
   curl -X POST http://localhost:8888/api/multi-mail/send \
     -H "Content-Type: application/json" \
     -d '{"to":["test@example.com"],"subject":"Test","body":"Test"}'
   ```

4. **更新API调用**
   - 将 `/api/mail/*` 改为 `/api/multi-mail/*`
   - 享受多服务器带来的优势

## 性能优化

### 1. 连接池

系统内部已实现连接池，避免频繁建立连接。

### 2. 异步发送

对于大批量发送，建议使用异步队列：

```go
// 添加到发送队列
jobs.NewBulkEmailJob(emails).Dispatch()
```

### 3. 批量发送优化

- 使用批量发送接口而非循环单发
- 合理设置批次大小（建议每批50-100封）
- 批次之间添加间隔

## 安全建议

1. **保护SMTP凭据**
   - 不要将 `.env` 文件提交到版本控制
   - 使用环境变量或密钥管理服务

2. **使用应用专用密码**
   - 不要使用邮箱登录密码
   - 使用SMTP专用密码
   - 定期更换密码

3. **限制API访问**
   - 使用API密钥认证
   - 限制IP访问
   - 实施请求频率限制

4. **监控异常**
   - 监控异常发送行为
   - 设置告警机制
   - 及时响应安全事件

## 支持和反馈

如有问题或建议，请：
- 查看日志文件
- 检查服务器状态API
- 联系技术支持
