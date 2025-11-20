package routes

import (
	"fmt"

	controllers "github.com/chenyusolar/aidecms/app/Http/Controllers"
	middleware "github.com/chenyusolar/aidecms/app/Http/Middleware"
	"github.com/chenyusolar/aidecms/config"
	"github.com/chenyusolar/aidecms/internal/app/adapters"
	"github.com/chenyusolar/aidecms/pkg/framework"
)

func APIRoutes(app *framework.Application) {
	userController := controllers.NewUserController(app)

	// 创建AI控制器
	manager, err := config.LoadAIManager()
	if err != nil {
		// 如果AI配置加载失败，记录错误但不影响其他路由
		fmt.Printf("Warning: Failed to load AI manager: %v\n", err)
		fmt.Println("AI routes will not be available.")
	} else {
		fmt.Printf("AI manager loaded successfully with %d clients\n", len(manager.ListClients()))
	}

	var aiController *controllers.AIController
	if manager != nil {
		aiController = controllers.NewAIController(manager)
	}

	// 创建邮件控制器
	mailController, err := controllers.NewMailController()
	if err != nil {
		fmt.Printf("Warning: Failed to create mail controller: %v\n", err)
		fmt.Println("Mail routes will not be available.")
	}

	// 创建CMS控制器
	mediaController := controllers.NewMediaController()
	postController := controllers.NewPostController()
	categoryController := controllers.NewCategoryController()
	tagController := controllers.NewTagController()
	menuController := controllers.NewMenuController()
	commentController := controllers.NewCommentController()

	// 创建SEO控制器
	seoController := controllers.NewSEOController("http://localhost:8888")

	// 创建Web3控制器
	web3Controller := &controllers.Web3Controller{}

	// 创建Exchange控制器
	exchangeController := &controllers.ExchangeController{}

	app.RegisterRoutes(func(r *framework.Router) {
		// SEO 路由（公开）
		r.GET("/sitemap.xml", adapters.HertzToFramework(seoController.Sitemap))
		r.GET("/sitemap-posts.xml", adapters.HertzToFramework(seoController.PostsSitemap))
		r.GET("/robots.txt", adapters.HertzToFramework(seoController.Robots))

		// 公开路由
		r.POST("/register", userController.Register)
		r.POST("/login", userController.Login)

		// AI 路由
		if aiController != nil {
			fmt.Println("Registering AI routes...")
			r.POST("/api/ai/chat", adapters.HertzToFramework(aiController.Chat))
			r.POST("/api/ai/completion", adapters.HertzToFramework(aiController.Completion))
			r.POST("/api/ai/embedding", adapters.HertzToFramework(aiController.Embedding))
		}

		// 邮件 API 路由
		if mailController != nil {
			fmt.Println("Registering mail routes...")
			r.POST("/api/mail/send", adapters.HertzToFramework(mailController.SendMail))
			r.POST("/api/mail/send-template", adapters.HertzToFramework(mailController.SendTemplate))
			r.POST("/api/mail/send-bulk", adapters.HertzToFramework(mailController.SendBulkMail))
			r.GET("/api/mail/test", adapters.HertzToFramework(mailController.TestConnection))
			r.GET("/api/mail/config", adapters.HertzToFramework(mailController.GetMailConfig))
			r.GET("/api/mail/validate", adapters.HertzToFramework(mailController.ValidateEmail))
		}

		// CMS 公开路由（只读）
		fmt.Println("Registering CMS routes...")
		r.GET("/api/posts", adapters.HertzToFramework(postController.List))
		r.GET("/api/posts/:id", adapters.HertzToFramework(postController.Get))
		r.GET("/api/categories", adapters.HertzToFramework(categoryController.List))
		r.GET("/api/categories/:id", adapters.HertzToFramework(categoryController.Get))
		r.GET("/api/tags", adapters.HertzToFramework(tagController.List))
		r.GET("/api/tags/:id", adapters.HertzToFramework(tagController.Get))
		r.GET("/api/media", adapters.HertzToFramework(mediaController.List))
		r.GET("/api/media/:id", adapters.HertzToFramework(mediaController.Get))
		r.GET("/api/menus", adapters.HertzToFramework(menuController.List))
		r.GET("/api/menus/:id", adapters.HertzToFramework(menuController.Get))
		r.GET("/api/comments", adapters.HertzToFramework(commentController.List))
		r.GET("/api/comments/:id", adapters.HertzToFramework(commentController.Get))
		r.POST("/api/comments", adapters.HertzToFramework(commentController.Create))

		// 需要认证的路由
		authGroup := r.Group("/user", middleware.JWTMiddleware())
		{
			authGroup.GET("/profile", userController.Profile)
			authGroup.PUT("/profile", userController.UpdateProfile)
		}

		// CMS 管理路由（需要认证）
		cmsGroup := r.Group("/api/cms", middleware.JWTMiddleware())
		{
			// 文章管理
			cmsGroup.POST("/posts", adapters.HertzToFramework(postController.Create))
			cmsGroup.PUT("/posts/:id", adapters.HertzToFramework(postController.Update))
			cmsGroup.DELETE("/posts/:id", adapters.HertzToFramework(postController.Delete))
			cmsGroup.POST("/posts/:id/publish", adapters.HertzToFramework(postController.Publish))

			// 分类管理
			cmsGroup.POST("/categories", adapters.HertzToFramework(categoryController.Create))
			cmsGroup.PUT("/categories/:id", adapters.HertzToFramework(categoryController.Update))
			cmsGroup.DELETE("/categories/:id", adapters.HertzToFramework(categoryController.Delete))

			// 标签管理
			cmsGroup.POST("/tags", adapters.HertzToFramework(tagController.Create))
			cmsGroup.PUT("/tags/:id", adapters.HertzToFramework(tagController.Update))
			cmsGroup.DELETE("/tags/:id", adapters.HertzToFramework(tagController.Delete))

			// 媒体管理
			cmsGroup.POST("/media/upload", adapters.HertzToFramework(mediaController.Upload))
			cmsGroup.PUT("/media/:id", adapters.HertzToFramework(mediaController.Update))
			cmsGroup.DELETE("/media/:id", adapters.HertzToFramework(mediaController.Delete))

			// 菜单管理
			cmsGroup.POST("/menus", adapters.HertzToFramework(menuController.Create))
			cmsGroup.PUT("/menus/:id", adapters.HertzToFramework(menuController.Update))
			cmsGroup.DELETE("/menus/:id", adapters.HertzToFramework(menuController.Delete))
			cmsGroup.POST("/menus/reorder", adapters.HertzToFramework(menuController.Reorder))

			// 评论管理
			cmsGroup.PUT("/comments/:id", adapters.HertzToFramework(commentController.Update))
			cmsGroup.DELETE("/comments/:id", adapters.HertzToFramework(commentController.Delete))
			cmsGroup.POST("/comments/:id/approve", adapters.HertzToFramework(commentController.Approve))
			cmsGroup.POST("/comments/:id/spam", adapters.HertzToFramework(commentController.MarkAsSpam))
		}

		// Web3 路由（公开）
		web3Group := r.Group("/api/web3")
		{
			// 区块链基本操作
			web3Group.GET("/:chain/balance/:address", adapters.HertzToFramework(web3Controller.GetBalance))
			web3Group.GET("/:chain/transaction/:hash", adapters.HertzToFramework(web3Controller.GetTransaction))
			web3Group.GET("/:chain/block-number", adapters.HertzToFramework(web3Controller.GetBlockNumber))
			web3Group.GET("/:chain/wallet/:address", adapters.HertzToFramework(web3Controller.GetWalletInfo))
			web3Group.GET("/:chain/validate/:address", adapters.HertzToFramework(web3Controller.ValidateAddress))

			// 支持的链列表
			web3Group.GET("/chains", adapters.HertzToFramework(web3Controller.GetSupportedChains))

			// 多链查询
			web3Group.POST("/multi-balance", adapters.HertzToFramework(web3Controller.GetMultiChainBalances))
		}

		// Exchange 路由（公开）
		exchangeGroup := r.Group("/api/exchange")
		{
			// 单交易所查询
			exchangeGroup.GET("/:exchange/balance/:currency", adapters.HertzToFramework(exchangeController.GetBalance))
			exchangeGroup.GET("/:exchange/balances", adapters.HertzToFramework(exchangeController.GetBalances))
			exchangeGroup.GET("/:exchange/price/:pair", adapters.HertzToFramework(exchangeController.GetPrice))

			// 支持的交易所列表
			exchangeGroup.GET("/supported", adapters.HertzToFramework(exchangeController.GetSupportedExchanges))

			// 多交易所查询
			exchangeGroup.GET("/all/balance/:currency", adapters.HertzToFramework(exchangeController.GetAllBalances))
			exchangeGroup.GET("/all/price/:pair", adapters.HertzToFramework(exchangeController.GetAllPrices))
		}
	})
}
