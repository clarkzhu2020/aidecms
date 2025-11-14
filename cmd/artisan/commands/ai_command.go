package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/clarkgo/clarkgo/pkg/ai"
)

// AICommand AI相关命令
type AICommand struct {
	manager *ai.Manager
}

// NewAICommand 创建AI命令
func NewAICommand() *AICommand {
	return &AICommand{
		manager: ai.NewManager(),
	}
}

// HandleAICommand 处理AI命令
func HandleAICommand(args []string) {
	cmd := NewAICommand()

	if len(args) < 1 {
		cmd.showHelp()
		return
	}

	subCommand := args[0]
	subArgs := args[1:]

	switch subCommand {
	case "setup":
		cmd.Setup(subArgs)
	case "chat":
		cmd.Chat(subArgs)
	case "completion":
		cmd.Completion(subArgs)
	case "models":
		cmd.ListModels(subArgs)
	case "test":
		cmd.Test(subArgs)
	case "config":
		cmd.Config(subArgs)
	default:
		fmt.Printf("Unknown AI command: %s\n", subCommand)
		cmd.showHelp()
	}
}

// Setup 设置AI配置
func (c *AICommand) Setup(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ai:setup <provider> <api_key> [model] [base_url]")
		fmt.Println("Supported providers: openai, anthropic, doubao, qianwen")
		return
	}

	provider := args[0]
	apiKey := args[1]
	model := "gpt-3.5-turbo" // 默认模型
	baseURL := ""

	if len(args) > 2 {
		model = args[2]
	}
	if len(args) > 3 {
		baseURL = args[3]
	}

	config := &ai.Config{
		Provider:    provider,
		APIKey:      apiKey,
		Model:       model,
		BaseURL:     baseURL,
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	// 创建配置目录
	if err := os.MkdirAll("config/ai", 0755); err != nil {
		fmt.Printf("Failed to create config directory: %v\n", err)
		return
	}

	// 保存配置
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal config: %v\n", err)
		return
	}

	configPath := fmt.Sprintf("config/ai/%s.json", provider)
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	fmt.Printf("AI configuration saved to %s\n", configPath)
	fmt.Printf("Provider: %s\n", provider)
	fmt.Printf("Model: %s\n", model)
	if baseURL != "" {
		fmt.Printf("Base URL: %s\n", baseURL)
	}
}

// loadAIConfig 加载AI配置
func (c *AICommand) loadAIConfig() error {
	configDir := "config/ai"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("AI config directory not found. Run 'ai:setup' first")
	}

	files, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		configPath := fmt.Sprintf("%s/%s", configDir, file.Name())
		configData, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Failed to read config %s: %v\n", file.Name(), err)
			continue
		}

		var config ai.Config
		if err := json.Unmarshal(configData, &config); err != nil {
			fmt.Printf("Failed to parse config %s: %v\n", file.Name(), err)
			continue
		}

		// 使用文件名（不包括.json扩展名）作为客户端名称
		clientName := strings.TrimSuffix(file.Name(), ".json")
		if err := c.manager.AddClient(clientName, &config); err != nil {
			fmt.Printf("Failed to add client %s: %v\n", clientName, err)
			continue
		}
	}

	return nil
}

// Chat 聊天命令
func (c *AICommand) Chat(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ai:chat <message> [model]")
		return
	}

	message := args[0]
	clientName := ""
	if len(args) > 1 {
		clientName = args[1]
	}

	// 加载配置
	if err := c.loadAIConfig(); err != nil {
		fmt.Printf("Error loading AI config: %v\n", err)
		return
	}

	// 获取客户端
	client, err := c.manager.GetClient(clientName)
	if err != nil {
		fmt.Printf("Error getting AI client: %v\n", err)
		return
	}

	// 执行聊天
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("You: %s\n", message)
	fmt.Print("AI: ")

	response, err := client.CreateCompletion(ctx, message)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(response)
}

// Completion 文本补全命令
func (c *AICommand) Completion(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ai:completion <prompt> [model] [temperature] [max_tokens]")
		return
	}

	prompt := args[0]
	clientName := ""
	temperature := 0.7
	maxTokens := 1000

	if len(args) > 1 {
		clientName = args[1]
	}
	if len(args) > 2 {
		if t, err := parseFloat(args[2]); err == nil {
			temperature = t
		}
	}
	if len(args) > 3 {
		if m, err := parseInt(args[3]); err == nil {
			maxTokens = m
		}
	}

	// 加载配置
	if err := c.loadAIConfig(); err != nil {
		fmt.Printf("Error loading AI config: %v\n", err)
		return
	}

	// 获取客户端
	client, err := c.manager.GetClient(clientName)
	if err != nil {
		fmt.Printf("Error getting AI client: %v\n", err)
		return
	}

	// 执行补全
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	options := []ai.Option{
		ai.WithTemperature(temperature),
		ai.WithMaxTokens(maxTokens),
	}

	response, err := client.CreateCompletion(ctx, prompt, options...)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(response)
}

// ListModels 列出可用模型
func (c *AICommand) ListModels(args []string) {
	// 加载配置
	if err := c.loadAIConfig(); err != nil {
		fmt.Printf("Error loading AI config: %v\n", err)
		return
	}

	models := c.manager.ListClients()
	if len(models) == 0 {
		fmt.Println("No AI models configured. Run 'ai:setup' to configure a model.")
		return
	}

	fmt.Println("Available AI models:")
	fmt.Println("==================")

	for _, name := range models {
		config, err := c.manager.GetConfig(name)
		if err != nil {
			continue
		}

		status := "✓"
		if name == c.manager.GetDefault() {
			status += " (default)"
		}

		fmt.Printf("%s %s\n", status, name)
		fmt.Printf("  Provider: %s\n", config.Provider)
		fmt.Printf("  Model: %s\n", config.Model)
		if config.BaseURL != "" {
			fmt.Printf("  Base URL: %s\n", config.BaseURL)
		}
		fmt.Println()
	}
}

// Test 测试AI连接
func (c *AICommand) Test(args []string) {
	clientName := ""
	if len(args) > 0 {
		clientName = args[0]
	}

	// 加载配置
	if err := c.loadAIConfig(); err != nil {
		fmt.Printf("Error loading AI config: %v\n", err)
		return
	}

	// 获取客户端
	client, err := c.manager.GetClient(clientName)
	if err != nil {
		fmt.Printf("Error getting AI client: %v\n", err)
		return
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testPrompt := "Hello, this is a test message. Please respond with 'Test successful'."

	fmt.Printf("Testing AI client: %s\n", func() string {
		if clientName == "" {
			return c.manager.GetDefault()
		}
		return clientName
	}())

	fmt.Print("Testing connection... ")

	response, err := client.CreateCompletion(ctx, testPrompt)
	if err != nil {
		fmt.Printf("✗ Failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Success\n")
	fmt.Printf("Response: %s\n", response)
}

// Config 配置管理
func (c *AICommand) Config(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ai:config <action> [args...]")
		fmt.Println("Actions:")
		fmt.Println("  list                - List all configurations")
		fmt.Println("  show <provider>     - Show specific configuration")
		fmt.Println("  delete <provider>   - Delete configuration")
		fmt.Println("  default <provider>  - Set default provider")
		return
	}

	action := args[0]
	actionArgs := args[1:]

	switch action {
	case "list":
		c.listConfigs()
	case "show":
		c.showConfig(actionArgs)
	case "delete":
		c.deleteConfig(actionArgs)
	case "default":
		c.setDefault(actionArgs)
	default:
		fmt.Printf("Unknown config action: %s\n", action)
	}
}

// listConfigs 列出所有配置
func (c *AICommand) listConfigs() {
	configDir := "config/ai"
	files, err := os.ReadDir(configDir)
	if err != nil {
		fmt.Printf("No AI configurations found\n")
		return
	}

	fmt.Println("AI Configurations:")
	fmt.Println("==================")

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			provider := strings.TrimSuffix(file.Name(), ".json")
			fmt.Printf("- %s\n", provider)
		}
	}
}

// showConfig 显示配置
func (c *AICommand) showConfig(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ai:config show <provider>")
		return
	}

	provider := args[0]
	configPath := fmt.Sprintf("config/ai/%s.json", provider)

	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Configuration for %s not found\n", provider)
		return
	}

	var config ai.Config
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("Failed to parse configuration: %v\n", err)
		return
	}

	// 隐藏API密钥
	config.APIKey = maskAPIKey(config.APIKey)

	configJSON, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Configuration for %s:\n%s\n", provider, configJSON)
}

// deleteConfig 删除配置
func (c *AICommand) deleteConfig(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ai:config delete <provider>")
		return
	}

	provider := args[0]
	configPath := fmt.Sprintf("config/ai/%s.json", provider)

	if err := os.Remove(configPath); err != nil {
		fmt.Printf("Failed to delete configuration for %s: %v\n", provider, err)
		return
	}

	fmt.Printf("Configuration for %s deleted\n", provider)
}

// setDefault 设置默认提供商
func (c *AICommand) setDefault(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ai:config default <provider>")
		return
	}

	provider := args[0]

	// 检查配置是否存在
	configPath := fmt.Sprintf("config/ai/%s.json", provider)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Configuration for %s not found\n", provider)
		return
	}

	// 保存默认配置
	defaultConfig := map[string]string{
		"default": provider,
	}

	defaultData, _ := json.MarshalIndent(defaultConfig, "", "  ")
	if err := os.WriteFile("config/ai/default.json", defaultData, 0644); err != nil {
		fmt.Printf("Failed to set default provider: %v\n", err)
		return
	}

	fmt.Printf("Default AI provider set to: %s\n", provider)
}

// showHelp 显示帮助信息
func (c *AICommand) showHelp() {
	fmt.Println("AI Commands:")
	fmt.Println("============")
	fmt.Println("  ai:setup <provider> <api_key> [model] [base_url]  - Setup AI configuration")
	fmt.Println("  ai:chat <message> [model]                         - Chat with AI")
	fmt.Println("  ai:completion <prompt> [model] [temp] [tokens]    - Text completion")
	fmt.Println("  ai:models                                         - List available models")
	fmt.Println("  ai:test [model]                                   - Test AI connection")
	fmt.Println("  ai:config <action> [args...]                     - Manage configurations")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ai:setup openai sk-xxx gpt-4")
	fmt.Println("  ai:chat \"Hello, how are you?\" openai")
	fmt.Println("  ai:completion \"Write a poem about\" openai 0.8 500")
	fmt.Println("  ai:test openai")
}

// parseFloat 解析浮点数
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// parseInt 解析整数
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// maskAPIKey 隐藏API密钥
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "****"
	}
	return apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
}
