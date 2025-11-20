package seo

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/chenyusolar/aidecms/internal/app/models"
	"gorm.io/gorm"
)

// URLSet Sitemap URL 集合
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL Sitemap URL
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float32 `xml:"priority,omitempty"`
}

// SitemapGenerator Sitemap 生成器
type SitemapGenerator struct {
	baseURL string
	db      *gorm.DB
}

// NewSitemapGenerator 创建 Sitemap 生成器
func NewSitemapGenerator(baseURL string, db *gorm.DB) *SitemapGenerator {
	return &SitemapGenerator{
		baseURL: baseURL,
		db:      db,
	}
}

// Generate 生成完整的 sitemap
func (g *SitemapGenerator) Generate() (string, error) {
	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, 0),
	}

	// 添加首页
	urlset.URLs = append(urlset.URLs, URL{
		Loc:        g.baseURL,
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	// 添加文章
	var posts []models.Post
	if err := g.db.Where("status = ?", "published").Find(&posts).Error; err != nil {
		return "", err
	}

	for _, post := range posts {
		lastMod := post.UpdatedAt.Format("2006-01-02")
		if post.PublishedAt != nil {
			lastMod = post.PublishedAt.Format("2006-01-02")
		}

		urlset.URLs = append(urlset.URLs, URL{
			Loc:        fmt.Sprintf("%s/posts/%s", g.baseURL, post.Slug),
			LastMod:    lastMod,
			ChangeFreq: "weekly",
			Priority:   0.8,
		})
	}

	// 添加分类
	var categories []models.Category
	if err := g.db.Find(&categories).Error; err != nil {
		return "", err
	}

	for _, category := range categories {
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        fmt.Sprintf("%s/categories/%s", g.baseURL, category.Slug),
			LastMod:    category.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   0.6,
		})
	}

	// 添加标签
	var tags []models.Tag
	if err := g.db.Find(&tags).Error; err != nil {
		return "", err
	}

	for _, tag := range tags {
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        fmt.Sprintf("%s/tags/%s", g.baseURL, tag.Slug),
			LastMod:    tag.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   0.5,
		})
	}

	// 生成 XML
	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}

// GenerateForPosts 仅生成文章的 sitemap
func (g *SitemapGenerator) GenerateForPosts() (string, error) {
	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, 0),
	}

	var posts []models.Post
	if err := g.db.Where("status = ?", "published").
		Order("published_at DESC").
		Find(&posts).Error; err != nil {
		return "", err
	}

	for _, post := range posts {
		lastMod := post.UpdatedAt.Format("2006-01-02")
		if post.PublishedAt != nil {
			lastMod = post.PublishedAt.Format("2006-01-02")
		}

		urlset.URLs = append(urlset.URLs, URL{
			Loc:        fmt.Sprintf("%s/posts/%s", g.baseURL, post.Slug),
			LastMod:    lastMod,
			ChangeFreq: "weekly",
			Priority:   0.8,
		})
	}

	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}
