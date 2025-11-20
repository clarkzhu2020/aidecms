# AideCMS 官网部署指南

## 📋 部署前准备

### 1. 构建项目

```bash
npm install
npm run build
```

构建后的文件将在 `dist` 目录中。

### 2. 验证构建

```bash
npm run preview
```

访问 http://localhost:4173 预览构建结果。

## 🚀 部署方式

### 方式一：Nginx 部署

#### 1. 复制文件到服务器

```bash
scp -r dist/* user@your-server:/var/www/aidecms-web/
```

#### 2. 配置 Nginx

复制 `nginx.conf.example` 到 Nginx 配置目录：

```bash
sudo cp nginx.conf.example /etc/nginx/sites-available/aidecms-web
sudo ln -s /etc/nginx/sites-available/aidecms-web /etc/nginx/sites-enabled/
```

修改配置中的域名和路径：

```nginx
server_name your-domain.com;
root /var/www/aidecms-web/dist;
```

#### 3. 测试并重载 Nginx

```bash
sudo nginx -t
sudo systemctl reload nginx
```

### 方式二：Vercel 部署

#### 1. 安装 Vercel CLI

```bash
npm install -g vercel
```

#### 2. 登录并部署

```bash
vercel login
vercel --prod
```

### 方式三：Netlify 部署

#### 1. 在 Netlify 创建新站点

#### 2. 配置构建设置

- Build command: `npm run build`
- Publish directory: `dist`

#### 3. 推送代码自动部署

### 方式四：Docker 部署

#### 1. 创建 Dockerfile

```dockerfile
FROM node:18-alpine as builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf.example /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

#### 2. 构建并运行

```bash
docker build -t aidecms-web .
docker run -d -p 80:80 aidecms-web
```

## 🔧 环境变量配置

如需配置 API 地址等环境变量，创建 `.env.production` 文件：

```env
VITE_API_URL=https://api.aidecms.com
```

## 📊 性能优化

### 1. 启用 Gzip 压缩

Nginx 配置已包含 Gzip 设置。

### 2. CDN 加速

将静态资源上传到 CDN：

```bash
# 修改 vite.config.ts
export default defineConfig({
  base: 'https://cdn.your-domain.com/'
})
```

### 3. 缓存策略

静态资源设置长期缓存（1年）。

## 🔒 安全配置

### 1. HTTPS 配置

使用 Let's Encrypt 获取免费证书：

```bash
sudo certbot --nginx -d your-domain.com
```

### 2. 安全头部

已在 nginx.conf.example 中配置：
- X-Frame-Options
- X-Content-Type-Options
- X-XSS-Protection

## 📱 监控和维护

### 1. 日志查看

```bash
# Nginx 访问日志
sudo tail -f /var/log/nginx/access.log

# Nginx 错误日志
sudo tail -f /var/log/nginx/error.log
```

### 2. 性能监控

建议使用：
- Google Analytics
- Sentry (错误监控)
- Lighthouse (性能测试)

## 🆘 常见问题

### Q: 页面刷新后 404

A: 确保 Nginx 配置了 SPA 路由：

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

### Q: 静态资源 404

A: 检查 `base` 配置和资源路径。

### Q: 构建失败

A: 检查 Node.js 版本（需要 18+）和依赖安装。

## 📞 支持

- GitHub Issues: https://github.com/chenyusolar/aidecms/issues
- 文档: https://github.com/chenyusolar/aidecms/tree/main/doc
