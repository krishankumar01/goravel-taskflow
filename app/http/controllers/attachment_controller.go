package controllers

import (
	"errors"
	"fmt"
	"mime"
	"os"

	"taskflow/app/facades"
	"taskflow/app/http/authuser"
	"taskflow/app/http/requests"
	"taskflow/app/models"
	"taskflow/app/policies"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"
	"github.com/goravel/framework/filesystem"
)

type AttachmentController struct {
	// Dependent services
}

func NewAttachmentController() *AttachmentController {
	return &AttachmentController{
		// Inject services
	}
}

func (r *AttachmentController) Store(ctx http.Context) http.Response {
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

	response := facades.Gate().WithContext(ctx).Inspect(policies.CreateAttachment, map[string]any{
		"task": &task,
	})
	if !response.Allowed() {
		return ctx.Response().Json(403, map[string]any{
			"error": response.Message(),
		})
	}

	var request requests.StoreAttachmentRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the file.",
		})
	}
	if validationErrors != nil {
		return ctx.Response().Json(422, map[string]any{
			"errors": validationErrors.All(),
		})
	}

	file, err := filesystem.NewFileFromRequest(request.File)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the file.",
		})
	}
	defer os.Remove(file.File())

	size, err := file.Size()
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the file.",
		})
	}
	mimeType, err := file.MimeType()
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the file.",
		})
	}

	path, err := file.Store(fmt.Sprintf("attachments/%d", task.ID))
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the file.",
		})
	}

	attachment := models.Attachment{
		TaskID:     task.ID,
		UploadedBy: authuser.From(ctx).ID,
		Name:       file.GetClientOriginalName(),
		Path:       path,
		Size:       size,
		MimeType:   mimeType,
	}
	if err := facades.Orm().WithContext(ctx).Query().Create(&attachment); err != nil {
		facades.Log().Error(err)
		if err := facades.Storage().Delete(path); err != nil {
			facades.Log().Error(err)
		}
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not save the file.",
		})
	}

	return ctx.Response().Json(201, attachment)
}

func (r *AttachmentController) Show(ctx http.Context) http.Response {
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

	response := facades.Gate().WithContext(ctx).Inspect(policies.ViewAttachment, map[string]any{
		"task": &task,
	})
	if !response.Allowed() {
		return ctx.Response().Json(403, map[string]any{
			"error": response.Message(),
		})
	}

	var attachment models.Attachment
	err = facades.Orm().WithContext(ctx).Query().
		Where("task_id", task.ID).
		FindOrFail(&attachment, ctx.Request().RouteInt("attachment"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The file does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the file.",
		})
	}

	content, err := facades.Storage().GetBytes(attachment.Path)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the file.",
		})
	}

	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": attachment.Name})

	return ctx.Response().
		Header("Content-Disposition", disposition).
		Data(200, attachment.MimeType, content)
}

func (r *AttachmentController) Destroy(ctx http.Context) http.Response {
	var attachment models.Attachment
	err := facades.Orm().WithContext(ctx).Query().
		Where("task_id", ctx.Request().RouteInt("id")).
		FindOrFail(&attachment, ctx.Request().RouteInt("attachment"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The file does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the file.",
		})
	}

	response := facades.Gate().WithContext(ctx).Inspect(policies.DeleteAttachment, map[string]any{
		"attachment": &attachment,
	})
	if !response.Allowed() {
		return ctx.Response().Json(403, map[string]any{
			"error": response.Message(),
		})
	}

	if _, err := facades.Orm().WithContext(ctx).Query().Delete(&attachment); err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not delete the file.",
		})
	}

	if err := facades.Storage().Delete(attachment.Path); err != nil {
		facades.Log().Error(err)
	}

	return ctx.Response().NoContent()
}
