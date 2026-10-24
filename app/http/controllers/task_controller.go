package controllers

import (
	"taskflow/app/facades"
	"taskflow/app/models"

	"github.com/goravel/framework/contracts/http"
)

type CreateTaskRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}

type TaskController struct {
	// Dependent services
}

func NewTaskController() *TaskController {
	return &TaskController{
		// Inject services
	}
}

func (r *TaskController) Store(ctx http.Context) http.Response {
	var request CreateTaskRequest
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Json(400, map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(200, request)
}

func (r *TaskController) Show(ctx http.Context) http.Response {
	validator, err := facades.Validation().Make(ctx, map[string]any{
		"id": ctx.Request().Route("id"),
	}, map[string]any{
		"id": "bail|integer|min:1",
	})
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": err.Error(),
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(422, map[string]any{
			"errors": validator.Errors().All(),
		})
	}

	var task models.Task
	err = facades.Orm().Query().
		With("Creator").
		With("Assignee").
		Find(&task, ctx.Request().RouteInt("id"))
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(200, task)
}

func (r *TaskController) Index(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]any{
		"status": "in:pending,in-progress,done",
		"page":   "bail|filled|integer|min:1",
	})
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": err.Error(),
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(422, map[string]any{
			"errors": validator.Errors().All(),
		})
	}

	status := ctx.Request().Query("status")
	page := ctx.Request().QueryInt("page", 1)
	filters := ctx.Request().QueryMap("filters")

	return ctx.Response().Json(200, map[string]any{
		"status":  status,
		"page":    page,
		"filters": filters,
	})
}
