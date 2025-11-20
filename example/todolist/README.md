# TodoList 示例应用

这是一个使用AideCMS框架开发的TodoList示例应用，展示了如何构建一个简单的RESTful API。

## 功能特性

- 创建Todo项
- 获取所有Todo项
- 获取单个Todo项
- 更新Todo项
- 删除Todo项

## 快速开始

1. 确保已安装Go 1.16+
2. 克隆项目
3. 进入示例目录：
   ```bash
   cd example/todolist
   ```
4. 运行应用：
   ```bash
   go run cmd/main.go
   ```

## API 端点

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /todos | 获取所有Todo项 |
| POST | /todos | 创建新的Todo项 |
| GET | /todos/:id | 获取单个Todo项 |
| PUT | /todos/:id | 更新Todo项 |
| DELETE | /todos/:id | 删除Todo项 |

## 请求示例

### 创建Todo
```bash
curl -X POST http://localhost:8888/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"新任务","completed":false}'
```

### 获取所有Todo
```bash
curl http://localhost:8888/todos
```

### 更新Todo
```bash
curl -X PUT http://localhost:8888/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"更新后的任务","completed":true}'
```

### 删除Todo
```bash
curl -X DELETE http://localhost:8888/todos/1
```

## 开发说明

1. 控制器位于 `internal/controllers/todo_controller.go`
2. 模型位于 `internal/models/todo.go`
3. 主入口文件位于 `cmd/main.go`

## 下一步

- 添加数据库持久化
- 实现用户认证
- 添加输入验证
- 编写单元测试