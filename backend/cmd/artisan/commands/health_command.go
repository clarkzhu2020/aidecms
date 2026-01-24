package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/health"
)

// HealthCommand demonstrates health check functionality
func HealthCommand(args []string) {
	if len(args) > 0 && args[0] == "demo" {
		runHealthDemo()
		return
	}

	fmt.Println("Health Check Commands:")
	fmt.Println("  artisan health demo - Run health check demonstration")
}

func runHealthDemo() {
	fmt.Println("=== Health Check Demo ===\n")

	// Create health checker
	hc := health.NewHealthChecker(5 * time.Second)
	ctx := context.Background()

	// Demo 1: Simple healthy service
	fmt.Println("1. Simple Health Check")
	fmt.Println("   - Checking a healthy service\n")

	healthyService := health.NewSimpleChecker("api_service", func(ctx context.Context) error {
		// Simulate service check
		time.Sleep(10 * time.Millisecond)
		return nil
	}).WithDetails(func(ctx context.Context) map[string]interface{} {
		return map[string]interface{}{
			"version": "1.0.0",
			"uptime":  3600,
		}
	})

	hc.Register(healthyService)

	result := healthyService.Check(ctx)
	fmt.Printf("   Service: %s\n", result.Name)
	fmt.Printf("   Status:  %s\n", result.Status)
	fmt.Printf("   Message: %s\n", result.Message)
	fmt.Printf("   Duration: %v\n", result.Duration)
	if result.Details != nil {
		fmt.Printf("   Details: %v\n", result.Details)
	}
	fmt.Println()

	// Demo 2: Degraded service (slow response)
	fmt.Println("2. Degraded Service Check")
	fmt.Println("   - Service is slow but functional\n")

	slowService := health.NewDegradableChecker("slow_api", func(ctx context.Context) error {
		// Simulate slow service
		time.Sleep(150 * time.Millisecond)
		return nil
	}, 100*time.Millisecond)

	hc.Register(slowService)

	result2 := slowService.Check(ctx)
	fmt.Printf("   Service: %s\n", result2.Name)
	fmt.Printf("   Status:  %s\n", result2.Status)
	fmt.Printf("   Message: %s\n", result2.Message)
	fmt.Printf("   Duration: %v\n", result2.Duration)
	fmt.Println()

	// Demo 3: Unhealthy service
	fmt.Println("3. Unhealthy Service Check")
	fmt.Println("   - Service is down or failing\n")

	failedService := health.NewSimpleChecker("database", func(ctx context.Context) error {
		// Simulate failed service
		return fmt.Errorf("connection refused")
	})

	hc.Register(failedService)

	result3 := failedService.Check(ctx)
	fmt.Printf("   Service: %s\n", result3.Name)
	fmt.Printf("   Status:  %s\n", result3.Status)
	fmt.Printf("   Message: %s\n", result3.Message)
	if result3.Error != "" {
		fmt.Printf("   Error:   %s\n", result3.Error)
	}
	fmt.Println()

	// Demo 4: Multiple checks
	fmt.Println("4. Multiple Service Checks")
	fmt.Println("   - Running all checks concurrently\n")

	start := time.Now()
	allResults := hc.Check(ctx)
	duration := time.Since(start)

	fmt.Printf("   Checked %d services in %v\n\n", len(allResults), duration)

	for name, res := range allResults {
		statusIcon := "✓"
		if res.Status == health.StatusUnhealthy {
			statusIcon = "✗"
		} else if res.Status == health.StatusDegraded {
			statusIcon = "⚠"
		}
		fmt.Printf("   %s %-15s: %-10s (%v)\n", statusIcon, name, res.Status, res.Duration)
	}
	fmt.Println()

	// Demo 5: Overall health status
	fmt.Println("5. Overall Health Status")
	overallStatus := hc.GetStatus(ctx)
	fmt.Printf("   Overall: %s\n", overallStatus)
	if overallStatus == health.StatusHealthy {
		fmt.Println("   ✓ All systems operational")
	} else if overallStatus == health.StatusDegraded {
		fmt.Println("   ⚠ Some systems degraded")
	} else {
		fmt.Println("   ✗ System unhealthy")
	}
	fmt.Println()

	// Demo 6: Health summary
	fmt.Println("6. Health Summary")
	summary := hc.GetSummary(ctx)
	fmt.Printf("   Status:          %s\n", summary["status"])
	fmt.Printf("   Total Checks:    %d\n", summary["total_checks"])
	fmt.Printf("   Healthy:         %d\n", summary["healthy_count"])
	fmt.Printf("   Degraded:        %d\n", summary["degraded_count"])
	fmt.Printf("   Unhealthy:       %d\n", summary["unhealthy_count"])
	fmt.Println()

	// Demo 7: Check specific service
	fmt.Println("7. Check Specific Service")
	specificResult, err := hc.CheckOne(ctx, "api_service")
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Service: %s\n", specificResult.Name)
		fmt.Printf("   Status:  %s\n", specificResult.Status)
		fmt.Printf("   Message: %s\n", specificResult.Message)
	}
	fmt.Println()

	// Demo 8: Custom checkers
	fmt.Println("8. Custom Health Checkers")
	fmt.Println("   - Demonstrating different checker types\n")

	// Memory checker
	memChecker := health.NewMemoryChecker(70.0, 90.0)
	hc.Register(memChecker)
	memResult := memChecker.Check(ctx)
	fmt.Printf("   Memory Check: %s (%s)\n", memResult.Status, memResult.Message)

	// Disk checker
	diskChecker := health.NewDiskSpaceChecker("/", 80.0, 95.0)
	hc.Register(diskChecker)
	diskResult := diskChecker.Check(ctx)
	fmt.Printf("   Disk Check:   %s (%s)\n", diskResult.Status, diskResult.Message)

	// HTTP service checker
	httpChecker := health.NewHTTPServiceChecker("external_api", "https://api.example.com/health", 200)
	hc.Register(httpChecker)
	httpResult := httpChecker.Check(ctx)
	fmt.Printf("   HTTP Check:   %s (%s)\n", httpResult.Status, httpResult.Message)
	fmt.Println()

	// Summary table
	fmt.Println("=== Health Check Types ===")
	fmt.Println("┌────────────────────┬─────────────────────────┬──────────────────────┐")
	fmt.Println("│ Checker Type       │ Use Case                │ Features             │")
	fmt.Println("├────────────────────┼─────────────────────────┼──────────────────────┤")
	fmt.Println("│ Simple             │ Basic service check     │ Healthy/Unhealthy    │")
	fmt.Println("├────────────────────┼─────────────────────────┼──────────────────────┤")
	fmt.Println("│ Degradable         │ Performance monitoring  │ +Degraded status     │")
	fmt.Println("├────────────────────┼─────────────────────────┼──────────────────────┤")
	fmt.Println("│ Database           │ DB connection check     │ Connection pool info │")
	fmt.Println("├────────────────────┼─────────────────────────┼──────────────────────┤")
	fmt.Println("│ Redis              │ Cache service check     │ Ping + stats         │")
	fmt.Println("├────────────────────┼─────────────────────────┼──────────────────────┤")
	fmt.Println("│ Memory             │ System resource check   │ Usage thresholds     │")
	fmt.Println("├────────────────────┼─────────────────────────┼──────────────────────┤")
	fmt.Println("│ Disk Space         │ Storage monitoring      │ Warning/Critical     │")
	fmt.Println("├────────────────────┼─────────────────────────┼──────────────────────┤")
	fmt.Println("│ HTTP Service       │ External API check      │ Status code validate │")
	fmt.Println("└────────────────────┴─────────────────────────┴──────────────────────┘")

	fmt.Println("\n=== Kubernetes Integration ===")
	fmt.Println("Endpoints:")
	fmt.Println("  GET /health          - Full health check")
	fmt.Println("  GET /health/ready    - Readiness probe")
	fmt.Println("  GET /health/live     - Liveness probe")
	fmt.Println("  GET /health/summary  - Summary view")

	fmt.Println("\n✓ Health check demo completed!")
}
