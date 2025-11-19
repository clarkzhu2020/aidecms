# Swagger API 文档使用指南

## 📚 概述

ClarkGo CMS 已集成 Swagger UI，提供交互式 API 文档界面。

## 🚀 访问 Swagger UI

启动服务后，访问：
```
http://localhost:8888/swagger/index.html
```

## 📝 已文档化的 API 端点

### Posts (文章管理)
- `GET /api/posts` - 获取文章列表
- `GET /api/posts/{id}` - 获取文章详情
- `POST /api/cms/posts` - 创建文章 🔒
- `PUT /api/cms/posts/{id}` - 更新文章 🔒
- `DELETE /api/cms/posts/{id}` - 删除文章 🔒
- `POST /api/cms/posts/{id}/publish` - 发布文章 🔒

### Categories (分类管理)
- `GET /api/categories` - 获取分类列表（支持树形结构）
- `POST /api/cms/categories` - 创建分类 🔒

### Tags (标签管理)
- `GET /api/tags` - 获取标签列表
- `POST /api/cms/tags` - 创建标签 🔒

### Media (媒体管理)
- `POST /api/cms/media/upload` - 上传文件 🔒

*🔒 表示需要 JWT 认证*

## 🔐 API 认证测试

### 1. 在 Swagger UI 中认证

1. 点击右上角的 **Authorize** 按钮
2. 输入格式：`Bearer YOUR_JWT_TOKEN`
3. 点击 **Authorize**
4. 点击 **Close**

### 2. 获取 JWT Token

首先需要登录获取 token（这个接口暂未添加 Swagger 注解，需要手动调用）：

```bash
curl -X POST http://localhost:8888/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123456"
  }'
```

响应中包含 token：
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. 在 Swagger UI 中使用

将获取到的 token 在 Authorize 对话框中输入：
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 📖 API 使用示例

### 示例 1: 获取文章列表

在 Swagger UI 中：
1. 展开 `GET /api/posts`
2. 点击 **Try it out**
3. 设置参数（可选）：
   - `page`: 1
   - `per_page`: 10
   - `status`: published
4. 点击 **Execute**
5. 查看响应

### 示例 2: 创建文章

1. 确保已经完成认证（Authorize）
2. 展开 `POST /api/cms/posts`
3. 点击 **Try it out**
4. 编辑请求体 JSON：
```json
{
  "title": "我的第一篇文章",
  "content": "这是文章内容...",
  "excerpt": "文章摘要",
  "status": "draft",
  "category_id": 1,
  "tags": [1, 2]
}
```
5. 点击 **Execute**
6. 查看响应

### 示例 3: 上传文件

1. 确保已认证
2. 展开 `POST /api/cms/media/upload`
3. 点击 **Try it out**
4. 点击 **Choose File** 选择文件
5. 点击 **Execute**
6. 查看上传结果

## 🔧 重新生成文档

当你修改了 API 接口或添加了新的注解后，需要重新生成 Swagger 文档：

```bash
# 安装 swag 工具（如果还没安装）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g main.go --output ./docs

# 重新编译项目
go build -o clarkgo main.go

# 启动服务
./clarkgo
```

## 📐 添加 Swagger 注解示例

### 为控制器方法添加注解

```go
// Create 创建文章
// @Summary      创建文章
// @Description  创建一篇新文章
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        post body CreatePostRequest true "文章信息"
// @Success      201 {object} response.Response{data=models.PostSwagger}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Security     BearerAuth
// @Router       /cms/posts [post]
func (c *PostController) Create(ctx context.Context, hCtx *app.RequestContext) {
    // 实现代码...
}
```

### 常用注解说明

- `@Summary` - 简短描述（显示在列表中）
- `@Description` - 详细描述
- `@Tags` - 分组标签
- `@Accept` - 接受的内容类型
- `@Produce` - 返回的内容类型
- `@Param` - 参数说明
  - 格式：`参数名 位置 类型 必需 "描述"`
  - 位置：path, query, body, header, formData
- `@Success` - 成功响应
- `@Failure` - 失败响应
- `@Security` - 安全认证（BearerAuth）
- `@Router` - 路由路径和方法

## 📊 响应模型

### 统一响应格式

```go
// 成功响应
{
  "success": true,
  "data": { /* 数据对象 */ },
  "message": "操作成功"
}

// 分页响应
{
  "success": true,
  "data": [ /* 数据数组 */ ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}

// 错误响应
{
  "success": false,
  "error": "BadRequest",
  "message": "错误描述"
}

// 验证错误响应
{
  "success": false,
  "error": "ValidationError",
  "message": "Validation failed",
  "errors": {
    "title": ["标题不能为空", "标题长度必须在3-200之间"]
  }
}
```

## 🎯 最佳实践

### 1. API 设计规范

- 使用 RESTful 风格
- GET：查询操作
- POST：创建操作
- PUT/PATCH：更新操作
- DELETE：删除操作

### 2. 路径命名

- 公开路由：`/api/资源`
- 需认证路由：`/api/cms/资源`

### 3. 状态码使用

- `200 OK` - 成功
- `201 Created` - 创建成功
- `400 Bad Request` - 请求参数错误
- `401 Unauthorized` - 未认证
- `403 Forbidden` - 无权限
- `404 Not Found` - 资源不存在
- `422 Unprocessable Entity` - 验证失败
- `500 Internal Server Error` - 服务器错误

## 🔍 调试技巧

### 1. 查看原始请求

在 Swagger UI 中执行请求后，可以看到：
- **Curl** - 等效的 curl 命令
- **Request URL** - 完整的请求 URL
- **Request Headers** - 请求头
- **Response Body** - 响应体
- **Response Headers** - 响应头

### 2. 导出 API 规范

下载 API 规范文件：
- JSON 格式：`http://localhost:8888/swagger/doc.json`
- YAML 格式：直接访问 `docs/swagger.yaml` 文件

### 3. 在 Postman 中使用

1. 访问 `http://localhost:8888/swagger/doc.json`
2. 复制 JSON 内容
3. 在 Postman 中：**Import** → **Paste Raw Text** → 粘贴 JSON
4. 导入后即可在 Postman 中测试所有 API

## 📚 参考资源

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Swagger UI 文档](https://swagger.io/docs/open-source-tools/swagger-ui/)
- [OpenAPI 3.0 规范](https://swagger.io/specification/)

## 🆘 常见问题

### Q: 修改了注解但文档没更新？
**A:** 需要重新运行 `swag init` 命令生成文档，然后重启服务。

### Q: 为什么有些接口没有显示？
**A:** 确保控制器方法上方有 Swagger 注解，并且注解格式正确。

### Q: 如何在 Swagger UI 中测试文件上传？
**A:** 
1. 确保已认证
2. 找到文件上传接口
3. 点击 **Try it out**
4. 点击 **Choose File** 按钮选择文件
5. 点击 **Execute**

### Q: Bearer token 格式是什么？
**A:** 格式为 `Bearer ` + 空格 + JWT token，例如：
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 🎉 下一步

- 为更多控制器添加 Swagger 注解
- 添加请求/响应示例
- 集成 API 测试套件
- 设置 API 版本控制
