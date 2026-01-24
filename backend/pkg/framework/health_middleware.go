package framework

import (
	"context"
	"encoding/json"

	"github.com/clarkzhu2020/aidecms/pkg/health"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// HealthEndpoint 健康检查端点中间件
func HealthEndpoint(checker *health.HealthChecker) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		results := checker.Check(ctx)
		status := checker.GetStatus(ctx)

		// Set HTTP status based on health status
		httpStatus := consts.StatusOK
		if status == health.StatusUnhealthy {
			httpStatus = consts.StatusServiceUnavailable
		} else if status == health.StatusDegraded {
			httpStatus = consts.StatusOK // Still return 200 for degraded
		}

		c.JSON(httpStatus, map[string]interface{}{
			"status": status,
			"checks": results,
		})
	}
}

// HealthSummaryEndpoint 健康检查摘要端点
func HealthSummaryEndpoint(checker *health.HealthChecker) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		summary := checker.GetSummary(ctx)
		status := summary["status"]

		// Set HTTP status based on health status
		httpStatus := consts.StatusOK
		if status == health.StatusUnhealthy {
			httpStatus = consts.StatusServiceUnavailable
		}

		c.JSON(httpStatus, summary)
	}
}

// ReadinessEndpoint 就绪检查端点 (用于 Kubernetes readiness probe)
func ReadinessEndpoint(checker *health.HealthChecker) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		status := checker.GetStatus(ctx)

		if status == health.StatusHealthy || status == health.StatusDegraded {
			c.JSON(consts.StatusOK, map[string]interface{}{
				"ready": true,
			})
		} else {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"ready": false,
			})
		}
	}
}

// LivenessEndpoint 存活检查端点 (用于 Kubernetes liveness probe)
func LivenessEndpoint() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Simple liveness check - if we can respond, we're alive
		c.JSON(consts.StatusOK, map[string]interface{}{
			"alive": true,
		})
	}
}

// HealthCheckMiddleware 健康检查中间件 (可选：记录不健康状态)
func HealthCheckMiddleware(checker *health.HealthChecker) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Check health before processing request
		status := checker.GetStatus(ctx)

		// Add health status to response headers
		c.Response.Header.Set("X-Health-Status", string(status))

		// If service is unhealthy, you might want to reject requests
		// For now, we just add the header and continue
		c.Next(ctx)
	}
}

// DetailedHealthEndpoint 详细健康检查端点 (包含单个检查器的详细信息)
func DetailedHealthEndpoint(checker *health.HealthChecker) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Get specific checker name from query param
		checkerName := c.Query("name")

		if checkerName != "" {
			// Check specific checker
			result, err := checker.CheckOne(ctx, checkerName)
			if err != nil {
				c.JSON(consts.StatusNotFound, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			httpStatus := consts.StatusOK
			if result.Status == health.StatusUnhealthy {
				httpStatus = consts.StatusServiceUnavailable
			}

			c.JSON(httpStatus, result)
		} else {
			// Return all checks with full details
			results := checker.Check(ctx)
			status := checker.GetStatus(ctx)

			httpStatus := consts.StatusOK
			if status == health.StatusUnhealthy {
				httpStatus = consts.StatusServiceUnavailable
			}

			// Format response with detailed information
			response := map[string]interface{}{
				"status": status,
				"checks": results,
			}

			// Pretty print if requested
			if c.Query("pretty") == "true" {
				prettyJSON, err := json.MarshalIndent(response, "", "  ")
				if err == nil {
					c.Data(httpStatus, "application/json", prettyJSON)
					return
				}
			}

			c.JSON(httpStatus, response)
		}
	}
}
