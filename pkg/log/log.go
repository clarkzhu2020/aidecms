package log

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Config 日志配置
type Config struct {
	Level      string
	Format     string
	Output     string
	OutputFile string
	TimeFormat string
}

// Logger 日志管理器
type Logger struct {
	Config *Config
}

// NewLogger 创建一个新的日志管理器
func NewLogger(config *Config) *Logger {
	return &Logger{
		Config: config,
	}
}

// Init 初始化日志
func (l *Logger) Init() error {
	// 设置日志级别
	level := hlog.LevelInfo
	switch l.Config.Level {
	case "debug":
		level = hlog.LevelDebug
	case "info":
		level = hlog.LevelInfo
	case "warn":
		level = hlog.LevelWarn
	case "error":
		level = hlog.LevelError
	case "fatal":
		level = hlog.LevelFatal
	}
	hlog.SetLevel(level)

	// 设置输出
	var output io.Writer
	switch l.Config.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "file":
		if l.Config.OutputFile == "" {
			l.Config.OutputFile = "logs/app.log"
		}

		// 创建日志目录
		dir := filepath.Dir(l.Config.OutputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// 打开日志文件
		file, err := os.OpenFile(l.Config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		output = file
	default:
		output = os.Stdout
	}
	hlog.SetOutput(output)

    // 如果配置为 JSON 格式，这里可以扩展使用 zap 或 slog
    // 目前 hlog 默认是 text 格式，这里保留扩展点
    // 实际生产环境建议使用 hertz-contrib/logger/zap
    
	return nil
}


// Debug 调试日志
func (l *Logger) Debug(args ...interface{}) {
	hlog.Debug(args...)
}

// Debugf 格式化调试日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	hlog.Debugf(format, args...)
}

// Info 信息日志
func (l *Logger) Info(args ...interface{}) {
	hlog.Info(args...)
}

// Infof 格式化信息日志
func (l *Logger) Infof(format string, args ...interface{}) {
	hlog.Infof(format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(args ...interface{}) {
	hlog.Warn(args...)
}

// Warnf 格式化警告日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	hlog.Warnf(format, args...)
}

// Error 错误日志
func (l *Logger) Error(args ...interface{}) {
	hlog.Error(args...)
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	hlog.Errorf(format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(args ...interface{}) {
	hlog.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func (l *Logger) Fatalf(format string, args ...interface{}) {
	hlog.Fatalf(format, args...)
}
