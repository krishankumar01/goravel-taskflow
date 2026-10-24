package controllers

import (
	"errors"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (r *UserController) Index(ctx http.Context) http.Response {
	var users []models.User
	err := facades.Orm().Query().
		With("AssignedTasks", func(query orm.Query) orm.Query {
			return query.Where("status <> ?", enums.StatusDone)
		}).
		With("AssignedTasks.Creator").
		Find(&users)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not list the users.",
		})
	}

	return ctx.Response().Json(200, users)
}

func (r *UserController) Show(ctx http.Context) http.Response {
	var user models.User
	err := facades.Orm().WithContext(ctx).Query().
		With("CreatedTasks").
		With("AssignedTasks").
		FindOrFail(&user, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(404, map[string]any{
			"error": "The user does not exist.",
		})
	}
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not load the user.",
		})
	}

	return ctx.Response().Json(200, user)
}
