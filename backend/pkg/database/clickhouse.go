package database

import (
	"context"

	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// ClickHouseConfig ClickHouse configuration
type ClickHouseConfig struct {
	Addr        []string
	Database    string
	Username    string
	Password    string
	Debug       bool
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration
}

// ClickHouseManager manages ClickHouse connection
type ClickHouseManager struct {
	Conn   driver.Conn
	Config *ClickHouseConfig
}

var globalClickHouse targetGlobalClickHouse

type targetGlobalClickHouse struct {
	conn driver.Conn
}

// SetClickHouse sets the global ClickHouse connection
func SetClickHouse(conn driver.Conn) {
	globalClickHouse.conn = conn
}

// GetClickHouse gets the global ClickHouse connection
func GetClickHouse() driver.Conn {
	return globalClickHouse.conn
}

// NewClickHouseManager creates a new ClickHouse manager
func NewClickHouseManager(config *ClickHouseConfig) *ClickHouseManager {
	return &ClickHouseManager{
		Config: config,
	}
}

// Connect connects to ClickHouse
func (m *ClickHouseManager) Connect() error {
	var opts clickhouse.Options
	opts.Addr = m.Config.Addr
	opts.Auth = clickhouse.Auth{
		Database: m.Config.Database,
		Username: m.Config.Username,
		Password: m.Config.Password,
	}
	opts.Debug = m.Config.Debug
	opts.MaxOpenConns = m.Config.MaxOpenConn
	opts.MaxIdleConns = m.Config.MaxIdleConn
	opts.ConnMaxLifetime = m.Config.MaxLifetime
	opts.Compression = &clickhouse.Compression{
		Method: clickhouse.CompressionLZ4,
	}

	// TLS config if needed (not implemented for local dev but good to have placeholder)
	// opts.TLS = &tls.Config{InsecureSkipVerify: true}

	conn, err := clickhouse.Open(&opts)
	if err != nil {
		return err
	}

	if err := conn.Ping(context.Background()); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return fmt.Errorf("clickhouse exception: %d %s %s", exception.Code, exception.Message, exception.StackTrace)
		}
		return err
	}

	hlog.Info("Successfully connected to ClickHouse")
	m.Conn = conn
	SetClickHouse(conn)
	
	return nil
}

// InitTable initializes the necessary tables in ClickHouse
func (m *ClickHouseManager) InitTable() error {
	if m.Conn == nil {
		return fmt.Errorf("clickhouse connection is nil")
	}

	query := `
		CREATE TABLE IF NOT EXISTS market_data (
			exchange String,
			pair     String,
			price    Float64,
			timestamp DateTime64(3, 'UTC')
		) ENGINE = MergeTree()
		ORDER BY (pair, timestamp)
	`
	if err := m.Conn.Exec(context.Background(), query); err != nil {
		return fmt.Errorf("failed to create market_data table: %v", err)
	}

	hlog.Info("Market data table initialized successfully")
	return nil
}

// Close closes the connection
func (m *ClickHouseManager) Close() error {
	if m.Conn != nil {
		return m.Conn.Close()
	}
	return nil
}

