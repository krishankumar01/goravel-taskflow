package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"

	"taskflow/app/facades"
	"taskflow/app/http/authuser"
	"taskflow/app/models"
)

type NotificationController struct {
	// Dependent services
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		// Inject services
	}
}

func (r *NotificationController) Index(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]any{
		"page":     "bail|filled|integer|min:1",
		"per_page": "bail|filled|integer|min:1|max:100",
	})
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not list the notifications.",
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(422, map[string]any{
			"errors": validator.Errors().All(),
		})
	}

	page := ctx.Request().QueryInt("page", 1)
	perPage := ctx.Request().QueryInt("per_page", 15)

	var notifications []models.Notification
	var total int64
	err = facades.Orm().WithContext(ctx).Query().
		Where("user_id", authuser.From(ctx).ID).
		OrderByDesc("id").
		Paginate(page, perPage, &notifications, &total)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not list the notifications.",
		})
	}

	return ctx.Response().Json(200, map[string]any{
		"data": notifications,
		"meta": map[string]any{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

func (r *NotificationController) Read(ctx http.Context) http.Response {
	var notification models.Notification
	err := facades.Orm().WithContext(ctx).Query().
		Where("user_id", authuser.From(ctx).ID).
		FindOrFail(&notification, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The notification does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the notification.",
		})
	}

	if notification.ReadAt == nil {
		notification.ReadAt = carbon.NewDateTime(carbon.Now().StartOfSecond())
		if err := facades.Orm().WithContext(ctx).Query().Save(&notification); err != nil {
			facades.Log().Error(err)
			return ctx.Response().Json(500, map[string]any{
				"error": "Could not save the notification.",
			})
		}
	}

	return ctx.Response().Json(200, notification)
}
