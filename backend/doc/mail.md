# 邮件系统

AideCMS 集成了强大的邮件发送功能，支持模板和队列。

## 配置

邮件配置在 `.env` 文件中:

```
MAIL_DRIVER=smtp
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=user@example.com
MAIL_PASSWORD=secret
MAIL_ENCRYPTION=tls
MAIL_FROM_ADDRESS=no-reply@example.com
MAIL_FROM_NAME="AideCMS"
```

## 基本使用

### 快速发送

```go
import "internal/app/services/mail"

err := mail.To("user@example.com").
    Subject("Welcome").
    Body("Hello, welcome to our service!").
    Send()
```

### 使用模板

1. 创建模板 `resources/views/emails/welcome.tpl`:

```html
<h1>Welcome, {{.name}}!</h1>
<p>Your account has been created.</p>
```

2. 发送模板邮件:

```go
err := mail.To("user@example.com").
    Template("emails.welcome", map[string]string{
        "name": "John",
    }).
    Send()
```

## 队列邮件

推荐使用队列发送邮件以避免阻塞:

```go
jobs.NewEmailJob("emails.welcome", "user@example.com", map[string]string{
    "name": "John",
}).Dispatch()
```

## 附件

添加附件:

```go
err := mail.To("user@example.com").
    Subject("Your documents").
    Attach("/path/to/file.pdf").
    Send()
```

## 测试

开发环境可以使用日志驱动(不实际发送):

```
MAIL_DRIVER=log
```

发送的邮件会记录在 `storage/logs/mail.log`。