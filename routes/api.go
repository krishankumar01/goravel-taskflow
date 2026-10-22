package routes

import (
	"taskflow/app/facades"
	"taskflow/app/http/controllers"

	"github.com/goravel/framework/contracts/route"
)

func Api() {
	taskController := controllers.NewTaskController()

	facades.Route().Prefix("api/v1").Group(func(router route.Router) {
		router.Post("tasks", taskController.Store).Name("tasks.create")
		router.Get("tasks/{id}", taskController.Show)
	})
}
