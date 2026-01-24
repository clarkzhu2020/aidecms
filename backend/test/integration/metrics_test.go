package integration

import (
	"context"
	"testing"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/framework"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestMetricsEndpoint(t *testing.T) {
	// Setup server
	app := framework.NewApplication().
		SetDebug(false).
		Boot()

	app.RegisterMiddleware(framework.PrometheusMiddleware())
	app.RegisterRoutes(func(r *framework.Router) {
		r.GET("/metrics", framework.PrometheusHandler())
		r.GET("/ping", func(ctx context.Context, c *framework.RequestContext) {
			c.String(200, "pong")
		})
	})

	go app.Server.Run()
	defer app.Server.Close()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	c, _ := client.NewClient()
	req := &protocol.Request{}
	resp := &protocol.Response{}

	// 1. Make a request to trigger metrics
	req.SetMethod("GET")
	req.SetRequestURI("http://localhost:8888/ping")
	err := c.Do(context.Background(), req, resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())

	// 2. Check metrics endpoint
	req.SetMethod("GET")
	req.SetRequestURI("http://localhost:8888/metrics")
	err = c.Do(context.Background(), req, resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())

	body := string(resp.Body())
	assert.Contains(t, body, "http_requests_total")
	assert.Contains(t, body, "http_request_duration_seconds")
}
