package framework

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// PrometheusMiddleware returns a middleware that records HTTP metrics
func PrometheusMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		c.Next(ctx)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response.StatusCode())
		method := string(c.Request.Method())
		path := string(c.Request.Path())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}

// PrometheusHandler returns the HTTP handler for Prometheus metrics
func PrometheusHandler() HandlerFunc {
	handler := adaptor.HertzHandler(promhttp.Handler())
	return func(ctx context.Context, c *RequestContext) {
		handler(ctx, c.RequestContext)
	}
}
