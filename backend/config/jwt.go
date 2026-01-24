package config

import (
	"os"
	"strconv"
)

// JWT holds the JWT configuration

type JWTConfig struct {
	SecretKey string
	ExpiresIn int // 过期时间(小时)
}

var JWT = &JWTConfig{
	SecretKey: "default_secret_key",
	ExpiresIn: 24,
}

func InitJWT() {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		JWT.SecretKey = secret
	}
	if expires := os.Getenv("JWT_EXPIRES"); expires != "" {
		if exp, err := strconv.Atoi(expires); err == nil {
			JWT.ExpiresIn = exp
		}
	}
}
