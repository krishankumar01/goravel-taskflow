package controllers

import (
	"errors"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/http/authuser"
	"taskflow/app/http/requests"
	"taskflow/app/models"
	"taskflow/app/policies"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"
	"github.com/goravel/framework/validation"
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
	if response := facades.Gate().WithContext(ctx).Inspect(policies.CreateTask, nil); !response.Allowed() {
		return ctx.Response().Json(403, map[string]any{
			"error": response.Message(),
		})
	}

	var request requests.StoreTaskRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the task.",
		})
	}
	if validationErrors != nil {
		return ctx.Response().Json(422, map[string]any{
			"errors": validationErrors.All(),
		})
	}

	task := models.Task{
		CreatedBy:   authuser.From(ctx).ID,
		AssignedTo:  request.AssignedTo,
		Title:       request.Title,
		Description: request.Description,
		Status:      enums.StatusPending,
		DueAt:       request.DueAt,
	}
	if err := facades.Orm().WithContext(ctx).Query().Create(&task); err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the task.",
		})
	}

	var saved models.Task
	err = facades.Orm().WithContext(ctx).Query().With("Creator").With("Assignee").Find(&saved, task.ID)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the task.",
		})
	}

	return ctx.Response().Json(201, saved)
}

func (r *TaskController) Show(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().WithContext(ctx).Query().
		With("Creator").
		With("Assignee").
		FindOrFail(&task, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The task does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the task.",
		})
	}

	return ctx.Response().Json(200, task)
}

func (r *TaskController) Index(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]any{
		"status":   "in:pending,in-progress,done",
		"page":     "bail|filled|integer|min:1",
		"per_page": "bail|filled|integer|min:1|max:100",
	})
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not list the tasks.",
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(422, map[string]any{
			"errors": validator.Errors().All(),
		})
	}

	page := ctx.Request().QueryInt("page", 1)
	perPage := ctx.Request().QueryInt("per_page", 15)

	query := facades.Orm().WithContext(ctx).Query()
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	var tasks []models.Task
	var total int64
	err = query.With("Creator").With("Assignee").OrderByDesc("id").Paginate(page, perPage, &tasks, &total)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not list the tasks.",
		})
	}

	return ctx.Response().Json(200, map[string]any{
		"data": tasks,
		"meta": map[string]any{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

func (r *TaskController) Update(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().WithContext(ctx).Query().FindOrFail(&task, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The task does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the task.",
		})
	}

	var request requests.UpdateTaskRequest
	if err := request.Authorize(ctx); err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the task.",
		})
	}

	validator, err := ctx.Request().Validate(
		request.Rules(ctx),
		validation.Filters(request.Filters(ctx)),
		validation.Attributes(request.Attributes(ctx)),
	)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the task.",
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(422, map[string]any{
			"errors": validator.Errors().All(),
		})
	}

	if changes := validator.Validated(); len(changes) > 0 {
		if _, err := facades.Orm().WithContext(ctx).Query().Model(&task).Update(changes); err != nil {
			facades.Log().Error(err)
			return ctx.Response().Json(500, map[string]any{
				"error": "Could not save the task.",
			})
		}
	}

	var saved models.Task
	err = facades.Orm().WithContext(ctx).Query().With("Creator").With("Assignee").Find(&saved, task.ID)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the task.",
		})
	}

	return ctx.Response().Json(200, saved)
}

func (r *TaskController) Destroy(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().WithContext(ctx).Query().FindOrFail(&task, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The task does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the task.",
		})
	}

	if response := facades.Gate().WithContext(ctx).Inspect(policies.DeleteTask, nil); !response.Allowed() {
		return ctx.Response().Json(403, map[string]any{
			"error": response.Message(),
		})
	}

	if _, err := facades.Orm().WithContext(ctx).Query().Delete(&task); err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not delete the task.",
		})
	}

	return ctx.Response().NoContent()
}
