package seo

import "strings"

// RobotsTxt Robots.txt 配置
type RobotsTxt struct {
	UserAgents []UserAgent
	Sitemap    string
}

// UserAgent 用户代理配置
type UserAgent struct {
	Agent      string
	Allow      []string
	Disallow   []string
	CrawlDelay int
}

// NewRobotsTxt 创建 robots.txt 配置
func NewRobotsTxt() *RobotsTxt {
	return &RobotsTxt{
		UserAgents: []UserAgent{},
	}
}

// AddUserAgent 添加用户代理
func (r *RobotsTxt) AddUserAgent(agent string) *RobotsTxt {
	r.UserAgents = append(r.UserAgents, UserAgent{
		Agent:    agent,
		Allow:    []string{},
		Disallow: []string{},
	})
	return r
}

// Allow 允许访问路径
func (r *RobotsTxt) Allow(paths ...string) *RobotsTxt {
	if len(r.UserAgents) > 0 {
		idx := len(r.UserAgents) - 1
		r.UserAgents[idx].Allow = append(r.UserAgents[idx].Allow, paths...)
	}
	return r
}

// Disallow 禁止访问路径
func (r *RobotsTxt) Disallow(paths ...string) *RobotsTxt {
	if len(r.UserAgents) > 0 {
		idx := len(r.UserAgents) - 1
		r.UserAgents[idx].Disallow = append(r.UserAgents[idx].Disallow, paths...)
	}
	return r
}

// SetCrawlDelay 设置爬取延迟
func (r *RobotsTxt) SetCrawlDelay(seconds int) *RobotsTxt {
	if len(r.UserAgents) > 0 {
		idx := len(r.UserAgents) - 1
		r.UserAgents[idx].CrawlDelay = seconds
	}
	return r
}

// SetSitemap 设置 sitemap URL
func (r *RobotsTxt) SetSitemap(url string) *RobotsTxt {
	r.Sitemap = url
	return r
}

// Generate 生成 robots.txt 内容
func (r *RobotsTxt) Generate() string {
	var builder strings.Builder

	for _, ua := range r.UserAgents {
		builder.WriteString("User-agent: " + ua.Agent + "\n")

		for _, allow := range ua.Allow {
			builder.WriteString("Allow: " + allow + "\n")
		}

		for _, disallow := range ua.Disallow {
			builder.WriteString("Disallow: " + disallow + "\n")
		}

		if ua.CrawlDelay > 0 {
			builder.WriteString("Crawl-delay: ")
			builder.WriteString(string(rune(ua.CrawlDelay + '0')))
			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	if r.Sitemap != "" {
		builder.WriteString("Sitemap: " + r.Sitemap + "\n")
	}

	return builder.String()
}

// DefaultRobotsTxt 创建默认的 robots.txt
func DefaultRobotsTxt(sitemapURL string) *RobotsTxt {
	return NewRobotsTxt().
		AddUserAgent("*").
		Allow("/").
		Disallow("/api/cms/").
		Disallow("/admin/").
		Disallow("/storage/").
		SetCrawlDelay(1).
		SetSitemap(sitemapURL)
}
