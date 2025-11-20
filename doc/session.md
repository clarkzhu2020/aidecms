# Session 功能

AideCMS 提供了完整的 Session 功能，支持多种存储驱动。

## 配置

在 `.env` 中添加配置:

```
SESSION_DRIVER=memory  # memory, redis, database
SESSION_NAME=aidecms_session
SESSION_LIFETIME=1440  # 分钟
SESSION_SECURE_COOKIE=false
SESSION_HTTP_ONLY=true
SESSION_DOMAIN=
SESSION_ENCRYPTION_KEY=your-secret-key

# Redis 配置 (当使用redis驱动时)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## 基本使用

### 启用 Session 中间件

在路由中启用:

```go
import "internal/app/http/middleware"

func RegisterRoutes(h *server.Hertz) {
    // 使用默认配置
    h.Use(middleware.SessionMiddleware(nil))
    
    // 或自定义配置
    config := &middleware.SessionConfig{
        Driver: "redis",
        CookieName: "my_session",
    }
    h.Use(middleware.SessionMiddleware(config))
}
```

### 控制器中使用

```go
func (c *UserController) Login(ctx context.Context, hCtx *app.RequestContext) {
    session, err := middleware.GetSession(hCtx)
    if err != nil {
        // 处理错误
    }
    
    // 设置Session值
    session.Set("user_id", "123")
    session.Set("username", "john")
    
    // 设置过期时间(覆盖全局配置)
    session.Store.Set(session.ID+":user_id", "123", 30*time.Minute)
}

func (c *UserController) Profile(ctx context.Context, hCtx *app.RequestContext) {
    session, _ := middleware.GetSession(hCtx)
    
    // 获取Session值
    userID, _ := session.Get("user_id")
    
    // 检查是否存在
    if session.Exists("username") {
        username, _ := session.Get("username")
    }
    
    // 删除Session值
    session.Delete("temp_data")
}
```

## 存储驱动

### 内存驱动

- 默认驱动
- 适合开发和测试环境
- 重启应用会丢失所有Session

### Redis 驱动

- 高性能分布式存储
- 适合生产环境
- 需要配置Redis连接

### 数据库驱动

- 使用数据库存储
- 适合需要持久化的场景
- 需要创建sessions表

## 最佳实践

1. 生产环境使用Redis驱动
2. 设置足够复杂的加密密钥
3. 敏感数据设置较短的过期时间
4. 用户登出时销毁Session
5. 避免在Session中存储大量数据

## 安全建议

1. 启用 SecureCookie (HTTPS环境)
2. 启用 HttpOnlyCookie
3. 定期更换加密密钥
4. 限制Session生命周期
5. 实现Session劫持检测