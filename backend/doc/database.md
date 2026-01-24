# 数据库

AideCMS 使用 GORM 作为 ORM，支持 MySQL、PostgreSQL、SQLite 等数据库。

## 配置

数据库配置在 `.env` 文件中:

```
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=aidecms
DB_USERNAME=root
DB_PASSWORD=
```

## 模型定义

```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name     string
    Email    string `gorm:"unique"`
    Posts    []Post
}

type Post struct {
    gorm.Model
    Title    string
    Content  string
    UserID   uint
    User     User
}
```

## 基本操作

### 创建记录

```go
user := models.User{
    Name:  "John",
    Email: "john@example.com",
}

result := db.Create(&user)
if result.Error != nil {
    // 处理错误
}
```

### 查询

```go
// 获取单个用户
var user models.User
db.First(&user, 1) // 通过主键查询
db.First(&user, "email = ?", "john@example.com")

// 获取用户列表
var users []models.User
db.Find(&users)

// Where 条件
db.Where("name = ?", "John").Find(&users)
```

### 更新

```go
db.Model(&user).Update("name", "Johnny")
```

### 删除

```go
db.Delete(&user)
```

## 关联关系

### 一对多

```go
// 获取用户的所有文章
db.Model(&user).Association("Posts").Find(&posts)

// 添加文章
db.Model(&user).Association("Posts").Append(&post)
```

### 预加载

```go
// 预加载用户的所有文章
db.Preload("Posts").First(&user)
```

## 迁移

使用 Artisan 命令创建迁移:

```bash
go run cmd/artisan/main.go make:migration create_users_table
```

运行迁移:

```bash
go run cmd/artisan/main.go migrate
```

回滚迁移:

```bash
go run cmd/artisan/main.go migrate:rollback
```

## 种子数据

创建种子:

```bash
go run cmd/artisan/main.go make:seeder UsersTableSeeder
```

运行种子:

```bash
go run cmd/artisan/main.go db:seed