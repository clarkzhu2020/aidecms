package css

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/ai"
	"github.com/clarkzhu2020/aidecms/pkg/css/conversation"
	"github.com/clarkzhu2020/aidecms/pkg/css/rag"
	"github.com/clarkzhu2020/aidecms/pkg/css/kb"
	"gorm.io/gorm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// InitCustomerServiceSystem 初始化客服系统
func InitCustomerServiceSystem(db *gorm.DB, aiManager *ai.Manager) (*CSSEngine, *WSManager) {
	// 加载配置
	config := loadConfig()

	hlog.Infof("[CSS] Initializing customer service system...")
	hlog.Infof("[CSS] Config: model=%s, top_k=%d, confidence=%.2f",
		config.DefaultModel, config.TopK, config.ConfidenceThreshold)

	// 1. 初始化对话管理器
	convManager := conversation.NewManager(db)
	hlog.Infof("[CSS] Conversation manager initialized")

	// 2. 初始化RAG检索器
	var ragRetriever *rag.Retriever
	if config.ConfidenceThreshold > 0 {
		vectorStore := rag.NewVectorStore(db)
		embedder := rag.NewEmbedder(aiManager)
		ragRetriever = rag.NewRetriever(vectorStore, embedder)
		hlog.Infof("[CSS] RAG retriever initialized")
	}

	// 3. 初始化知识库服务
	kbService := kb.NewService(db)
	hlog.Infof("[CSS] Knowledge base service initialized")

	// 4. 初始化核心引擎
	engine := NewCSSEngine(aiManager, ragRetriever, convManager, config)
	hlog.Infof("[CSS] Core engine initialized")

	// 5. 初始化WebSocket管理器
	wsManager := NewWSManager(engine)
	go wsManager.Start()
	hlog.Infof("[CSS] WebSocket manager started")

	hlog.Infof("[CSS] Customer service system initialized successfully")

	return engine, wsManager
}

// 配置辅助函数
func getConfig(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getConfigInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getConfigFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getConfigBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getConfigList(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// loadConfig 从环境变量加载配置
func loadConfig() *Config {
	return &Config{
		// AI 配置
		DefaultModel:     getConfig("CSS_DEFAULT_MODEL", "qianwen"),
		Temperature:     getConfigFloat("CSS_TEMPERATURE", 0.7),
		MaxTokens:      getConfigInt("CSS_MAX_TOKENS", 1000),

		// RAG 配置
		TopK:              getConfigInt("CSS_TOP_K", 5),
		ConfidenceThreshold: getConfigFloat("CSS_CONFIDENCE_THRESHOLD", 0.6),

		// 转接配置
		TransferOnLowConfidence: getConfigBool("CSS_TRANSFER_ON_LOW_CONFIDENCE", true),
		MaxRetries:             getConfigInt("CSS_MAX_RETRIES", 3),
		TransferKeywords:        getConfigList("CSS_TRANSFER_KEYWORDS", []string{"投诉", "退款", "人工", "转接"}),

		// 会话配置
		SessionTimeout: time.Duration(getConfigInt("CSS_SESSION_TIMEOUT", 1800)) * time.Second,
		MaxHistory:     getConfigInt("CSS_MAX_HISTORY", 10),
	}
}
