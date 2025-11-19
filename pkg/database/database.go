package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	Driver      string
	Host        string
	Port        string
	Database    string
	Username    string
	Password    string
	Charset     string
	Prefix      string
	MaxIdleConn int
	MaxOpenConn int
	MaxLifetime time.Duration
	Debug       bool
}

// Database 数据库连接管理器
type Database struct {
	DB     *gorm.DB
	Config *Config
}

// 全局数据库实例
var globalDB *gorm.DB

// NewDatabase 创建一个新的数据库连接管理器
func NewDatabase(config *Config) *Database {
	return &Database{
		Config: config,
	}
}

// SetDB 设置全局数据库实例
func SetDB(db *gorm.DB) {
	globalDB = db
}

// GetDB 获取全局数据库实例
func GetDB() *gorm.DB {
	return globalDB
}

// Connect 连接数据库
func (d *Database) Connect() error {
	var err error
	var db *gorm.DB

	switch d.Config.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			d.Config.Username,
			d.Config.Password,
			d.Config.Host,
			d.Config.Port,
			d.Config.Database,
			d.Config.Charset,
		)
		db, err = gorm.Open(mysql.Open(dsn), d.getGormConfig())
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			d.Config.Host,
			d.Config.Username,
			d.Config.Password,
			d.Config.Database,
			d.Config.Port,
		)
		db, err = gorm.Open(postgres.Open(dsn), d.getGormConfig())
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(d.Config.Database), d.getGormConfig())
	default:
		return fmt.Errorf("unsupported database driver: %s", d.Config.Driver)
	}

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(d.Config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(d.Config.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(d.Config.MaxLifetime)

	d.DB = db
	return nil
}

// getGormConfig 获取GORM配置
func (d *Database) getGormConfig() *gorm.Config {
	config := &gorm.Config{}

	if d.Config.Debug {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	return config
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移模型
func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}

// Transaction 事务
func (d *Database) Transaction(fn func(tx *gorm.DB) error) error {
	return d.DB.Transaction(fn)
}
