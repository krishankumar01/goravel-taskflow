package middleware

import (
	"errors"

	"taskflow/app/facades"
	"taskflow/app/http/authuser"
	"taskflow/app/models"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"
)

type Auth struct{}

func (m *Auth) Signature() string {
	return "Auth"
}

func (m *Auth) Handle(ctx http.Context) {
	if _, err := facades.Auth(ctx).Parse(ctx.Request().Header("Authorization")); err != nil {
		ctx.Response().Json(401, map[string]any{
			"error": err.Error(),
		}).Abort()
		return
	}

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
			ctx.Response().Json(401, map[string]any{
				"error": "The user no longer exists.",
			}).Abort()
			return
		}

		facades.Log().Error(err)
		ctx.Response().Json(500, map[string]any{
			"error": "Could not load the user.",
		}).Abort()
		return
	}

	authuser.Set(ctx, &user)
	ctx.Request().Next()
}
