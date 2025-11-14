package routes

import (
	"fmt"

	controllers "github.com/clarkgo/clarkgo/app/Http/Controllers"
	middleware "github.com/clarkgo/clarkgo/app/Http/Middleware"
	"github.com/clarkgo/clarkgo/config"
	"github.com/clarkgo/clarkgo/internal/app/adapters"
	"github.com/clarkgo/clarkgo/pkg/framework"
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

	app.RegisterRoutes(func(r *framework.Router) {
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

		// 需要认证的路由
		authGroup := r.Group("/user", middleware.JWTMiddleware())
		{
			authGroup.GET("/profile", userController.Profile)
			authGroup.PUT("/profile", userController.UpdateProfile)
		}
	})
}
