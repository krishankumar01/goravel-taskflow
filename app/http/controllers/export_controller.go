package controllers

import (
	"errors"
	"fmt"
	"mime"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/http/authuser"
	"taskflow/app/jobs"
	"taskflow/app/models"
)

type ExportController struct {
	// Dependent services
}

func NewExportController() *ExportController {
	return &ExportController{
		// Inject services
	}
}

func (r *ExportController) Store(ctx http.Context) http.Response {
	export := models.Export{
		UserID: authuser.From(ctx).ID,
		Status: enums.ExportPending,
	}
	if err := facades.Orm().WithContext(ctx).Query().Create(&export); err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not start the export.",
		})
	}

	err := facades.Queue().Job(jobs.NewExportTasks(export.ID)).Dispatch()
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not start the export.",
		})
	}

	return ctx.Response().Json(202, export)
}

func (r *ExportController) Show(ctx http.Context) http.Response {
	var export models.Export
	err := facades.Orm().WithContext(ctx).Query().
		Where("user_id", authuser.From(ctx).ID).
		FindOrFail(&export, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The export does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the export.",
		})
	}

	return ctx.Response().Json(200, export)
}

func (r *ExportController) Download(ctx http.Context) http.Response {
	var export models.Export
	err := facades.Orm().WithContext(ctx).Query().
		Where("user_id", authuser.From(ctx).ID).
		FindOrFail(&export, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The export does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the export.",
		})
	}
	if export.Status != enums.ExportDone {
		return ctx.Response().Json(409, map[string]any{
			"error": "The export is not ready yet.",
		})
	}

	content, err := facades.Storage().GetBytes(export.Path)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the export.",
		})
	}

	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": fmt.Sprintf("tasks-%d.csv", export.ID)})

	return ctx.Response().
		Header("Content-Disposition", disposition).
		Data(200, "text/csv", content)
}
