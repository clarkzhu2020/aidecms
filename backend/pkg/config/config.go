package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Config 配置管理器
type Config struct {
	items map[string]interface{}
	paths []string
}

// NewConfig 创建一个新的配置管理器
func NewConfig(paths []string) *Config {
	return &Config{
		items: make(map[string]interface{}),
		paths: paths,
	}
}

// Load 加载配置文件
func (c *Config) Load() error {
	for _, path := range c.paths {
		files, err := os.ReadDir(path)
		if err != nil {
			hlog.Warnf("Failed to read config directory %s: %v", path, err)
			continue
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			configName := strings.TrimSuffix(file.Name(), ".json")
			configPath := filepath.Join(path, file.Name())

			data, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to read config file %s: %w", configPath, err)
			}

			var configData interface{}
			if err := json.Unmarshal(data, &configData); err != nil {
				return fmt.Errorf("failed to parse config file %s: %w", configPath, err)
			}

			c.items[configName] = configData
		}
	}

	return nil
}

// Get 获取配置项
func (c *Config) Get(key string, defaultValue ...interface{}) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return nil
	}

	// 获取顶级配置
	configName := keys[0]
	config, ok := c.items[configName]
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	// 如果只有一级，直接返回
	if len(keys) == 1 {
		return config
	}

	// 递归获取嵌套配置
	current := config
	for _, k := range keys[1:] {
		m, ok := current.(map[string]interface{})
		if !ok {
			if len(defaultValue) > 0 {
				return defaultValue[0]
			}
			return nil
		}

		current, ok = m[k]
		if !ok {
			if len(defaultValue) > 0 {
				return defaultValue[0]
			}
			return nil
		}
	}

	return current
}

// GetString 获取字符串配置
func (c *Config) GetString(key string, defaultValue ...string) string {
	value := c.Get(key)
	if value == nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	if str, ok := value.(string); ok {
		return str
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string, defaultValue ...int) int {
	value := c.Get(key)
	if value == nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	value := c.Get(key)
	if value == nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	if b, ok := value.(bool); ok {
		return b
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return false
}

// Set 设置配置项
func (c *Config) Set(key string, value interface{}) {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return
	}

	// 设置顶级配置
	configName := keys[0]
	if len(keys) == 1 {
		c.items[configName] = value
		return
	}

	// 获取或创建顶级配置
	config, ok := c.items[configName]
	if !ok {
		config = make(map[string]interface{})
		c.items[configName] = config
	}

	// 递归设置嵌套配置
	current, ok := config.(map[string]interface{})
	if !ok {
		current = make(map[string]interface{})
		c.items[configName] = current
	}

	for i, k := range keys[1:] {
		if i == len(keys)-2 {
			// 最后一个键，设置值
			current[k] = value
			break
		}

		// 获取或创建子配置
		next, ok := current[k]
		if !ok {
			next = make(map[string]interface{})
			current[k] = next
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			nextMap = make(map[string]interface{})
			current[k] = nextMap
		}

		current = nextMap
	}
}
