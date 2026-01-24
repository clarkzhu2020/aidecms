# AideCMS AI 集成指南

AideCMS 集成了 CloudWeGo Eino 框架，提供强大的 AI 大模型交互能力。支持多种主流大模型提供商，包括 OpenAI、Anthropic、豆包、通义千问等。

## 特性

- 🤖 **多模型支持**: 支持 OpenAI、Anthropic、豆包、通义千问、ChatGLM 等主流大模型
- 🔄 **统一接口**: 提供一致的 API 接口，轻松切换不同模型
- 💬 **对话管理**: 自动管理对话上下文，支持多轮对话
- 🌊 **流式输出**: 支持流式响应，提供更好的用户体验
- 🛠️ **命令行工具**: 内置 Artisan 命令，便于测试和管理
- 🔧 **中间件支持**: 提供 HTTP 中间件，快速构建 AI 应用
- 📊 **嵌入向量**: 支持文本嵌入向量生成
- ⚙️ **灵活配置**: 支持多环境配置和动态切换

## 快速开始

### 1. 配置 AI 服务

使用 Artisan 命令配置 AI 服务：

```bash
# 配置 OpenAI
go run cmd/artisan/main.go ai:setup openai sk-your-api-key gpt-4

# 配置豆包
go run cmd/artisan/main.go ai:setup doubao your-api-key ep-xxx

# 配置通义千问
go run cmd/artisan/main.go ai:setup qianwen your-api-key qwen-max
```

### 2. 测试连接

```bash
# 测试默认模型
go run cmd/artisan/main.go ai:test

# 测试指定模型
go run cmd/artisan/main.go ai:test openai
```

### 3. 命令行聊天

```bash
# 简单对话
go run cmd/artisan/main.go ai:chat "Hello, how are you?"

# 指定模型对话
go run cmd/artisan/main.go ai:chat "写一首关于春天的诗" qianwen

# 文本补全
go run cmd/artisan/main.go ai:completion "Once upon a time" openai 0.8 500
```

## API 使用

### 1. 在应用中集成 AI 管理器

```go
package main

import (
    "github.com/chenyusolar/aidecms/pkg/framework"
    "github.com/chenyusolar/aidecms/pkg/framework/middleware"
    "github.com/chenyusolar/aidecms/config"
    "github.com/chenyusolar/aidecms/app/Http/Controllers"
)

func main() {
    app := framework.NewApplication().Boot()
    
    // 加载 AI 管理器
    aiManager, err := config.LoadAIManager()
    if err != nil {
        panic(err)
    }
    
    // 注册 AI 中间件
    aiMiddleware := middleware.NewAIMiddleware(aiManager)
    app.RegisterMiddleware(aiMiddleware.Handler())
    
    // 创建 AI 控制器
    aiController := controllers.NewAIController(aiManager)
    
    // 注册 AI 路由
    app.RegisterRoutes(func(router *framework.Router) {
        api := router.Group("/api/ai")
        {
            api.POST("/chat", aiController.Chat)
            api.POST("/completion", aiController.Completion)
            api.POST("/embedding", aiController.Embedding)
            api.POST("/conversation", aiController.Conversation)
            api.GET("/models", aiController.Models)
            api.GET("/health", aiController.Health)
        }
    })
    
    app.Run()
}
```

### 2. 聊天 API

**请求**:
```bash
curl -X POST http://localhost:8888/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好，请介绍一下 Go 语言",
    "model": "qianwen",
    "temperature": 0.7,
    "max_tokens": 1000
  }'
```

**响应**:
```json
{
  "message": "Go 是 Google 开发的开源编程语言...",
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 200,
    "total_tokens": 210
  },
  "model": "qianwen"
}
```

### 3. 流式聊天 API

**请求**:
```bash
curl -X POST http://localhost:8888/api/ai/chat \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{
    "message": "写一个 Go 程序示例",
    "stream": true
  }'
```

**响应** (Server-Sent Events):
```
data: {"message": "这里", "model": "qianwen"}

data: {"message": "这里是一个", "model": "qianwen"}

data: {"message": "这里是一个简单的 Go", "model": "qianwen"}

data: [DONE]
```

### 4. 对话上下文 API

支持多轮对话，自动管理上下文：

```bash
# 开始新对话
curl -X POST http://localhost:8888/api/ai/conversation \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "user123_session1",
    "message": "我想学习 Go 语言",
    "model": "qianwen"
  }'

# 继续对话
curl -X POST http://localhost:8888/api/ai/conversation \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "user123_session1",
    "message": "Go 语言有什么特点？"
  }'

# 查看对话历史
curl "http://localhost:8888/api/ai/conversation/user123_session1/history"

# 清空对话历史
curl -X DELETE "http://localhost:8888/api/ai/conversation/user123_session1/history"
```

### 5. 文本嵌入 API

```bash
curl -X POST http://localhost:8888/api/ai/embedding \
  -H "Content-Type: application/json" \
  -d '{
    "input": ["Hello world", "你好世界"],
    "model": "openai"
  }'
```

## 高级用法

### 1. 自定义中间件

```go
// 创建带系统提示的聊天中间件
chatConfig := &middleware.ChatMiddlewareConfig{
    ClientName:   "qianwen",
    SystemPrompt: "你是一个专业的 Go 语言助手，请用中文回答问题。",
    AutoResponse: true,
}
chatMiddleware := middleware.NewChatMiddleware(aiManager, chatConfig)

// 应用到特定路由组
chatGroup := router.Group("/chat")
chatGroup.Use(chatMiddleware.Handler())
```

### 2. 程序化使用

```go
package main

import (
    "context"
    "fmt"
    "github.com/chenyusolar/aidecms/pkg/ai"
)

func example() {
    // 创建 AI 客户端
    config := &ai.Config{
        Provider:    "openai",
        APIKey:      "sk-your-key",
        Model:       "gpt-4",
        Temperature: 0.7,
        MaxTokens:   1000,
    }
    
    client, err := ai.NewClient(config)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 简单聊天
    response, err := client.CreateCompletion(context.Background(), 
        "请解释一下 Go 语言的并发模型")
    if err != nil {
        panic(err)
    }
    fmt.Println(response)
    
    // 对话聊天
    messages := []*ai.Message{
        {Role: "system", Content: "你是一个 Go 语言专家"},
        {Role: "user", Content: "什么是 goroutine？"},
    }
    
    chatResp, err := client.Chat(context.Background(), &ai.ChatRequest{
        Messages: messages,
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(chatResp.Message.Content)
}
```

### 3. 对话管理

```go
// 创建对话客户端
modelConfig := &ai.ModelConfig{
    Provider: ai.ProviderOpenAI,
    APIKey:   "sk-your-key",
    Model:    "gpt-4",
}

einoClient, _ := ai.NewEinoClient(modelConfig)
conversationClient := ai.NewConversationClient(einoClient, 100)

// 多轮对话
response1, _ := conversationClient.Chat(ctx, "你好，我想学习编程")
response2, _ := conversationClient.Chat(ctx, "推荐学习 Go 语言吗？")

// 获取历史记录
history := conversationClient.GetHistory()
```

## 配置管理

### 1. 命令行配置管理

```bash
# 列出所有配置
go run cmd/artisan/main.go ai:config list

# 显示特定配置
go run cmd/artisan/main.go ai:config show openai

# 删除配置
go run cmd/artisan/main.go ai:config delete openai

# 设置默认提供商
go run cmd/artisan/main.go ai:config default qianwen
```

### 2. 配置文件

配置文件位于 `config/ai/` 目录下，每个提供商一个 JSON 文件：

**config/ai/openai.json**:
```json
{
  "provider": "openai",
  "api_key": "sk-your-api-key",
  "model": "gpt-4",
  "base_url": "",
  "temperature": 0.7,
  "max_tokens": 1000,
  "options": {}
}
```

**config/ai/qianwen.json**:
```json
{
  "provider": "qianwen",
  "api_key": "your-api-key",
  "model": "qwen-max",
  "base_url": "https://dashscope.aliyuncs.com/api/v1",
  "temperature": 0.8,
  "max_tokens": 2000,
  "options": {
    "region": "cn-beijing"
  }
}
```

## 支持的模型提供商

### 1. OpenAI
- **模型**: gpt-4, gpt-4-turbo, gpt-3.5-turbo
- **功能**: 聊天、补全、嵌入
- **配置**: API Key + Base URL (可选)

### 2. Anthropic
- **模型**: claude-3-opus, claude-3-sonnet, claude-3-haiku
- **功能**: 聊天、补全
- **配置**: API Key

### 3. 字节豆包
- **模型**: ep-xxx 格式的端点
- **功能**: 聊天、补全、嵌入
- **配置**: API Key + 端点 ID

### 4. 通义千问
- **模型**: qwen-max, qwen-plus, qwen-turbo
- **功能**: 聊天、补全、嵌入
- **配置**: API Key + 地域 (可选)

### 5. 智谱 ChatGLM
- **模型**: glm-4, glm-3-turbo
- **功能**: 聊天、补全
- **配置**: API Key

## 最佳实践

### 1. 错误处理
```go
response, err := client.Chat(ctx, req)
if err != nil {
    // 记录错误日志
    log.Printf("AI chat failed: %v", err)
    
    // 返回友好的错误信息
    return "抱歉，AI 服务暂时不可用，请稍后再试"
}
```

### 2. 超时控制
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := client.Chat(ctx, req)
```

### 3. 流式响应
```go
responseCh, errorCh := client.StreamChat(ctx, req)
for {
    select {
    case response := <-responseCh:
        // 处理流式响应
        fmt.Print(response.Message.Content)
    case err := <-errorCh:
        // 处理错误
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### 4. 配置管理
- 使用环境变量存储敏感信息（API Key）
- 为不同环境配置不同的模型
- 合理设置超时和重试机制

### 5. 监控和日志
- 记录 AI 请求和响应日志
- 监控 API 调用频率和成功率
- 设置告警机制

## 故障排除

### 常见问题

1. **配置不存在**
   ```bash
   Error loading AI config: AI config directory not found. Run 'ai:setup' first
   ```
   解决方案: 运行 `ai:setup` 命令配置 AI 服务

2. **API Key 无效**
   ```bash
   Error: Model not available: failed to create client: api_key is required
   ```
   解决方案: 检查 API Key 是否正确配置

3. **网络连接失败**
   ```bash
   Error: Chat failed: context deadline exceeded
   ```
   解决方案: 检查网络连接和 Base URL 配置

4. **模型不存在**
   ```bash
   Error: Model not available: client not found
   ```
   解决方案: 检查模型名称是否正确配置

### 调试技巧

1. 使用 `ai:test` 命令测试连接
2. 查看日志文件 `storage/logs/artisan.log`
3. 使用 `ai:config show` 检查配置
4. 使用 `ai:models` 列出可用模型

## 更多信息

- [CloudWeGo Eino 官方文档](https://github.com/cloudwego/eino)
- [AideCMS 项目主页](https://github.com/clarkzhu2020/aidecms)
- [API 参考文档](./api.md)
- [部署指南](./deployment.md)