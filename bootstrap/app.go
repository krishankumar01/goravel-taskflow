package bootstrap

import (
	"context"
	"strings"

	"github.com/goravel/framework/auth/access"
	contractsaccess "github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/database/schema"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/foundation"
	"github.com/goravel/framework/http/limit"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/http/authuser"
	"taskflow/app/models"
	"taskflow/app/policies"
	"taskflow/config"
	"taskflow/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithRules(Rules).
		WithSeeders(Seeders).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Api()
			routes.Web()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		WithCallback(func() {
			facades.Schema().Extend(schema.Extension{
				Models: []any{models.User{}, models.Task{}},
			})

			facades.RateLimiter().For("login", func(ctx http.Context) http.Limit {
				email := strings.ToLower(strings.TrimSpace(ctx.Request().Input("email")))

				return limit.PerMinute(5).By(email).Response(func(ctx http.Context) {
					ctx.Response().Json(429, map[string]any{
						"error": "Too many login attempts. Try again in a minute.",
					}).Abort()
				})
			})

			facades.Gate().Before(func(ctx context.Context, ability string, arguments map[string]any) contractsaccess.Response {
				if user := authuser.From(ctx); user != nil && user.Role == enums.RoleManager {
					return access.NewAllowResponse()
				}

				return nil
			})

			facades.Gate().Define(policies.ViewUsers, func(ctx context.Context, arguments map[string]any) contractsaccess.Response {
				return access.NewDenyResponse("Only managers can see the users.")
			})
		}).
		Create()
}
