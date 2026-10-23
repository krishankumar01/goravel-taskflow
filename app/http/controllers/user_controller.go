package controllers

import (
	"taskflow/app/facades"
	"taskflow/app/models"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (r *UserController) Index(ctx http.Context) http.Response {
	var users []models.User
	err := facades.Orm().Query().
		With("AssignedTasks", func(query orm.Query) orm.Query {
			return query.Where("status <> ?", models.StatusDone)
		}).
		With("AssignedTasks.Creator").
		Find(&users)
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(200, users)
}

func (r *UserController) Show(ctx http.Context) http.Response {
	var user models.User
	err := facades.Orm().Query().
		With("CreatedTasks").
		With("AssignedTasks").
		Find(&user, ctx.Request().RouteInt("id"))
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(200, user)
}
