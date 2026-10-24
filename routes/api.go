package routes

import (
	"taskflow/app/facades"
	"taskflow/app/http/controllers"

	"github.com/goravel/framework/contracts/route"
)

func Api() {
	authController := controllers.NewAuthController()
	taskController := controllers.NewTaskController()
	userController := controllers.NewUserController()

	facades.Route().Prefix("api/v1").Group(func(router route.Router) {
		router.Post("register", authController.Register)
		router.Post("login", authController.Login)

		router.Post("tasks", taskController.Store).Name("tasks.create")
		router.Get("tasks/{id}", taskController.Show)
		router.Get("tasks", taskController.Index)

		router.Get("users", userController.Index)
		router.Get("users/{id}", userController.Show)
	})
}
