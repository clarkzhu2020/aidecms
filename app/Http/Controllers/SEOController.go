package controllers

import (
	"context"
	"net/http"

	"github.com/chenyusolar/aidecms/pkg/database"
	"github.com/chenyusolar/aidecms/pkg/seo"
	"github.com/cloudwego/hertz/pkg/app"
)

// SEOController SEO控制器
type SEOController struct {
	sitemapGen  *seo.SitemapGenerator
	robotsTxt   *seo.RobotsTxt
	redirectMgr *seo.RedirectManager
}

// NewSEOController 创建SEO控制器
func NewSEOController(baseURL string) *SEOController {
	db := database.GetDB()

	return &SEOController{
		sitemapGen:  seo.NewSitemapGenerator(baseURL, db),
		robotsTxt:   seo.DefaultRobotsTxt(baseURL + "/sitemap.xml"),
		redirectMgr: seo.NewRedirectManager(db),
	}
}

// Sitemap 生成sitemap.xml
// @Summary      生成sitemap
// @Description  生成站点地图XML
// @Tags         SEO
// @Produce      xml
// @Success      200 {string} string "XML格式的sitemap"
// @Router       /sitemap.xml [get]
func (c *SEOController) Sitemap(ctx context.Context, hCtx *app.RequestContext) {
	sitemap, err := c.sitemapGen.Generate()
	if err != nil {
		hCtx.String(http.StatusInternalServerError, "Failed to generate sitemap")
		return
	}

	hCtx.Header("Content-Type", "application/xml")
	hCtx.String(http.StatusOK, sitemap)
}

// PostsSitemap 生成posts sitemap
// @Summary      生成文章sitemap
// @Description  生成文章站点地图XML
// @Tags         SEO
// @Produce      xml
// @Success      200 {string} string "XML格式的sitemap"
// @Router       /sitemap-posts.xml [get]
func (c *SEOController) PostsSitemap(ctx context.Context, hCtx *app.RequestContext) {
	sitemap, err := c.sitemapGen.GenerateForPosts()
	if err != nil {
		hCtx.String(http.StatusInternalServerError, "Failed to generate posts sitemap")
		return
	}

	hCtx.Header("Content-Type", "application/xml")
	hCtx.String(http.StatusOK, sitemap)
}

// Robots 生成robots.txt
// @Summary      生成robots.txt
// @Description  生成搜索引擎爬虫规则文件
// @Tags         SEO
// @Produce      plain
// @Success      200 {string} string "robots.txt内容"
// @Router       /robots.txt [get]
func (c *SEOController) Robots(ctx context.Context, hCtx *app.RequestContext) {
	robots := c.robotsTxt.Generate()
	hCtx.Header("Content-Type", "text/plain")
	hCtx.String(http.StatusOK, robots)
}
