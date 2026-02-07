package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/config"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/clarkzhu2020/aidecms/pkg/log"
	"github.com/clarkzhu2020/aidecms/pkg/redis"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Application 是AideCMS框架的核心结构
type Application struct {
	Server     *server.Hertz
	Router     *Router
	Config     *config.Config
	DB         *database.Database
	ClickHouse *database.ClickHouseManager
	Redis      *redis.Client
	Logger     *log.Logger
	ConfigPath string
	AppName    string
	AppVersion string
	Env        string
	Debug      bool
	booted     bool
}

// NewApplication 创建一个新的应用实例
func NewApplication() *Application {
	app := &Application{
		AppName:    "AideCMS",
		AppVersion: "1.0.0",
		Env:        "development",
		Debug:      true,
		ConfigPath: "config",
		booted:     false,
	}

	return app
}

// Boot 启动应用程序，初始化各种组件
func (app *Application) Boot() *Application {
	if app.booted {
		return app
	}

	// 加载环境变量
	if err := config.LoadEnv(".env"); err != nil {
		hlog.Warnf("Failed to load .env file: %v", err)
	}

	// Update app configuration from environment
	app.AppName = config.GetEnv("APP_NAME", "AideCMS")
	app.AppVersion = config.GetEnv("APP_VERSION", "1.0.0")
	app.Env = config.GetEnv("APP_ENV", "development")
	app.Debug = config.GetEnvBool("APP_DEBUG", true)

	// 加载配置
	app.loadConfig()

	// 初始化日志
	app.initLogger()

	// 初始化服务器
	app.initServer()

	// 初始化路由
	app.initRouter()

	// 初始化数据库
	app.initDatabase()

	// 初始化ClickHouse
	app.initClickHouse()

	// 初始化Redis
	app.initRedis()

	// 初始化PayPal
	app.initPayPal()

	// 初始化Stripe
	app.initStripe()

	// 初始化MoonPay
	app.initMoonPay()

	// 初始化Coinbase
	app.initCoinbase()

	// 初始化KuCoin
	app.initKuCoin()

	app.booted = true
	return app
}

// loadConfig 加载配置
func (app *Application) loadConfig() {
	app.Config = config.NewConfig([]string{app.ConfigPath})
	if err := app.Config.Load(); err != nil {
		hlog.Warnf("Failed to load config: %v", err)
	}
}

// initLogger 初始化日志
func (app *Application) initLogger() {
	level := hlog.LevelInfo
	if app.Debug {
		level = hlog.LevelDebug
	}
	hlog.SetLevel(level)
}

// initServer 初始化服务器
func (app *Application) initServer() {
	app.Server = server.Default()
}

// initRouter 初始化路由
func (app *Application) initRouter() {
	app.Router = NewRouter(app.Server)
}

// SetDebug 设置调试模式
func (app *Application) SetDebug(debug bool) *Application {
	app.Debug = debug
	return app
}

// SetEnv 设置环境
func (app *Application) SetEnv(env string) *Application {
	app.Env = env
	return app
}

// SetConfigPath 设置配置文件路径
func (app *Application) SetConfigPath(path string) *Application {
	app.ConfigPath = path
	return app
}

// Version 获取应用版本
func (app *Application) Version() string {
	return fmt.Sprintf("%s v%s", app.AppName, app.AppVersion)
}

// LoadConfigFile 加载指定的配置文件
func (app *Application) LoadConfigFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// RegisterRoutes 注册路由
func (app *Application) RegisterRoutes(fn func(*Router)) {
	fn(app.Router)
}

// RegisterMiddleware 注册全局中间件
func (app *Application) RegisterMiddleware(handlers ...app.HandlerFunc) {
	app.Server.Use(handlers...)
}

// Static 注册静态文件目录
func (app *Application) Static(path, root string) {
	app.Router.Static(path, root)
}

// StaticFile 注册静态文件
func (app *Application) StaticFile(path, filepath string) {
	app.Router.StaticFile(path, filepath)
}

// GetPublicPath 获取公共目录路径
func (app *Application) GetPublicPath() string {
	return filepath.Join(app.GetBasePath(), "public")
}

// GetStoragePath 获取存储目录路径
func (app *Application) GetStoragePath() string {
	return filepath.Join(app.GetBasePath(), "storage")
}

// GetBasePath 获取应用基础路径
func (app *Application) GetBasePath() string {
	dir, _ := os.Getwd()
	return dir
}

// initDatabase 初始化数据库
func (app *Application) initDatabase() {
	// 从环境变量获取数据库配置
	dbConfig := &database.Config{
		Driver:      config.GetEnv("DB_CONNECTION", "postgres"),
		Host:        config.GetEnv("DB_HOST", "127.0.0.1"),
		Port:        config.GetEnv("DB_PORT", "5432"),
		Database:    config.GetEnv("DB_DATABASE", "aidecms"),
		Username:    config.GetEnv("DB_USERNAME", "aidecms"),
		Password:    config.GetEnv("DB_PASSWORD", ""),
		Charset:     config.GetEnv("DB_CHARSET", "utf8mb4"),
		MaxIdleConn: config.GetEnvInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConn: config.GetEnvInt("DB_MAX_OPEN_CONNS", 100),
		MaxLifetime: time.Duration(config.GetEnvInt("DB_MAX_LIFETIME", 3600)) * time.Second,
		Debug:       app.Debug,
	}

	app.DB = database.NewDatabase(dbConfig)

	// 增加重试机制，解决 Docker 环境下数据库启动慢的问题
	maxRetries := 5
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = app.DB.Connect(); err == nil {
			hlog.Infof("Successfully connected to database on attempt %d", i+1)
			return
		}
		hlog.Warnf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	hlog.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
}

// Run 运行应用程序
func (app *Application) Run() {
	if !app.booted {
		app.Boot()
	}

	// 启动服务器
	go func() {
		host := app.Config.GetString("server.host", "0.0.0.0")
		port := app.Config.GetInt("server.port", 8888)
		addr := fmt.Sprintf("%s:%d", host, port)

		// 打印所有注册的路由
		if app.Debug {
			app.Router.PrintRoutes()
		}

		hlog.Infof("Server is running on %s", addr)
		err := app.Server.Run()
		if err != nil {
			hlog.Fatal("Server run error: ", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	hlog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		hlog.Fatal("Server forced to shutdown: ", err)
	}

	// 关闭数据库连接
	if app.DB != nil {
		if err := app.DB.Close(); err != nil {
			hlog.Errorf("Failed to close database connection: %v", err)
		}
	}

	// 关闭ClickHouse连接
	if app.ClickHouse != nil {
		if err := app.ClickHouse.Close(); err != nil {
			hlog.Errorf("Failed to close ClickHouse connection: %v", err)
		}
	}

	hlog.Info("Server exiting")
}



// initClickHouse 初始化ClickHouse连接
func (app *Application) initClickHouse() {
	if !config.GetEnvBool("CLICKHOUSE_ENABLED", false) {
		return
	}

	clickhouseConfig := &database.ClickHouseConfig{
		Addr:        []string{fmt.Sprintf("%s:%s", config.GetEnv("CLICKHOUSE_HOST", "localhost"), config.GetEnv("CLICKHOUSE_PORT", "9000"))},
		Database:    config.GetEnv("CLICKHOUSE_DATABASE", "default"),
		Username:    config.GetEnv("CLICKHOUSE_USERNAME", "default"),
		Password:    config.GetEnv("CLICKHOUSE_PASSWORD", ""),
		Debug:       app.Debug,
		MaxOpenConn: config.GetEnvInt("CLICKHOUSE_MAX_OPEN_CONNS", 10),
		MaxIdleConn: config.GetEnvInt("CLICKHOUSE_MAX_IDLE_CONNS", 5),
		MaxLifetime: time.Duration(config.GetEnvInt("CLICKHOUSE_MAX_LIFETIME", 3600)) * time.Second,
	}

	app.ClickHouse = database.NewClickHouseManager(clickhouseConfig)
	if err := app.ClickHouse.Connect(); err != nil {
		hlog.Warnf("Failed to connect to ClickHouse: %v", err)
		return
	}

	// 初始化表
	if err := app.ClickHouse.InitTable(); err != nil {
		hlog.Warnf("Failed to initialize ClickHouse tables: %v", err)
	}
}



// initRedis 初始化Redis连接
func (app *Application) initRedis() {
	if !config.GetEnvBool("REDIS_ENABLED", false) {
		return
	}

	redisConfig := &redis.Config{
		Host:     config.GetEnv("REDIS_HOST", "localhost"),
		Port:     config.GetEnv("REDIS_PORT", "6379"),
		Password: config.GetEnv("REDIS_PASSWORD", ""),
		DB:       config.GetEnvInt("REDIS_DB", 0),
	}

	app.Redis = redis.NewClient(redisConfig)
	if err := app.Redis.Connect(); err != nil {
		hlog.Warnf("Failed to connect to Redis: %v", err)
	}
}

// initPayPal 初始化PayPal支付服务
func (app *Application) initPayPal() {
	if err := config.InitPayPal(); err != nil {
		hlog.Warnf("Failed to initialize PayPal: %v", err)
	}
}

// initStripe 初始化Stripe支付服务
func (app *Application) initStripe() {
	if err := config.InitStripe(); err != nil {
		hlog.Warnf("Failed to initialize Stripe: %v", err)
	}
}

// initMoonPay 初始化MoonPay支付服务
func (app *Application) initMoonPay() {
	config.InitMoonPay()
}

// initCoinbase 初始化Coinbase支付和交易服务
func (app *Application) initCoinbase() {
	config.InitCoinbase()
}

// initKuCoin 初始化KuCoin交易服务
func (app *Application) initKuCoin() {
	if err := config.InitKuCoin(); err != nil {
		hlog.Warnf("Failed to initialize KuCoin: %v", err)
	}
}
