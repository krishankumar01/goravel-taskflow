package middleware

import (
	"taskflow/app/facades"

	"github.com/goravel/framework/contracts/http"
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

	ctx.Request().Next()
}
