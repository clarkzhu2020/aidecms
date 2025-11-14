# 邮件服务API文档

ClarkGo 框架现已集成完整的邮件服务API，支持SMTP、TLS/SSL加密、模板渲染、批量发送等功能。

## 配置说明

在 `.env` 文件中配置邮件服务参数：

```bash
# 邮件配置
MAIL_DRIVER=smtp
MAIL_HOST=mailhog                    # 邮件服务器地址
MAIL_PORT=1025                       # 邮件服务器端口
MAIL_USERNAME=your_email@example.com # 邮箱用户名
MAIL_PASSWORD=your_password          # 邮箱密码
MAIL_ENCRYPTION=tls                  # 加密方式：tls/ssl/none
MAIL_FROM_ADDRESS=noreply@example.com # 发件人邮箱
MAIL_FROM_NAME=ClarkGo               # 发件人姓名
```

## API 端点

### 1. 获取邮件配置 - GET /api/mail/config

获取当前邮件服务配置信息（敏感信息已脱敏）。

**请求示例：**
```bash
curl -X GET http://localhost:8888/api/mail/config
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "driver": "smtp",
    "host": "mailhog", 
    "port": 1025,
    "encryption": "tls",
    "from_email": "noreply@example.com",
    "from_name": "ClarkGo",
    "username": "***"
  }
}
```

### 2. 测试邮件连接 - GET /api/mail/test

测试邮件服务器连接状态，不实际发送邮件。

**请求示例：**
```bash
curl -X GET http://localhost:8888/api/mail/test
```

**响应示例：**
```json
{
  "success": true,
  "message": "Mail server connection is working"
}
```

### 3. 验证邮箱地址 - GET /api/mail/validate

验证邮箱地址格式是否正确。

**请求参数：**
- `email` (query): 要验证的邮箱地址

**请求示例：**
```bash
curl -X GET "http://localhost:8888/api/mail/validate?email=test@example.com"
```

**响应示例：**
```json
{
  "success": true,
  "email": "test@example.com",
  "message": "Email address is valid"
}
```

### 4. 发送普通邮件 - POST /api/mail/send

发送单封邮件，支持纯文本和HTML格式。

**请求体参数：**
```json
{
  "to": ["recipient@example.com"],           // 收件人列表
  "cc": ["cc@example.com"],                  // 抄送列表（可选）
  "bcc": ["bcc@example.com"],                // 密送列表（可选）
  "subject": "邮件主题",                      // 邮件主题
  "body": "邮件内容",                        // 邮件正文
  "isHTML": false                            // 是否为HTML格式
}
```

**请求示例：**
```bash
curl -X POST http://localhost:8888/api/mail/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "测试邮件",
    "body": "这是一封测试邮件",
    "isHTML": false
  }'
```

### 5. 使用模板发送邮件 - POST /api/mail/send-template

使用HTML模板发送邮件，支持模板变量替换。

**请求体参数：**
```json
{
  "to": ["recipient@example.com"],
  "subject": "邮件主题", 
  "template": "<h1>Hello {{.Name}}</h1><p>Welcome to {{.App}}</p>",
  "data": {
    "Name": "张三",
    "App": "ClarkGo"
  }
}
```

**请求示例：**
```bash
curl -X POST http://localhost:8888/api/mail/send-template \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "欢迎使用ClarkGo",
    "template": "<h1>Hello {{.Name}}</h1><p>Welcome to {{.App}}</p>",
    "data": {
      "Name": "用户",
      "App": "ClarkGo框架"
    }
  }'
```

### 6. 批量发送邮件 - POST /api/mail/send-bulk

批量发送邮件，每个收件人可以有不同的内容。

**请求体参数：**
```json
{
  "emails": [
    {
      "to": ["user1@example.com"],
      "subject": "个人化主题1",
      "body": "个人化内容1",
      "isHTML": false
    },
    {
      "to": ["user2@example.com"], 
      "subject": "个人化主题2",
      "body": "个人化内容2",
      "isHTML": false
    }
  ]
}
```

## 错误处理

所有API在出错时返回标准错误格式：

```json
{
  "success": false,
  "error": "错误类型",
  "message": "详细错误信息"
}
```

常见错误代码：
- `400` - 请求参数错误
- `500` - 服务器内部错误（配置错误、连接失败等）

## 附件支持

邮件服务支持Base64编码的附件。在发送邮件请求中添加 `attachments` 字段：

```json
{
  "to": ["user@example.com"],
  "subject": "带附件的邮件",
  "body": "请查收附件",
  "attachments": [
    {
      "filename": "document.pdf",
      "content": "base64编码的文件内容",
      "contentType": "application/pdf"
    }
  ]
}
```

## 安全注意事项

1. 确保 `.env` 文件不被提交到版本控制系统
2. 使用强密码和安全的SMTP服务
3. 在生产环境中启用TLS/SSL加密
4. 对邮件发送频率进行限制以防止滥用

## 测试

使用提供的测试脚本可以快速验证邮件API功能：

```bash
python3 test_mail_api.py
```

测试脚本会依次测试所有API端点并显示结果。