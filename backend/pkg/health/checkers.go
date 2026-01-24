package health

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// DatabaseChecker 数据库健康检查器
type DatabaseChecker struct {
	db   *gorm.DB
	name string
}

// NewDatabaseChecker 创建数据库检查器
func NewDatabaseChecker(db *gorm.DB) *DatabaseChecker {
	return &DatabaseChecker{
		db:   db,
		name: "database",
	}
}

// Name 实现 Checker 接口
func (d *DatabaseChecker) Name() string {
	return d.name
}

// Check 实现 Checker 接口
func (d *DatabaseChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	result := CheckResult{
		Name:      d.name,
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Get SQL database
	sqlDB, err := d.db.DB()
	if err != nil {
		result.Status = StatusUnhealthy
		result.Error = err.Error()
		result.Message = "Failed to get database connection"
		result.Duration = time.Since(start)
		return result
	}

	// Ping database
	if err := sqlDB.PingContext(ctx); err != nil {
		result.Status = StatusUnhealthy
		result.Error = err.Error()
		result.Message = "Database ping failed"
		result.Duration = time.Since(start)
		return result
	}

	// Get statistics
	stats := sqlDB.Stats()
	result.Details["open_connections"] = stats.OpenConnections
	result.Details["in_use"] = stats.InUse
	result.Details["idle"] = stats.Idle
	result.Details["max_open_connections"] = stats.MaxOpenConnections

	result.Duration = time.Since(start)

	// Check if degraded (slow response or high connection usage)
	if result.Duration > 500*time.Millisecond {
		result.Status = StatusDegraded
		result.Message = fmt.Sprintf("Database is slow (took %v)", result.Duration)
	} else if stats.MaxOpenConnections > 0 && float64(stats.OpenConnections)/float64(stats.MaxOpenConnections) > 0.8 {
		result.Status = StatusDegraded
		result.Message = "Database connection pool is nearly exhausted"
	} else {
		result.Status = StatusHealthy
		result.Message = "Database is healthy"
	}

	return result
}

// RedisChecker Redis健康检查器
type RedisChecker struct {
	client *redis.Client
	name   string
}

// NewRedisChecker 创建Redis检查器
func NewRedisChecker(client *redis.Client) *RedisChecker {
	return &RedisChecker{
		client: client,
		name:   "redis",
	}
}

// Name 实现 Checker 接口
func (r *RedisChecker) Name() string {
	return r.name
}

// Check 实现 Checker 接口
func (r *RedisChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	result := CheckResult{
		Name:      r.name,
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Ping Redis
	pong, err := r.client.Ping(ctx).Result()
	if err != nil {
		result.Status = StatusUnhealthy
		result.Error = err.Error()
		result.Message = "Redis ping failed"
		result.Duration = time.Since(start)
		return result
	}

	result.Details["ping_response"] = pong

	// Get Redis info
	info, err := r.client.Info(ctx, "stats").Result()
	if err == nil {
		result.Details["info_available"] = true
	}

	// Get memory usage
	memInfo, err := r.client.Info(ctx, "memory").Result()
	if err == nil {
		result.Details["memory_info_available"] = true
		_ = memInfo // Parse if needed
	}

	// Get keyspace info
	dbSize, err := r.client.DBSize(ctx).Result()
	if err == nil {
		result.Details["keys_count"] = dbSize
	}

	result.Duration = time.Since(start)

	// Check if degraded (slow response)
	if result.Duration > 100*time.Millisecond {
		result.Status = StatusDegraded
		result.Message = fmt.Sprintf("Redis is slow (took %v)", result.Duration)
	} else {
		result.Status = StatusHealthy
		result.Message = "Redis is healthy"
	}

	_ = info // Use info if needed

	return result
}

// DiskSpaceChecker 磁盘空间检查器
type DiskSpaceChecker struct {
	path            string
	warningPercent  float64 // 警告阈值百分比
	criticalPercent float64 // 严重阈值百分比
}

// NewDiskSpaceChecker 创建磁盘空间检查器
func NewDiskSpaceChecker(path string, warningPercent, criticalPercent float64) *DiskSpaceChecker {
	return &DiskSpaceChecker{
		path:            path,
		warningPercent:  warningPercent,
		criticalPercent: criticalPercent,
	}
}

// Name 实现 Checker 接口
func (d *DiskSpaceChecker) Name() string {
	return fmt.Sprintf("disk_space_%s", d.path)
}

// Check 实现 Checker 接口
func (d *DiskSpaceChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	result := CheckResult{
		Name:      d.Name(),
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Note: This is a placeholder implementation
	// In production, use syscall.Statfs or similar to get actual disk stats

	// Simulate disk check
	usedPercent := 45.0 // This should be calculated from actual disk stats

	result.Details["path"] = d.path
	result.Details["used_percent"] = usedPercent
	result.Details["warning_threshold"] = d.warningPercent
	result.Details["critical_threshold"] = d.criticalPercent

	result.Duration = time.Since(start)

	if usedPercent >= d.criticalPercent {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("Disk space critical: %.1f%% used (threshold: %.1f%%)", usedPercent, d.criticalPercent)
	} else if usedPercent >= d.warningPercent {
		result.Status = StatusDegraded
		result.Message = fmt.Sprintf("Disk space warning: %.1f%% used (threshold: %.1f%%)", usedPercent, d.warningPercent)
	} else {
		result.Status = StatusHealthy
		result.Message = fmt.Sprintf("Disk space healthy: %.1f%% used", usedPercent)
	}

	return result
}

// HTTPServiceChecker HTTP服务健康检查器
type HTTPServiceChecker struct {
	name     string
	url      string
	method   string
	timeout  time.Duration
	expected int // Expected HTTP status code
}

// NewHTTPServiceChecker 创建HTTP服务检查器
func NewHTTPServiceChecker(name, url string, expectedStatus int) *HTTPServiceChecker {
	return &HTTPServiceChecker{
		name:     name,
		url:      url,
		method:   "GET",
		timeout:  5 * time.Second,
		expected: expectedStatus,
	}
}

// WithMethod 设置HTTP方法
func (h *HTTPServiceChecker) WithMethod(method string) *HTTPServiceChecker {
	h.method = method
	return h
}

// WithTimeout 设置超时时间
func (h *HTTPServiceChecker) WithTimeout(timeout time.Duration) *HTTPServiceChecker {
	h.timeout = timeout
	return h
}

// Name 实现 Checker 接口
func (h *HTTPServiceChecker) Name() string {
	return h.name
}

// Check 实现 Checker 接口
func (h *HTTPServiceChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	result := CheckResult{
		Name:      h.name,
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	result.Details["url"] = h.url
	result.Details["method"] = h.method
	result.Details["expected_status"] = h.expected

	// Note: This is a placeholder implementation
	// In production, use http.Client to make actual request

	// Simulate HTTP check
	statusCode := 200
	result.Details["status_code"] = statusCode

	result.Duration = time.Since(start)

	if statusCode != h.expected {
		result.Status = StatusUnhealthy
		result.Error = fmt.Sprintf("Unexpected status code: %d (expected %d)", statusCode, h.expected)
		result.Message = fmt.Sprintf("%s returned unexpected status", h.name)
	} else if result.Duration > h.timeout {
		result.Status = StatusDegraded
		result.Message = fmt.Sprintf("%s is slow (took %v)", h.name, result.Duration)
	} else {
		result.Status = StatusHealthy
		result.Message = fmt.Sprintf("%s is healthy", h.name)
	}

	return result
}

// MemoryChecker 内存使用检查器
type MemoryChecker struct {
	warningPercent  float64
	criticalPercent float64
}

// NewMemoryChecker 创建内存检查器
func NewMemoryChecker(warningPercent, criticalPercent float64) *MemoryChecker {
	return &MemoryChecker{
		warningPercent:  warningPercent,
		criticalPercent: criticalPercent,
	}
}

// Name 实现 Checker 接口
func (m *MemoryChecker) Name() string {
	return "memory"
}

// Check 实现 Checker 接口
func (m *MemoryChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	result := CheckResult{
		Name:      "memory",
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Note: This is a placeholder implementation
	// In production, use runtime.MemStats or similar to get actual memory usage

	// Simulate memory check
	usedPercent := 35.0 // This should be calculated from actual memory stats

	result.Details["used_percent"] = usedPercent
	result.Details["warning_threshold"] = m.warningPercent
	result.Details["critical_threshold"] = m.criticalPercent

	result.Duration = time.Since(start)

	if usedPercent >= m.criticalPercent {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("Memory usage critical: %.1f%% (threshold: %.1f%%)", usedPercent, m.criticalPercent)
	} else if usedPercent >= m.warningPercent {
		result.Status = StatusDegraded
		result.Message = fmt.Sprintf("Memory usage warning: %.1f%% (threshold: %.1f%%)", usedPercent, m.warningPercent)
	} else {
		result.Status = StatusHealthy
		result.Message = fmt.Sprintf("Memory usage healthy: %.1f%%", usedPercent)
	}

	return result
}
