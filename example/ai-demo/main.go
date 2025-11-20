package main

import (
	"context"
	"fmt"
	"log"

	"github.com/clarkzhu2020/aidecms/pkg/ai"
)

func main() {
	fmt.Println("AideCMS AI Demo")
	fmt.Println("===============")

	// 创建 AI 管理器
	aiManager := ai.NewManager()

	// 示例：添加一个测试配置（实际应用中应该从配置文件加载）
	testConfig := &ai.Config{
		Provider:    "openai",
		APIKey:      "test-key",
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	err := aiManager.AddClient("test", testConfig)
	if err != nil {
		log.Printf("Failed to add test client: %v", err)
	}

	// 显示可用的 AI 模型
	models := aiManager.ListClients()
	if len(models) > 0 {
		fmt.Printf("Available AI models: %v\n", models)
		fmt.Printf("Default model: %s\n", aiManager.GetDefault())
	} else {
		fmt.Println("No AI models configured.")
		fmt.Println("Run the following command to setup an AI model:")
		fmt.Println("  go run cmd/artisan/main.go ai:setup openai your-api-key gpt-4")
	}

	// 示例：如何使用AI客户端
	if len(models) > 0 {
		client, err := aiManager.GetDefaultClient()
		if err != nil {
			fmt.Printf("Error getting default client: %v\n", err)
			return
		}

		fmt.Println("\nTesting AI client (mock response)...")
		ctx := context.Background()
		response, err := client.CreateCompletion(ctx, "Hello, this is a test.")
		if err != nil {
			fmt.Printf("AI test failed: %v\n", err)
		} else {
			fmt.Printf("AI response: %s\n", response)
		}
	}

	fmt.Println("\nAI Integration Features:")
	fmt.Println("- Multiple AI provider support (OpenAI, Anthropic, Doubao, etc.)")
	fmt.Println("- Unified API interface")
	fmt.Println("- Conversation management with context")
	fmt.Println("- Streaming responses")
	fmt.Println("- Embedding generation")
	fmt.Println("- Command-line tools")
	fmt.Println("- HTTP middleware")
	fmt.Println()
	fmt.Println("To start the main web server with AI endpoints:")
	fmt.Println("  go run main.go")
} // 简单的聊天界面 HTML
const chatHTML = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AideCMS AI Chat</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .chat-container {
            background: white;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .chat-header {
            background: #007AFF;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .chat-messages {
            height: 400px;
            overflow-y: auto;
            padding: 20px;
            border-bottom: 1px solid #eee;
        }
        .message {
            margin: 10px 0;
            padding: 10px;
            border-radius: 10px;
            max-width: 80%;
        }
        .user-message {
            background: #007AFF;
            color: white;
            margin-left: auto;
        }
        .ai-message {
            background: #f0f0f0;
            color: #333;
        }
        .chat-input {
            padding: 20px;
            display: flex;
            gap: 10px;
        }
        .chat-input input {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 16px;
        }
        .chat-input button {
            padding: 10px 20px;
            background: #007AFF;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
        }
        .chat-input button:hover {
            background: #0056b3;
        }
        .chat-input button:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
    </style>
</head>
<body>
    <div class="chat-container">
        <div class="chat-header">
            <h2>AideCMS AI Chat Demo</h2>
            <p>基于 CloudWeGo Eino 的 AI 聊天演示</p>
        </div>
        <div class="chat-messages" id="messages"></div>
        <div class="chat-input">
            <input type="text" id="messageInput" placeholder="输入您的消息..." onkeypress="handleEnter(event)">
            <button onclick="sendMessage()" id="sendButton">发送</button>
        </div>
    </div>

    <script>
        const messagesContainer = document.getElementById('messages');
        const messageInput = document.getElementById('messageInput');
        const sendButton = document.getElementById('sendButton');
        let sessionId = 'demo_' + Date.now();

        function addMessage(content, isUser) {
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message ' + (isUser ? 'user-message' : 'ai-message');
            messageDiv.textContent = content;
            messagesContainer.appendChild(messageDiv);
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }

        function handleEnter(event) {
            if (event.key === 'Enter') {
                sendMessage();
            }
        }

        async function sendMessage() {
            const message = messageInput.value.trim();
            if (!message) return;

            // 添加用户消息到界面
            addMessage(message, true);
            messageInput.value = '';
            sendButton.disabled = true;
            sendButton.textContent = '发送中...';

            try {
                // 发送请求到 AI API
                const response = await fetch('/api/ai/conversation', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        session_id: sessionId,
                        message: message
                    })
                });

                const data = await response.json();
                
                if (response.ok) {
                    // 添加 AI 回复到界面
                    addMessage(data.message, false);
                } else {
                    addMessage('错误: ' + (data.error || '服务暂时不可用'), false);
                }
            } catch (error) {
                console.error('Error:', error);
                addMessage('错误: 网络连接失败', false);
            } finally {
                sendButton.disabled = false;
                sendButton.textContent = '发送';
            }
        }

        // 页面加载时的欢迎消息
        addMessage('你好！我是 AideCMS AI 助手。请问有什么可以帮助您的？', false);
    </script>
</body>
</html>
`
