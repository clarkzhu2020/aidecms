package controllers

import (
	"context"
	"example/todolist/internal/models"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
)

type TodoController struct {
	// 这里可以注入服务或存储库
}

func NewTodoController() *TodoController {
	return &TodoController{}
}

func (c *TodoController) ListTodos(ctx context.Context, rc *app.RequestContext) {
	// 获取所有Todo项
	todos := []models.Todo{
		{ID: 1, Title: "学习Hertz框架", Completed: false},
		{ID: 2, Title: "开发TodoList应用", Completed: true},
	}

	rc.JSON(200, map[string]interface{}{
		"data": todos,
	})
}

func (c *TodoController) CreateTodo(ctx context.Context, rc *app.RequestContext) {
	var todo models.Todo
	if err := rc.Bind(&todo); err != nil {
		rc.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}

	// 这里应该是保存到数据库的逻辑
	fmt.Printf("创建Todo: %+v\n", todo)

	rc.JSON(201, map[string]interface{}{
		"data": todo,
	})
}

func (c *TodoController) GetTodo(ctx context.Context, rc *app.RequestContext) {
	id := rc.Param("id")

	// 模拟从数据库获取
	todo := models.Todo{
		ID:        1,
		Title:     "示例Todo",
		Completed: false,
	}

	rc.JSON(200, map[string]interface{}{
		"data": todo,
	})
}

func (c *TodoController) UpdateTodo(ctx context.Context, rc *app.RequestContext) {
	id := rc.Param("id")
	var todo models.Todo
	if err := rc.Bind(&todo); err != nil {
		rc.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}

	// 更新逻辑
	fmt.Printf("更新Todo ID %s: %+v\n", id, todo)

	rc.JSON(200, map[string]interface{}{
		"data": todo,
	})
}

func (c *TodoController) DeleteTodo(ctx context.Context, rc *app.RequestContext) {
	id := rc.Param("id")

	// 删除逻辑
	fmt.Printf("删除Todo ID %s\n", id)

	rc.JSON(204, nil)
}
