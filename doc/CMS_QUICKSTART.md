# ClarkGo CMS 快速开始指南

## 📦 1. 安装和初始化

### 安装依赖
```bash
cd /path/to/clarkgo
go mod tidy
```

### 编译项目
```bash
go build -o clarkgo main.go
```

### 初始化CMS数据库
```bash
# 这会创建所有CMS表并初始化默认角色和权限
go run cmd/artisan/main.go cms:init
```

输出示例：
```
Initializing CMS...
Creating CMS tables...
✓ CMS tables created successfully

Creating default roles and permissions...
✓ Default roles and permissions created successfully

✓ CMS initialization completed successfully!

Default roles created:
  - super_admin: Full system access
  - admin: Administrative access
  - editor: Content management
  - author: Create and manage own posts
  - user: Basic user access
```

---

## 🚀 2. 启动服务

```bash
./clarkgo
```

服务将在 `http://localhost:8888` 启动

---

## 👤 3. 创建管理员用户

### 注册用户
```bash
curl -X POST http://localhost:8888/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123456",
    "first_name": "Admin",
    "last_name": "User"
  }'
```

### 为用户分配管理员角色

使用Go代码或创建一个命令工具：

```go
// 示例代码（在临时脚本中运行）
permService := services.NewPermissionService()
adminRole, _ := permService.GetRoleByName("admin")
permService.AssignRoleToUser(1, adminRole.ID) // 1是用户ID
```

或者创建Artisan命令：
```bash
go run cmd/artisan/main.go user:assign-role admin admin@example.com
```

---

## 🔑 4. 登录获取Token

```bash
curl -X POST http://localhost:8888/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123456"
  }'
```

响应：
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com"
    }
  }
}
```

**保存这个token，后续请求都需要它！**

---

## 📁 5. 上传媒体文件

```bash
curl -X POST http://localhost:8888/api/cms/media/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "files=@/path/to/image.jpg" \
  -F "files=@/path/to/document.pdf"
```

响应：
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "file_name": "uuid-generated-name.jpg",
      "original_name": "image.jpg",
      "file_url": "/uploads/2024/01/02/uuid.jpg",
      "file_size": 102400,
      "file_type": "image",
      "thumbnails": "{\"small\":\"...\",\"medium\":\"...\",\"large\":\"...\"}",
      "width": 1920,
      "height": 1080
    }
  ]
}
```

---

## 🏷️ 6. 创建分类

```bash
curl -X POST http://localhost:8888/api/cms/categories \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "技术文章",
    "description": "关于技术的文章",
    "meta_title": "技术文章 - 我的博客",
    "meta_description": "最新的技术文章和教程"
  }'
```

---

## 🔖 7. 创建标签

```bash
curl -X POST http://localhost:8888/api/cms/tags \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go语言"
  }'
```

---

## 📝 8. 创建文章

```bash
curl -X POST http://localhost:8888/api/cms/posts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "ClarkGo框架入门指南",
    "content": "# ClarkGo是什么\n\nClarkGo是一个基于Hertz的高性能Go Web框架...",
    "excerpt": "快速了解ClarkGo框架的核心特性",
    "featured_image": "/uploads/2024/01/02/image.jpg",
    "status": "published",
    "category_id": 1,
    "tags": [1, 2],
    "meta_title": "ClarkGo框架入门指南 - 完整教程",
    "meta_description": "本文详细介绍ClarkGo框架的使用方法",
    "meta_keywords": "ClarkGo, Go, Web框架, 教程"
  }'
```

---

## 📖 9. 查看内容

### 获取文章列表
```bash
curl http://localhost:8888/api/posts?page=1&per_page=10&status=published
```

### 获取文章详情
```bash
curl http://localhost:8888/api/posts/1
```

### 获取分类列表（树形）
```bash
curl http://localhost:8888/api/categories?tree=true
```

### 获取标签列表
```bash
curl http://localhost:8888/api/tags
```

### 获取媒体库
```bash
curl http://localhost:8888/api/media?file_type=image&page=1
```

---

## 🔧 10. 更新和删除

### 更新文章
```bash
curl -X PUT http://localhost:8888/api/cms/posts/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "ClarkGo框架完整指南（更新版）",
    "content": "更新后的内容..."
  }'
```

### 发布草稿文章
```bash
curl -X POST http://localhost:8888/api/cms/posts/1/publish \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 删除文章
```bash
curl -X DELETE http://localhost:8888/api/cms/posts/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📊 11. 常用查询示例

### 按分类筛选文章
```bash
curl "http://localhost:8888/api/posts?category_id=1&page=1"
```

### 按作者筛选文章
```bash
curl "http://localhost:8888/api/posts?author_id=1&page=1"
```

### 只查看草稿
```bash
curl "http://localhost:8888/api/posts?status=draft" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 按文件类型筛选媒体
```bash
curl "http://localhost:8888/api/media?file_type=image"
curl "http://localhost:8888/api/media?file_type=document"
```

---

## 🔐 12. 权限测试

### 测试无权限访问
```bash
# 使用普通用户token尝试删除文章（应该返回403）
curl -X DELETE http://localhost:8888/api/cms/posts/1 \
  -H "Authorization: Bearer NORMAL_USER_TOKEN"
```

预期响应：
```json
{
  "success": false,
  "error": "Forbidden",
  "message": "You don't have permission to delete post"
}
```

---

## 📝 13. 完整工作流示例

```bash
# 1. 登录
TOKEN=$(curl -s -X POST http://localhost:8888/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123456"}' \
  | jq -r '.data.token')

# 2. 上传特色图片
IMAGE_URL=$(curl -s -X POST http://localhost:8888/api/cms/media/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@image.jpg" \
  | jq -r '.data[0].file_url')

# 3. 创建分类
CATEGORY_ID=$(curl -s -X POST http://localhost:8888/api/cms/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"技术"}' \
  | jq -r '.data.id')

# 4. 创建标签
TAG_ID=$(curl -s -X POST http://localhost:8888/api/cms/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Go"}' \
  | jq -r '.data.id')

# 5. 创建并发布文章
curl -X POST http://localhost:8888/api/cms/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"我的第一篇文章\",
    \"content\": \"这是内容...\",
    \"featured_image\": \"$IMAGE_URL\",
    \"status\": \"published\",
    \"category_id\": $CATEGORY_ID,
    \"tags\": [$TAG_ID]
  }"

# 6. 查看文章
curl http://localhost:8888/api/posts
```

---

## 🐛 14. 故障排查

### 数据库连接问题
检查`.env`文件中的数据库配置：
```bash
DB_TYPE=sqlite
SQLITE_DATABASE=database/data.db
```

### JWT认证失败
确保：
1. Token正确复制（没有多余空格）
2. Token未过期
3. Authorization头格式正确：`Bearer TOKEN`

### 上传失败
检查：
1. 存储目录权限：`chmod -R 755 storage/uploads`
2. 文件大小是否超过限制（默认10MB）
3. 文件类型是否允许

### 权限被拒绝
确认：
1. 用户已分配正确的角色
2. 角色拥有所需权限
3. 使用正确的JWT token

---

## 📚 15. 下一步

- 阅读 [完整CMS文档](./CMS_IMPLEMENTATION.md)
- 查看 [API文档](./api.md)
- 了解 [权限系统](./CMS_IMPLEMENTATION.md#3-rbac权限管理系统)
- 探索 [邮件功能](./mail-api.md)
- 使用 [AI功能](./ai.md)

---

## 🎉 恭喜！

你已经成功设置并使用了ClarkGo CMS系统！现在可以开始构建你的内容管理应用了。

如有问题，请查阅文档或提交Issue。
