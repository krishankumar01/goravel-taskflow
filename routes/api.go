package routes

import (
	"taskflow/app/facades"
	"taskflow/app/http/controllers"
	"taskflow/app/http/middleware"

	"github.com/goravel/framework/contracts/route"
	httpmiddleware "github.com/goravel/framework/http/middleware"
)

func Api() {
	authController := controllers.NewAuthController()
	taskController := controllers.NewTaskController()
	userController := controllers.NewUserController()

	facades.Route().Prefix("api/v1").Group(func(router route.Router) {
		router.Post("register", authController.Register)
		router.Middleware(httpmiddleware.Throttle("login")).Post("login", authController.Login)
		router.Post("refresh", authController.Refresh)

		router.Middleware(&middleware.Auth{}).Group(func(router route.Router) {
			router.Get("me", authController.Me)
			router.Post("logout", authController.Logout)

			router.Post("tasks", taskController.Store).Name("tasks.create")
			router.Get("tasks/{id}", taskController.Show)
			router.Get("tasks", taskController.Index)

			router.Get("users", userController.Index)
			router.Get("users/{id}", userController.Show)
		})
	})
}
