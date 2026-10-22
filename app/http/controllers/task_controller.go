package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type TaskController struct {
	// Dependent services
}

func NewTaskController() *TaskController {
	return &TaskController{
		// Inject services
	}
}

func (r *TaskController) Store(ctx http.Context) http.Response {
	return ctx.Response().String(200, "task created")
}

func (r *TaskController) Show(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")

	return ctx.Response().Json(200, map[string]any{
		"id": id,
	})
}
