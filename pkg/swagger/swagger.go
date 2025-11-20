package swagger

import (
	"context"
	"net/http"
	"net/url"

	"github.com/chenyusolar/aidecms/pkg/framework"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SwaggerHandler 创建 Swagger UI 处理器
func SwaggerHandler() framework.HandlerFunc {
	return func(ctx context.Context, c *framework.RequestContext) {
		// 获取请求路径
		path := string(c.RequestContext.Request.URI().Path())

		// 如果是根路径，重定向到 index.html
		if path == "/swagger/" || path == "/swagger" {
			c.Redirect(consts.StatusMovedPermanently, "/swagger/index.html")
			return
		}

		// 创建一个简单的适配器来处理 Swagger UI
		handler := httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
		)

		// 转换 URL
		reqURL, _ := url.Parse(string(c.RequestContext.Request.URI().RequestURI()))
		req := &http.Request{
			Method:     string(c.RequestContext.Request.Method()),
			URL:        reqURL,
			RequestURI: string(c.RequestContext.Request.RequestURI()),
			Header:     convertHeader(&c.RequestContext.Request.Header),
		}

		// 使用 Hertz 适配器
		handler.ServeHTTP(&responseWriter{c}, req)
	}
}

// responseWriter 适配器，将 Hertz 的 ResponseWriter 适配为 http.ResponseWriter
type responseWriter struct {
	ctx *framework.RequestContext
}

func (w *responseWriter) Header() http.Header {
	header := make(http.Header)
	w.ctx.RequestContext.Response.Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.ctx.RequestContext.Response.BodyWriter().Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ctx.RequestContext.Response.SetStatusCode(statusCode)
}

// convertHeader 转换 Hertz 的 Header 为标准 http.Header
func convertHeader(h *protocol.RequestHeader) http.Header {
	header := make(http.Header)
	h.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}
