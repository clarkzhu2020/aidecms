package routes

import (
	"fmt"

	controllers "github.com/clarkzhu2020/aidecms/app/Http/Controllers"
	middleware "github.com/clarkzhu2020/aidecms/app/Http/Middleware"
	"github.com/clarkzhu2020/aidecms/config"
	"github.com/clarkzhu2020/aidecms/internal/app/adapters"
	"github.com/clarkzhu2020/aidecms/internal/app/services"
	"github.com/clarkzhu2020/aidecms/pkg/framework"
)

func APIRoutes(app *framework.Application) {
	userController := controllers.NewUserController(app)

	// AI
	manager, err := config.LoadAIManager()
	if err != nil {
		fmt.Printf("Warning: Failed to load AI manager: %v\n", err)
	}
	var aiController *controllers.AIController
	if manager != nil {
		aiController = controllers.NewAIController(manager)
	}

	// Mail
	mailController, _ := controllers.NewMailController()

	// CMS
	mediaController := controllers.NewMediaController()
	postController := controllers.NewPostController()
	categoryController := controllers.NewCategoryController()
	tagController := controllers.NewTagController()
	menuController := controllers.NewMenuController()
	commentController := controllers.NewCommentController()

	// Payments
	paymentController := controllers.NewPaymentController()
	stripePaymentController := controllers.NewStripePaymentController()
	moonPayController := controllers.NewMoonPayController()
	coinbaseController := controllers.NewCoinbaseController()
	kucoinController := &controllers.KuCoinController{}

	// SEO
	seoController := controllers.NewSEOController("http://localhost:8888")

	// Web3 & Exchange
	web3Controller := &controllers.Web3Controller{}
	mdService := services.NewMarketDataService(app.ClickHouse)
	exController := controllers.NewExchangeController(mdService)

	app.RegisterRoutes(func(r *framework.Router) {
		// SEO
		r.GET("/sitemap.xml", adapters.HertzToFramework(seoController.Sitemap))
		r.GET("/sitemap-posts.xml", adapters.HertzToFramework(seoController.PostsSitemap))
		r.GET("/robots.txt", adapters.HertzToFramework(seoController.Robots))

		// Auth
		r.POST("/register", userController.Register)
		r.POST("/login", userController.Login)

		// AI
		if aiController != nil {
			r.POST("/api/ai/chat", adapters.HertzToFramework(aiController.Chat))
			r.POST("/api/ai/completion", adapters.HertzToFramework(aiController.Completion))
			r.POST("/api/ai/embedding", adapters.HertzToFramework(aiController.Embedding))
		}

		// Mail
		if mailController != nil {
			r.POST("/api/mail/send", adapters.HertzToFramework(mailController.SendMail))
			r.POST("/api/mail/send-template", adapters.HertzToFramework(mailController.SendTemplate))
			r.POST("/api/mail/send-bulk", adapters.HertzToFramework(mailController.SendBulkMail))
			r.GET("/api/mail/test", adapters.HertzToFramework(mailController.TestConnection))
			r.GET("/api/mail/config", adapters.HertzToFramework(mailController.GetMailConfig))
			r.GET("/api/mail/validate", adapters.HertzToFramework(mailController.ValidateEmail))
		}

		// CMS Public
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

		// Payments
		r.GET("/api/payments", adapters.HertzToFramework(paymentController.ListPayments))
		r.GET("/api/payments/:id", adapters.HertzToFramework(paymentController.GetPayment))
		r.POST("/api/payments", adapters.HertzToFramework(paymentController.CreatePayment))
		r.POST("/api/payments/capture/:orderID", adapters.HertzToFramework(paymentController.CapturePayment))
		r.POST("/api/payments/:id/refund", adapters.HertzToFramework(paymentController.RefundPayment))
		r.POST("/api/payments/webhook", adapters.HertzToFramework(paymentController.HandleWebhook))

		// Stripe Payments
		r.GET("/api/stripe/payments", adapters.HertzToFramework(stripePaymentController.ListPayments))
		r.GET("/api/stripe/payments/:id", adapters.HertzToFramework(stripePaymentController.GetPayment))
		r.POST("/api/stripe/payments/intent", adapters.HertzToFramework(stripePaymentController.CreatePaymentIntent))
		r.POST("/api/stripe/payments/:paymentIntentID/confirm", adapters.HertzToFramework(stripePaymentController.ConfirmPayment))
		r.POST("/api/stripe/payments/:id/refund", adapters.HertzToFramework(stripePaymentController.RefundPayment))
		r.POST("/api/stripe/webhook", adapters.HertzToFramework(stripePaymentController.HandleWebhook))

		// MoonPay Payments
		r.GET("/api/moonpay/transactions", adapters.HertzToFramework(moonPayController.ListTransactions))
		r.GET("/api/moonpay/transactions/:id", adapters.HertzToFramework(moonPayController.GetTransaction))
		r.POST("/api/moonpay/transactions", adapters.HertzToFramework(moonPayController.CreateTransaction))
		r.POST("/api/moonpay/quote", adapters.HertzToFramework(moonPayController.GetQuote))
		r.POST("/api/moonpay/widget-url", adapters.HertzToFramework(moonPayController.GenerateWidgetURL))
		r.POST("/api/moonpay/webhook", adapters.HertzToFramework(moonPayController.HandleWebhook))

		// Coinbase Payments (Payment Links)
		r.GET("/api/coinbase/payment-links", adapters.HertzToFramework(coinbaseController.ListPaymentLinks))
		r.GET("/api/coinbase/payment-links/:id", adapters.HertzToFramework(coinbaseController.GetPaymentLink))
		r.POST("/api/coinbase/payment-links", adapters.HertzToFramework(coinbaseController.CreatePaymentLink))
		r.DELETE("/api/coinbase/payment-links/:id", adapters.HertzToFramework(coinbaseController.DeletePaymentLink))

		// Coinbase Trading
		r.GET("/api/coinbase/orders", adapters.HertzToFramework(coinbaseController.ListOrders))
		r.GET("/api/coinbase/orders/:id", adapters.HertzToFramework(coinbaseController.GetOrder))
		r.POST("/api/coinbase/orders", adapters.HertzToFramework(coinbaseController.CreateOrder))
		r.POST("/api/coinbase/orders/:id/cancel", adapters.HertzToFramework(coinbaseController.CancelOrder))
		r.GET("/api/coinbase/products", adapters.HertzToFramework(coinbaseController.GetProducts))
		r.GET("/api/coinbase/products/:productId", adapters.HertzToFramework(coinbaseController.GetProduct))
		r.GET("/api/coinbase/products/:productId/ticker", adapters.HertzToFramework(coinbaseController.GetTicker))
		r.POST("/api/coinbase/webhook", adapters.HertzToFramework(coinbaseController.HandleWebhook))

		// KuCoin Trading
		r.GET("/api/kucoin/server-time", adapters.HertzToFramework(kucoinController.GetServerTime))
		r.GET("/api/kucoin/symbols", adapters.HertzToFramework(kucoinController.GetSymbols))
		r.GET("/api/kucoin/ticker", adapters.HertzToFramework(kucoinController.GetTicker))
		r.GET("/api/kucoin/orderbook", adapters.HertzToFramework(kucoinController.GetOrderBook))
		r.GET("/api/kucoin/trades", adapters.HertzToFramework(kucoinController.GetMarketTrades))
		r.GET("/api/kucoin/klines", adapters.HertzToFramework(kucoinController.GetKlines))
		r.GET("/api/kucoin/24h-stats", adapters.HertzToFramework(kucoinController.Get24HStats))
		r.GET("/api/kucoin/accounts", adapters.HertzToFramework(kucoinController.GetAccounts))
		r.GET("/api/kucoin/accounts/:accountId", adapters.HertzToFramework(kucoinController.GetAccountDetail))
		r.POST("/api/kucoin/orders", adapters.HertzToFramework(kucoinController.CreateOrder))
		r.DELETE("/api/kucoin/orders/:orderId", adapters.HertzToFramework(kucoinController.CancelOrder))
		r.GET("/api/kucoin/orders/:orderId", adapters.HertzToFramework(kucoinController.GetOrder))
		r.GET("/api/kucoin/orders", adapters.HertzToFramework(kucoinController.GetOpenOrders))
		r.GET("/api/kucoin/orders/closed", adapters.HertzToFramework(kucoinController.GetClosedOrders))
		r.POST("/api/kucoin/accounts/sync", adapters.HertzToFramework(kucoinController.SyncAccounts))
		r.POST("/api/kucoin/balance/snapshot", adapters.HertzToFramework(kucoinController.CreateBalanceSnapshot))

		// User
		authGroup := r.Group("/user", middleware.JWTMiddleware())
		{
			authGroup.GET("/profile", userController.Profile)
			authGroup.PUT("/profile", userController.UpdateProfile)
		}

		// CMS Admin
		cmsGroup := r.Group("/api/cms", middleware.JWTMiddleware())
		{
			cmsGroup.POST("/posts", adapters.HertzToFramework(postController.Create))
			cmsGroup.PUT("/posts/:id", adapters.HertzToFramework(postController.Update))
			cmsGroup.DELETE("/posts/:id", adapters.HertzToFramework(postController.Delete))
			cmsGroup.POST("/posts/:id/publish", adapters.HertzToFramework(postController.Publish))
			cmsGroup.POST("/categories", adapters.HertzToFramework(categoryController.Create))
			cmsGroup.PUT("/categories/:id", adapters.HertzToFramework(categoryController.Update))
			cmsGroup.DELETE("/categories/:id", adapters.HertzToFramework(categoryController.Delete))
			cmsGroup.POST("/tags", adapters.HertzToFramework(tagController.Create))
			cmsGroup.PUT("/tags/:id", adapters.HertzToFramework(tagController.Update))
			cmsGroup.DELETE("/tags/:id", adapters.HertzToFramework(tagController.Delete))
			cmsGroup.POST("/media/upload", adapters.HertzToFramework(mediaController.Upload))
			cmsGroup.PUT("/media/:id", adapters.HertzToFramework(mediaController.Update))
			cmsGroup.DELETE("/media/:id", adapters.HertzToFramework(mediaController.Delete))
			cmsGroup.POST("/menus", adapters.HertzToFramework(menuController.Create))
			cmsGroup.PUT("/menus/:id", adapters.HertzToFramework(menuController.Update))
			cmsGroup.DELETE("/menus/:id", adapters.HertzToFramework(menuController.Delete))
			cmsGroup.POST("/menus/reorder", adapters.HertzToFramework(menuController.Reorder))
			cmsGroup.PUT("/comments/:id", adapters.HertzToFramework(commentController.Update))
			cmsGroup.DELETE("/comments/:id", adapters.HertzToFramework(commentController.Delete))
			cmsGroup.POST("/comments/:id/approve", adapters.HertzToFramework(commentController.Approve))
			cmsGroup.POST("/comments/:id/spam", adapters.HertzToFramework(commentController.MarkAsSpam))
		}

		// Web3
		web3Group := r.Group("/api/web3")
		{
			web3Group.GET("/:chain/balance/:address", adapters.HertzToFramework(web3Controller.GetBalance))
			web3Group.GET("/:chain/transaction/:hash", adapters.HertzToFramework(web3Controller.GetTransaction))
			web3Group.GET("/:chain/block-number", adapters.HertzToFramework(web3Controller.GetBlockNumber))
			web3Group.GET("/:chain/wallet/:address", adapters.HertzToFramework(web3Controller.GetWalletInfo))
			web3Group.GET("/:chain/validate/:address", adapters.HertzToFramework(web3Controller.ValidateAddress))
			web3Group.GET("/chains", adapters.HertzToFramework(web3Controller.GetSupportedChains))
			web3Group.POST("/multi-balance", adapters.HertzToFramework(web3Controller.GetMultiChainBalances))
		}

		// Exchange
		exGroup := r.Group("/api/exchange")
		{
			exGroup.GET("/:exchange/balance/:currency", adapters.HertzToFramework(exController.GetBalance))
			exGroup.GET("/:exchange/balances", adapters.HertzToFramework(exController.GetBalances))
			exGroup.GET("/:exchange/price/:pair", adapters.HertzToFramework(exController.GetPrice))
			exGroup.GET("/supported", adapters.HertzToFramework(exController.GetSupportedExchanges))
			exGroup.GET("/all/balance/:currency", adapters.HertzToFramework(exController.GetAllBalances))
			exGroup.GET("/all/price/:pair", adapters.HertzToFramework(exController.GetAllPrices))
		}
	})
}
