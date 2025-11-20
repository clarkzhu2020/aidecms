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

	// 初始化Redis
	app.initRedis()

	app.booted = true
	return app
}

// loadConfig 加载配置
func (app *Application) loadConfig() {
	app.Config = config.NewConfig([]string{app.ConfigPath})
	if err := app.Config.Load(); err != nil {
		hlog.Fatalf("Failed to load config: %v", err)
	}

	// 从配置中更新应用设置
	if appName := app.Config.GetString("app.name"); appName != "" {
		app.AppName = appName
	}

	if env := app.Config.GetString("app.env"); env != "" {
		app.Env = env
	}

	app.Debug = app.Config.GetBool("app.debug", app.Debug)
}

// initLogger 初始化日志
func (app *Application) initLogger() {
	logConfig := &log.Config{
		Level:      app.Config.GetString("log.level", "info"),
		Format:     app.Config.GetString("log.format", "text"),
		Output:     app.Config.GetString("log.output", "stdout"),
		OutputFile: app.Config.GetString("log.output_file", "logs/app.log"),
		TimeFormat: app.Config.GetString("log.time_format", "2006-01-02 15:04:05"),
	}

	app.Logger = log.NewLogger(logConfig)
	if err := app.Logger.Init(); err != nil {
		hlog.Fatalf("Failed to initialize logger: %v", err)
	}
}

// initServer 初始化Hertz服务器
func (app *Application) initServer() {
	host := app.Config.GetString("server.host", "0.0.0.0")
	port := app.Config.GetInt("server.port", 8888)
	addr := fmt.Sprintf("%s:%d", host, port)

	if app.Debug {
		app.Server = server.Default(server.WithHostPorts(addr))
	} else {
		app.Server = server.New(server.WithHostPorts(addr))
	}
}

// initRouter 初始化路由
func (app *Application) initRouter() {
	if app.Server == nil {
		app.initServer()
	}
	app.Router = NewRouter(app.Server)
}

// initDatabase 初始化数据库
func (app *Application) initDatabase() {
	// 强制使用SQLite
	dbConfig := &database.Config{
		Driver:   "sqlite",
		Database: "storage/database/data.db",
		Debug:    app.Debug,
	}

	// 确保数据库目录存在
	if err := os.MkdirAll("storage/database", 0755); err != nil {
		hlog.Fatalf("Failed to create database directory: %v", err)
	}

	app.DB = database.NewDatabase(dbConfig)
	if err := app.DB.Connect(); err != nil {
		hlog.Fatalf("Failed to connect to database: %v", err)
	}
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

	hlog.Info("Server exiting")
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

// initRedis 初始化Redis连接
func (app *Application) initRedis() {
	if !app.Config.GetBool("redis.enabled", false) {
		return
	}

	redisConfig := &redis.Config{
		Host:     app.Config.GetString("redis.host", "localhost"),
		Port:     app.Config.GetString("redis.port", "6379"),
		Password: app.Config.GetString("redis.password", ""),
		DB:       app.Config.GetInt("redis.db", 0),
	}

	app.Redis = redis.NewClient(redisConfig)
	if err := app.Redis.Connect(); err != nil {
		hlog.Warnf("Failed to connect to Redis: %v", err)
	}
}
