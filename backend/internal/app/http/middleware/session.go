package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// Session 配置
type SessionConfig struct {
	Driver          string        // 存储驱动: memory, redis, database
	CookieName      string        // Cookie名称
	CookieDomain    string        // Cookie域名
	SecureCookie    bool          // 是否仅HTTPS
	HTTPOnlyCookie  bool          // 是否仅HTTP访问
	CookieLifetime  time.Duration // Cookie有效期
	SessionLifetime time.Duration // Session有效期
	EncryptionKey   string        // 加密密钥
}

// Session 存储接口
type SessionStore interface {
	Get(key string) (string, error)
	Set(key string, value string, expire time.Duration) error
	Delete(key string) error
	Exists(key string) bool
}

// Session 实例
type Session struct {
	ID      string
	Store   SessionStore
	Config  *SessionConfig
	Context context.Context
}

// 默认配置
var defaultConfig = &SessionConfig{
	Driver:          "memory",
	CookieName:      "aidecms_session",
	CookieDomain:    "",
	SecureCookie:    false,
	HTTPOnlyCookie:  true,
	CookieLifetime:  24 * time.Hour,
	SessionLifetime: 24 * time.Hour,
}

// 创建Session中间件
func SessionMiddleware(config *SessionConfig) app.HandlerFunc {
	if config == nil {
		config = defaultConfig
	}

	// 初始化存储驱动
	var store SessionStore
	switch config.Driver {
	case "memory":
		store = NewMemoryStore()
	case "redis":
		store = NewRedisStore()
	case "database":
		store = NewDatabaseStore()
	default:
		panic("unsupported session driver: " + config.Driver)
	}

	return func(ctx context.Context, c *app.RequestContext) {
		// 从Cookie获取Session ID
		sessionID := string(c.Cookie(config.CookieName))
		if sessionID == "" {
			// 创建新Session
			sessionID = generateSessionID()
			c.SetCookie(
				config.CookieName,
				sessionID,
				int(config.CookieLifetime.Seconds()),
				"/",
				config.CookieDomain,
				1, // Lax mode
				config.SecureCookie,
				config.HTTPOnlyCookie,
			)
		}

		// 创建Session实例
		session := &Session{
			ID:      string(sessionID),
			Store:   store,
			Config:  config,
			Context: ctx,
		}

		// 存入上下文
		c.Set("session", session)

		// 继续处理请求
		c.Next(ctx)
	}
}

// 生成Session ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// 从上下文中获取Session
func GetSession(c *app.RequestContext) (*Session, error) {
	session, ok := c.Get("session")
	if !ok {
		return nil, errors.New("session not found")
	}
	return session.(*Session), nil
}

// 设置Session值
func (s *Session) Set(key, value string) error {
	return s.Store.Set(s.ID+":"+key, value, s.Config.SessionLifetime)
}

// 获取Session值
func (s *Session) Get(key string) (string, error) {
	return s.Store.Get(s.ID + ":" + key)
}

// 删除Session值
func (s *Session) Delete(key string) error {
	return s.Store.Delete(s.ID + ":" + key)
}

// 检查Session值是否存在
func (s *Session) Exists(key string) bool {
	return s.Store.Exists(s.ID + ":" + key)
}

// 销毁整个Session
func (s *Session) Destroy() error {
	// 实现需要根据存储驱动具体处理
	return nil
}
