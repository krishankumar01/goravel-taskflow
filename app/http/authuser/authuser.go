package authuser

import (
	"context"

	"taskflow/app/models"

	"github.com/goravel/framework/contracts/http"
)

type key struct{}

func Set(ctx http.Context, user *models.User) {
	ctx.WithValue(key{}, user)
}

func From(ctx context.Context) *models.User {
	user, _ := ctx.Value(key{}).(*models.User)

	return user
}
