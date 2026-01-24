package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSimpleChecker(t *testing.T) {
	// Test healthy checker
	healthyChecker := NewSimpleChecker("test", func(ctx context.Context) error {
		return nil
	})

	result := healthyChecker.Check(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("Expected StatusHealthy, got %s", result.Status)
	}
	if result.Name != "test" {
		t.Errorf("Expected name 'test', got %s", result.Name)
	}

	// Test unhealthy checker
	unhealthyChecker := NewSimpleChecker("test2", func(ctx context.Context) error {
		return errors.New("service unavailable")
	})

	result2 := unhealthyChecker.Check(context.Background())
	if result2.Status != StatusUnhealthy {
		t.Errorf("Expected StatusUnhealthy, got %s", result2.Status)
	}
	if result2.Error != "service unavailable" {
		t.Errorf("Expected error 'service unavailable', got %s", result2.Error)
	}
}

func TestSimpleCheckerWithDetails(t *testing.T) {
	checker := NewSimpleChecker("test", func(ctx context.Context) error {
		return nil
	}).WithDetails(func(ctx context.Context) map[string]interface{} {
		return map[string]interface{}{
			"version": "1.0.0",
			"uptime":  3600,
		}
	})

	result := checker.Check(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("Expected StatusHealthy, got %s", result.Status)
	}
	if result.Details == nil {
		t.Error("Expected details, got nil")
	}
	if result.Details["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %v", result.Details["version"])
	}
}

func TestDegradableChecker(t *testing.T) {
	// Test healthy checker
	fastChecker := NewDegradableChecker("test", func(ctx context.Context) error {
		return nil
	}, 100*time.Millisecond)

	result := fastChecker.Check(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("Expected StatusHealthy, got %s", result.Status)
	}

	// Test degraded checker (slow response)
	slowChecker := NewDegradableChecker("test2", func(ctx context.Context) error {
		time.Sleep(150 * time.Millisecond)
		return nil
	}, 100*time.Millisecond)

	result2 := slowChecker.Check(context.Background())
	if result2.Status != StatusDegraded {
		t.Errorf("Expected StatusDegraded, got %s", result2.Status)
	}
	if result2.Duration < 100*time.Millisecond {
		t.Errorf("Expected duration >= 100ms, got %v", result2.Duration)
	}

	// Test unhealthy checker
	failedChecker := NewDegradableChecker("test3", func(ctx context.Context) error {
		return errors.New("connection failed")
	}, 100*time.Millisecond)

	result3 := failedChecker.Check(context.Background())
	if result3.Status != StatusUnhealthy {
		t.Errorf("Expected StatusUnhealthy, got %s", result3.Status)
	}
}

func TestHealthChecker_Register(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)

	checker1 := NewSimpleChecker("check1", func(ctx context.Context) error {
		return nil
	})

	checker2 := NewSimpleChecker("check2", func(ctx context.Context) error {
		return nil
	})

	hc.Register(checker1)
	hc.Register(checker2)

	results := hc.Check(context.Background())
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if _, exists := results["check1"]; !exists {
		t.Error("Expected check1 in results")
	}
	if _, exists := results["check2"]; !exists {
		t.Error("Expected check2 in results")
	}
}

func TestHealthChecker_Check(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)

	hc.Register(NewSimpleChecker("healthy", func(ctx context.Context) error {
		return nil
	}))

	hc.Register(NewSimpleChecker("unhealthy", func(ctx context.Context) error {
		return errors.New("failed")
	}))

	results := hc.Check(context.Background())

	if results["healthy"].Status != StatusHealthy {
		t.Errorf("Expected healthy status, got %s", results["healthy"].Status)
	}

	if results["unhealthy"].Status != StatusUnhealthy {
		t.Errorf("Expected unhealthy status, got %s", results["unhealthy"].Status)
	}
}

func TestHealthChecker_GetStatus(t *testing.T) {
	tests := []struct {
		name     string
		checkers []Checker
		expected Status
	}{
		{
			name: "all healthy",
			checkers: []Checker{
				NewSimpleChecker("c1", func(ctx context.Context) error { return nil }),
				NewSimpleChecker("c2", func(ctx context.Context) error { return nil }),
			},
			expected: StatusHealthy,
		},
		{
			name: "one unhealthy",
			checkers: []Checker{
				NewSimpleChecker("c1", func(ctx context.Context) error { return nil }),
				NewSimpleChecker("c2", func(ctx context.Context) error { return errors.New("error") }),
			},
			expected: StatusUnhealthy,
		},
		{
			name: "one degraded",
			checkers: []Checker{
				NewSimpleChecker("c1", func(ctx context.Context) error { return nil }),
				NewDegradableChecker("c2", func(ctx context.Context) error {
					time.Sleep(150 * time.Millisecond)
					return nil
				}, 100*time.Millisecond),
			},
			expected: StatusDegraded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := NewHealthChecker(5 * time.Second)
			for _, checker := range tt.checkers {
				hc.Register(checker)
			}

			status := hc.GetStatus(context.Background())
			if status != tt.expected {
				t.Errorf("Expected status %s, got %s", tt.expected, status)
			}
		})
	}
}

func TestHealthChecker_GetSummary(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)

	hc.Register(NewSimpleChecker("c1", func(ctx context.Context) error { return nil }))
	hc.Register(NewSimpleChecker("c2", func(ctx context.Context) error { return nil }))
	hc.Register(NewSimpleChecker("c3", func(ctx context.Context) error { return errors.New("error") }))

	summary := hc.GetSummary(context.Background())

	if summary["status"] != StatusUnhealthy {
		t.Errorf("Expected status unhealthy, got %v", summary["status"])
	}
	if summary["total_checks"] != 3 {
		t.Errorf("Expected 3 total checks, got %v", summary["total_checks"])
	}
	if summary["healthy_count"] != 2 {
		t.Errorf("Expected 2 healthy checks, got %v", summary["healthy_count"])
	}
	if summary["unhealthy_count"] != 1 {
		t.Errorf("Expected 1 unhealthy check, got %v", summary["unhealthy_count"])
	}
}

func TestHealthChecker_CheckOne(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)

	hc.Register(NewSimpleChecker("check1", func(ctx context.Context) error {
		return nil
	}))

	hc.Register(NewSimpleChecker("check2", func(ctx context.Context) error {
		return errors.New("failed")
	}))

	// Check existing checker
	result, err := hc.CheckOne(context.Background(), "check1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Status != StatusHealthy {
		t.Errorf("Expected healthy status, got %s", result.Status)
	}

	// Check non-existing checker
	_, err = hc.CheckOne(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent checker")
	}
}

func TestHealthChecker_Cache(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)
	hc.SetCacheTTL(500 * time.Millisecond)

	callCount := 0
	hc.Register(NewSimpleChecker("cached", func(ctx context.Context) error {
		callCount++
		return nil
	}))

	// First call - should execute
	hc.Check(context.Background())
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second call - should use cache
	hc.Check(context.Background())
	if callCount != 1 {
		t.Errorf("Expected still 1 call (cached), got %d", callCount)
	}

	// Wait for cache expiration
	time.Sleep(600 * time.Millisecond)

	// Third call - should execute again
	hc.Check(context.Background())
	if callCount != 2 {
		t.Errorf("Expected 2 calls after cache expiry, got %d", callCount)
	}
}

func TestHealthChecker_ClearCache(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)
	hc.SetCacheTTL(10 * time.Second) // Long TTL

	callCount := 0
	hc.Register(NewSimpleChecker("test", func(ctx context.Context) error {
		callCount++
		return nil
	}))

	// First call
	hc.Check(context.Background())
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second call (should be cached)
	hc.Check(context.Background())
	if callCount != 1 {
		t.Errorf("Expected still 1 call (cached), got %d", callCount)
	}

	// Clear cache
	hc.ClearCache()

	// Third call (should execute)
	hc.Check(context.Background())
	if callCount != 2 {
		t.Errorf("Expected 2 calls after cache clear, got %d", callCount)
	}
}

func TestHealthChecker_Timeout(t *testing.T) {
	hc := NewHealthChecker(100 * time.Millisecond) // Short timeout

	hc.Register(NewSimpleChecker("slow", func(ctx context.Context) error {
		select {
		case <-time.After(500 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}))

	results := hc.Check(context.Background())
	result := results["slow"]

	// Should fail due to timeout
	if result.Status == StatusHealthy {
		t.Error("Expected unhealthy status due to timeout")
	}
}

func TestConcurrentHealthChecks(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)

	// Register multiple checkers
	for i := 0; i < 10; i++ {
		name := string(rune('a' + i))
		hc.Register(NewSimpleChecker(name, func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}))
	}

	start := time.Now()
	results := hc.Check(context.Background())
	duration := time.Since(start)

	// All checks should run concurrently
	// Total time should be ~10ms (not 100ms if sequential)
	if duration > 100*time.Millisecond {
		t.Errorf("Checks took too long: %v (expected ~10ms due to concurrency)", duration)
	}

	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}
}
