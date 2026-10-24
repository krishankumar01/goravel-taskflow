package controllers

import (
	"taskflow/app/facades"
	"taskflow/app/http/requests"
	"taskflow/app/models"

	"github.com/goravel/framework/contracts/http"
)

type AuthController struct {
	// Dependent services
}

func NewAuthController() *AuthController {
	return &AuthController{
		// Inject services
	}
}

func (r *AuthController) Register(ctx http.Context) http.Response {
	var request requests.RegisterRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not create the account.",
		})
	}
	if validationErrors != nil {
		return ctx.Response().Json(422, map[string]any{
			"errors": validationErrors.All(),
		})
	}

	password, err := facades.Hash().Make(request.Password)
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not create the account.",
		})
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: password,
	}
	if err := facades.Orm().WithContext(ctx).Query().Create(&user); err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not create the account.",
		})
	}

	return ctx.Response().Json(201, user)
}
