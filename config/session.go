package config

import "time"

type Session struct {
	Driver          string        `mapstructure:"driver"`
	CookieName      string        `mapstructure:"cookie_name"`
	CookieDomain    string        `mapstructure:"cookie_domain"`
	SecureCookie    bool          `mapstructure:"secure_cookie"`
	HTTPOnlyCookie  bool          `mapstructure:"http_only_cookie"`
	CookieLifetime  time.Duration `mapstructure:"cookie_lifetime"`
	SessionLifetime time.Duration `mapstructure:"session_lifetime"`
	EncryptionKey   string        `mapstructure:"encryption_key"`
	Redis           struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
}

var SessionConfig = &Session{
	Driver:          "memory",
	CookieName:      "aidecms_session",
	CookieDomain:    "",
	SecureCookie:    false,
	HTTPOnlyCookie:  true,
	CookieLifetime:  24 * time.Hour,
	SessionLifetime: 24 * time.Hour,
	EncryptionKey:   "default-secret-key",
	Redis: struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	},
}
