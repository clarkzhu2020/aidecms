package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clarkgo/clarkgo/pkg/ai"
	envConfig "github.com/clarkgo/clarkgo/pkg/config"
)

// AIConfig AI配置结构
type AIConfig struct {
	DefaultProvider string                 `json:"default_provider"`
	Providers       map[string]*ai.Config  `json:"providers"`
	GlobalOptions   map[string]interface{} `json:"global_options"`
}

// LoadAIConfig 加载AI配置
func LoadAIConfig() (*AIConfig, error) {
	// 首先加载 .env 文件
	if err := envConfig.LoadEnv(".env"); err != nil {
		// 如果加载失败，尝试加载 .env.example
		envConfig.LoadEnv(".env.example")
	}

	configPath := "config/ai.json"

	// 优先从环境变量加载配置
	config := loadConfigFromEnv()

	// 如果存在配置文件，则合并配置文件中的设置
	if _, err := os.Stat(configPath); err == nil {
		fileConfig, err := loadConfigFromFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load AI config file: %w", err)
		}

		// 合并配置
		mergeConfigs(config, fileConfig)
	}

	return config, nil
}

// loadConfigFromEnv 从环境变量加载配置
func loadConfigFromEnv() *AIConfig {
	config := &AIConfig{
		Providers:     make(map[string]*ai.Config),
		GlobalOptions: make(map[string]interface{}),
	}

	// 设置默认提供商
	config.DefaultProvider = envConfig.GetEnv("AI_DEFAULT_PROVIDER", "")

	// 加载各个提供商的配置
	providers := []string{"openai", "anthropic", "doubao", "qianwen", "chatglm", "baichuan", "minimax"}

	for _, provider := range providers {
		providerConfig := loadProviderConfig(provider)
		if providerConfig != nil {
			config.Providers[provider] = providerConfig
		}
	} // 加载全局选项
	config.GlobalOptions["request_timeout"] = envConfig.GetEnvInt("AI_REQUEST_TIMEOUT", 30)
	config.GlobalOptions["max_retries"] = envConfig.GetEnvInt("AI_MAX_RETRIES", 3)
	config.GlobalOptions["retry_delay"] = envConfig.GetEnvInt("AI_RETRY_DELAY", 1)
	config.GlobalOptions["rate_limit_enabled"] = envConfig.GetEnvBool("AI_RATE_LIMIT_ENABLED", true)
	config.GlobalOptions["rate_limit_rpm"] = envConfig.GetEnvInt("AI_RATE_LIMIT_RPM", 60)
	config.GlobalOptions["conversation_max_history"] = envConfig.GetEnvInt("AI_CONVERSATION_MAX_HISTORY", 50)
	config.GlobalOptions["embedding_batch_size"] = envConfig.GetEnvInt("AI_EMBEDDING_BATCH_SIZE", 100)

	// 功能开关
	config.GlobalOptions["chat_enabled"] = envConfig.GetEnvBool("AI_CHAT_ENABLED", true)
	config.GlobalOptions["completion_enabled"] = envConfig.GetEnvBool("AI_COMPLETION_ENABLED", true)
	config.GlobalOptions["embedding_enabled"] = envConfig.GetEnvBool("AI_EMBEDDING_ENABLED", true)
	config.GlobalOptions["streaming_enabled"] = envConfig.GetEnvBool("AI_STREAMING_ENABLED", true)

	// 安全设置
	config.GlobalOptions["content_filter_enabled"] = envConfig.GetEnvBool("AI_CONTENT_FILTER_ENABLED", true)
	config.GlobalOptions["api_key_encryption"] = envConfig.GetEnvBool("AI_API_KEY_ENCRYPTION", false)
	config.GlobalOptions["audit_log_enabled"] = envConfig.GetEnvBool("AI_AUDIT_LOG_ENABLED", true)

	// 缓存设置
	config.GlobalOptions["cache_enabled"] = envConfig.GetEnvBool("AI_CACHE_ENABLED", true)
	config.GlobalOptions["cache_ttl"] = envConfig.GetEnvInt("AI_CACHE_TTL", 3600)
	config.GlobalOptions["cache_max_size"] = envConfig.GetEnvInt("AI_CACHE_MAX_SIZE", 1000)

	// 监控设置
	config.GlobalOptions["metrics_enabled"] = envConfig.GetEnvBool("AI_METRICS_ENABLED", true)
	config.GlobalOptions["health_check_enabled"] = envConfig.GetEnvBool("AI_HEALTH_CHECK_ENABLED", true)
	config.GlobalOptions["performance_logging"] = envConfig.GetEnvBool("AI_PERFORMANCE_LOGGING", true)

	return config
}

// loadProviderConfig 加载特定提供商的配置
func loadProviderConfig(provider string) *ai.Config {
	upperProvider := strings.ToUpper(provider)

	// 检查是否有 API 密钥
	apiKey := envConfig.GetEnv(upperProvider+"_API_KEY", "")
	if apiKey == "" || strings.Contains(apiKey, "your-") {
		return nil // 没有有效的 API 密钥，跳过此提供商
	}

	providerConfig := &ai.Config{
		Provider:    provider,
		APIKey:      apiKey,
		BaseURL:     envConfig.GetEnv(upperProvider+"_API_BASE", ""),
		Model:       envConfig.GetEnv(upperProvider+"_MODEL", ""),
		Temperature: envConfig.GetEnvFloat64(upperProvider+"_TEMPERATURE", 0.7),
		MaxTokens:   envConfig.GetEnvInt(upperProvider+"_MAX_TOKENS", 2000),
		Options:     make(map[string]interface{}),
	}

	// 特殊处理 MiniMax 的 API Secret
	if provider == "minimax" {
		if apiSecret := envConfig.GetEnv("MINIMAX_API_SECRET", ""); apiSecret != "" {
			providerConfig.Options["api_secret"] = apiSecret
		}
	}

	return providerConfig
} // loadConfigFromFile 从配置文件加载配置
func loadConfigFromFile(configPath string) (*AIConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config AIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// mergeConfigs 合并配置
func mergeConfigs(envConfig, fileConfig *AIConfig) {
	// 如果环境变量中没有设置默认提供商，使用文件配置
	if envConfig.DefaultProvider == "" && fileConfig.DefaultProvider != "" {
		envConfig.DefaultProvider = fileConfig.DefaultProvider
	}

	// 合并提供商配置
	for name, providerConfig := range fileConfig.Providers {
		if _, exists := envConfig.Providers[name]; !exists {
			envConfig.Providers[name] = providerConfig
		}
	}

	// 合并全局选项
	for key, value := range fileConfig.GlobalOptions {
		if _, exists := envConfig.GlobalOptions[key]; !exists {
			envConfig.GlobalOptions[key] = value
		}
	}
}

// SaveAIConfig 保存AI配置
func SaveAIConfig(config *AIConfig) error {
	configDir := filepath.Dir("config/ai.json")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal AI config: %w", err)
	}

	if err := os.WriteFile("config/ai.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write AI config: %w", err)
	}

	return nil
}

// LoadAIManager 从配置文件加载AI管理器
func LoadAIManager() (*ai.Manager, error) {
	config, err := LoadAIConfig()
	if err != nil {
		return nil, err
	}

	manager := ai.NewManager()

	// 添加所有配置的提供商
	for name, providerConfig := range config.Providers {
		if err := manager.AddClient(name, providerConfig); err != nil {
			return nil, fmt.Errorf("failed to add AI client %s: %w", name, err)
		}
	}

	// 设置默认提供商
	if config.DefaultProvider != "" {
		if err := manager.SetDefault(config.DefaultProvider); err != nil {
			// 如果指定的默认提供商不存在，使用第一个可用的
			clients := manager.ListClients()
			if len(clients) > 0 {
				manager.SetDefault(clients[0])
			}
		}
	}

	return manager, nil
}

// GetAIGlobalOption 获取全局AI选项
func GetAIGlobalOption(key string, defaultValue interface{}) interface{} {
	config, err := LoadAIConfig()
	if err != nil {
		return defaultValue
	}

	if value, exists := config.GlobalOptions[key]; exists {
		return value
	}
	return defaultValue
}

// IsAIFeatureEnabled 检查AI功能是否启用
func IsAIFeatureEnabled(feature string) bool {
	switch feature {
	case "chat":
		return GetAIGlobalOption("chat_enabled", true).(bool)
	case "completion":
		return GetAIGlobalOption("completion_enabled", true).(bool)
	case "embedding":
		return GetAIGlobalOption("embedding_enabled", true).(bool)
	case "streaming":
		return GetAIGlobalOption("streaming_enabled", true).(bool)
	default:
		return false
	}
}
