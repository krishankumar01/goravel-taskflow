package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/support"

	"taskflow/app/http/controllers"
	"taskflow/app/facades"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Static("public", "./public")

	userController := controllers.NewUserController()
	facades.Route().Get("/users", userController.Index)

	facades.Route().Get("hello", func(ctx http.Context) http.Response {
		return ctx.Response().Json(200, map[string]any{
			"message": "Hello, world",
		})
	})

	facades.Route().Prefix("api/v1").Group(func(router route.Router) {
		router.Post("tasks", func(ctx http.Context) http.Response {
			return ctx.Response().String(200, "task created")
		})
	})
}
