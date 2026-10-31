package controllers

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/http"
	redisfacades "github.com/goravel/redis/facades"

	"taskflow/app/facades"
)

type HealthController struct {
	// Dependent services
}

func NewHealthController() *HealthController {
	return &HealthController{
		// Inject services
	}
}

func (r *HealthController) Show(ctx http.Context) http.Response {
	c, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	checks := map[string]string{"database": "ok", "redis": "ok"}
	status := http.StatusOK

	db, err := facades.Orm().DB()
	if err == nil {
		err = db.PingContext(c)
	}
	if err != nil {
		checks["database"] = err.Error()
		status = http.StatusServiceUnavailable
	}

	client, err := redisfacades.Instance("default")
	if err == nil {
		err = client.Ping(c).Err()
	}
	if err != nil {
		checks["redis"] = err.Error()
		status = http.StatusServiceUnavailable
	}

	return ctx.Response().Json(status, checks)
}
