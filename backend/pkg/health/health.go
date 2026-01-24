package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Status 健康状态
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

// CheckResult 健康检查结果
type CheckResult struct {
	Name      string                 `json:"name"`
	Status    Status                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// Checker 健康检查接口
type Checker interface {
	// Name 返回检查器名称
	Name() string
	// Check 执行健康检查
	Check(ctx context.Context) CheckResult
}

// HealthChecker 健康检查器管理
type HealthChecker struct {
	checkers []Checker
	mu       sync.RWMutex
	timeout  time.Duration
	cache    map[string]*cachedResult
	cacheTTL time.Duration
}

type cachedResult struct {
	result    CheckResult
	expiresAt time.Time
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		checkers: make([]Checker, 0),
		timeout:  timeout,
		cache:    make(map[string]*cachedResult),
		cacheTTL: 10 * time.Second,
	}
}

// Register 注册健康检查
func (h *HealthChecker) Register(checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers = append(h.checkers, checker)
}

// Check 执行所有健康检查
func (h *HealthChecker) Check(ctx context.Context) map[string]CheckResult {
	h.mu.RLock()
	checkers := make([]Checker, len(h.checkers))
	copy(checkers, h.checkers)
	h.mu.RUnlock()

	results := make(map[string]CheckResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, checker := range checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()

			// Check cache
			if cached := h.getCached(c.Name()); cached != nil {
				mu.Lock()
				results[c.Name()] = *cached
				mu.Unlock()
				return
			}

			// Execute check with timeout
			checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
			defer cancel()

			result := c.Check(checkCtx)

			// Cache result
			h.setCached(c.Name(), result)

			mu.Lock()
			results[c.Name()] = result
			mu.Unlock()
		}(checker)
	}

	wg.Wait()
	return results
}

// CheckOne 执行单个健康检查
func (h *HealthChecker) CheckOne(ctx context.Context, name string) (CheckResult, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, checker := range h.checkers {
		if checker.Name() == name {
			// Check cache
			if cached := h.getCached(name); cached != nil {
				return *cached, nil
			}

			checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
			defer cancel()

			result := checker.Check(checkCtx)
			h.setCached(name, result)
			return result, nil
		}
	}

	return CheckResult{}, fmt.Errorf("checker not found: %s", name)
}

// GetStatus 获取整体健康状态
func (h *HealthChecker) GetStatus(ctx context.Context) Status {
	results := h.Check(ctx)

	healthyCount := 0
	degradedCount := 0
	unhealthyCount := 0

	for _, result := range results {
		switch result.Status {
		case StatusHealthy:
			healthyCount++
		case StatusDegraded:
			degradedCount++
		case StatusUnhealthy:
			unhealthyCount++
		}
	}

	// If any check is unhealthy, overall status is unhealthy
	if unhealthyCount > 0 {
		return StatusUnhealthy
	}

	// If any check is degraded, overall status is degraded
	if degradedCount > 0 {
		return StatusDegraded
	}

	return StatusHealthy
}

// GetSummary 获取健康检查摘要
func (h *HealthChecker) GetSummary(ctx context.Context) map[string]interface{} {
	results := h.Check(ctx)
	status := h.GetStatus(ctx)

	healthyCount := 0
	degradedCount := 0
	unhealthyCount := 0

	for _, result := range results {
		switch result.Status {
		case StatusHealthy:
			healthyCount++
		case StatusDegraded:
			degradedCount++
		case StatusUnhealthy:
			unhealthyCount++
		}
	}

	return map[string]interface{}{
		"status":          status,
		"timestamp":       time.Now(),
		"total_checks":    len(results),
		"healthy_count":   healthyCount,
		"degraded_count":  degradedCount,
		"unhealthy_count": unhealthyCount,
		"checks":          results,
	}
}

// SetCacheTTL 设置缓存TTL
func (h *HealthChecker) SetCacheTTL(ttl time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cacheTTL = ttl
}

// ClearCache 清除缓存
func (h *HealthChecker) ClearCache() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cache = make(map[string]*cachedResult)
}

func (h *HealthChecker) getCached(name string) *CheckResult {
	h.mu.RLock()
	defer h.mu.RUnlock()

	cached, exists := h.cache[name]
	if !exists {
		return nil
	}

	if time.Now().After(cached.expiresAt) {
		return nil
	}

	return &cached.result
}

func (h *HealthChecker) setCached(name string, result CheckResult) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cache[name] = &cachedResult{
		result:    result,
		expiresAt: time.Now().Add(h.cacheTTL),
	}
}

// SimpleChecker 简单检查器
type SimpleChecker struct {
	name     string
	checkFn  func(ctx context.Context) error
	detailFn func(ctx context.Context) map[string]interface{}
}

// NewSimpleChecker 创建简单检查器
func NewSimpleChecker(name string, checkFn func(ctx context.Context) error) *SimpleChecker {
	return &SimpleChecker{
		name:    name,
		checkFn: checkFn,
	}
}

// WithDetails 添加详细信息函数
func (s *SimpleChecker) WithDetails(detailFn func(ctx context.Context) map[string]interface{}) *SimpleChecker {
	s.detailFn = detailFn
	return s
}

// Name 实现 Checker 接口
func (s *SimpleChecker) Name() string {
	return s.name
}

// Check 实现 Checker 接口
func (s *SimpleChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	result := CheckResult{
		Name:      s.name,
		Timestamp: start,
	}

	err := s.checkFn(ctx)
	result.Duration = time.Since(start)

	if err != nil {
		result.Status = StatusUnhealthy
		result.Error = err.Error()
		result.Message = fmt.Sprintf("%s check failed", s.name)
	} else {
		result.Status = StatusHealthy
		result.Message = fmt.Sprintf("%s is healthy", s.name)
	}

	// Add details if function provided
	if s.detailFn != nil {
		result.Details = s.detailFn(ctx)
	}

	return result
}

// DegradableChecker 可降级检查器
type DegradableChecker struct {
	name             string
	checkFn          func(ctx context.Context) error
	detailFn         func(ctx context.Context) map[string]interface{}
	degradeThreshold time.Duration // 响应时间超过此值视为降级
}

// NewDegradableChecker 创建可降级检查器
func NewDegradableChecker(name string, checkFn func(ctx context.Context) error, degradeThreshold time.Duration) *DegradableChecker {
	return &DegradableChecker{
		name:             name,
		checkFn:          checkFn,
		degradeThreshold: degradeThreshold,
	}
}

// WithDetails 添加详细信息函数
func (d *DegradableChecker) WithDetails(detailFn func(ctx context.Context) map[string]interface{}) *DegradableChecker {
	d.detailFn = detailFn
	return d
}

// Name 实现 Checker 接口
func (d *DegradableChecker) Name() string {
	return d.name
}

// Check 实现 Checker 接口
func (d *DegradableChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	result := CheckResult{
		Name:      d.name,
		Timestamp: start,
	}

	err := d.checkFn(ctx)
	result.Duration = time.Since(start)

	if err != nil {
		result.Status = StatusUnhealthy
		result.Error = err.Error()
		result.Message = fmt.Sprintf("%s check failed", d.name)
	} else if result.Duration > d.degradeThreshold {
		result.Status = StatusDegraded
		result.Message = fmt.Sprintf("%s is slow (took %v, threshold %v)", d.name, result.Duration, d.degradeThreshold)
	} else {
		result.Status = StatusHealthy
		result.Message = fmt.Sprintf("%s is healthy", d.name)
	}

	// Add details if function provided
	if d.detailFn != nil {
		result.Details = d.detailFn(ctx)
	}

	return result
}
