package main

import (
	"example/todolist/internal/controllers"
	"example/todolist/pkg/framework"
)

func main() {
	// 创建应用实例
	application := framework.NewApplication()

	// 设置调试模式
	application.SetDebug(true)

	// 注册路由
	application.RegisterRoutes(func(r *framework.Router) {
		// TodoList路由
		todoController := controllers.NewTodoController()
		r.GET("/todos", todoController.ListTodos)
		r.POST("/todos", todoController.CreateTodo)
		r.GET("/todos/:id", todoController.GetTodo)
		r.PUT("/todos/:id", todoController.UpdateTodo)
		r.DELETE("/todos/:id", todoController.DeleteTodo)
	})

	// 启动应用
	application.Run()
}
