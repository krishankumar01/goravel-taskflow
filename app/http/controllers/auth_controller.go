package controllers

import (
	"sync"

	"taskflow/app/facades"
	"taskflow/app/http/requests"
	"taskflow/app/models"

	"github.com/goravel/framework/contracts/http"
)

// noUserHash is a bcrypt hash of a password nobody has, made with the same cost as the
// real ones. An unknown email is checked against it, so it takes as long to refuse as a
// wrong password.
var noUserHash = sync.OnceValue(func() string {
	hash, err := facades.Hash().Make("no user has this password")
	if err != nil {
		panic(err)
	}

	return hash
})

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

func (r *AuthController) Login(ctx http.Context) http.Response {
	var request requests.LoginRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not log in.",
		})
	}
	if validationErrors != nil {
		return ctx.Response().Json(422, map[string]any{
			"errors": validationErrors.All(),
		})
	}

	var user models.User
	if err := facades.Orm().WithContext(ctx).Query().Where("email", request.Email).First(&user); err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not log in.",
		})
	}

	hash := user.Password
	if user.ID == 0 {
		hash = noUserHash()
	}
	if !facades.Hash().Check(request.Password, hash) || user.ID == 0 {
		return ctx.Response().Json(401, map[string]any{
			"error": "The email or password is wrong.",
		})
	}

	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": "Could not log in.",
		})
	}

	return ctx.Response().Json(200, map[string]any{
		"token": token,
	})
}
