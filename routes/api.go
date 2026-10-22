package routes

import (
	"taskflow/app/facades"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
)

func Api() {
	facades.Route().Prefix("api/v1").Group(func(router route.Router) {
		router.Post("tasks", func(ctx http.Context) http.Response {
			return ctx.Response().String(200, "task created")
		}).Name("tasks.create")
	})
}
