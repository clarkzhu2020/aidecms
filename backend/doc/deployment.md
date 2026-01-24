# 部署指南

AideCMS 支持多种部署方式，包括传统服务器部署和容器化部署。

## 环境要求

- Go 1.18+
- MySQL 5.7+ 或 PostgreSQL 10+
- Redis (可选，用于队列和缓存)

## 构建应用

1. 设置环境变量:
```bash
cp .env.example .env.production
```

2. 修改生产环境配置:
```
APP_ENV=production
APP_DEBUG=false
```

3. 构建应用:
```bash
go build -o aidecms cmd/aidecms/main.go
```

## 传统部署

### 使用 Supervisor

`/etc/supervisor/conf.d/aidecms.conf`:
```ini
[program:aidecms]
command=/path/to/aidecms
directory=/path/to/aidecms
user=www-data
autostart=true
autorestart=true
stderr_logfile=/var/log/aidecms.err.log
stdout_logfile=/var/log/aidecms.out.log
```

启动服务:
```bash
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start aidecms
```

### Nginx 配置

`/etc/nginx/sites-available/aidecms`:
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://127.0.0.1:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 容器化部署

### Dockerfile

```dockerfile
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o aidecms cmd/aidecms/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/clarkgo .
COPY --from=builder /app/.env.production .env

EXPOSE 8888
CMD ["./aidecms"]
```

### docker-compose.yml

```yaml
version: '3'

services:
  app:
    build: .
    ports:
      - "8888:8888"
    depends_on:
      - db
      - redis
    environment:
      - APP_ENV=production
      - DB_HOST=db
      - REDIS_HOST=redis

  db:
    image: mysql:5.7
    environment:
      - MYSQL_ROOT_PASSWORD=secret
      - MYSQL_DATABASE=clarkgo
    volumes:
      - db_data:/var/lib/mysql

  redis:
    image: redis:alpine

volumes:
  db_data:
```

## 部署 Artisan 命令

### 数据库迁移

```bash
docker-compose exec app ./aidecms migrate
```

### 队列处理

```bash
docker-compose exec app ./aidecms queue:work
```

## 性能优化

1. 启用 Gzip 压缩:
```go
h.Use(hertzMiddleware.Gzip())
```

2. 使用连接池:
```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    ConnPool: &sql.DB{
        MaxIdleConns: 10,
        MaxOpenConns: 100,
    },
})
```

3. 启用 HTTP/2:
```go
h := server.Default(
    server.WithHostPorts(":8888"),
    server.WithH2C(true),
)
```

## 监控

集成 Prometheus 监控:

1. 添加中间件:
```go
import "github.com/hertz-contrib/monitor-prometheus"

h.Use(monitor_prometheus.NewMonitor())
```

2. 访问指标:
```
http://localhost:8888/metrics