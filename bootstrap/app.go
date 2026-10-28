package bootstrap

import (
	"context"
	"strings"

	"github.com/goravel/framework/auth/access"
	contractsaccess "github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/database/schema"
	contractsevent "github.com/goravel/framework/contracts/event"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/foundation"
	"github.com/goravel/framework/http/limit"

	"taskflow/app/enums"
	"taskflow/app/events"
	"taskflow/app/facades"
	"taskflow/app/http/authuser"
	"taskflow/app/listeners"
	"taskflow/app/models"
	"taskflow/app/policies"
	"taskflow/config"
	"taskflow/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithJobs(Jobs).
		WithRules(Rules).
		WithSeeders(Seeders).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Api()
			routes.Web()
		}).
		WithEvents(func() map[contractsevent.Event][]contractsevent.Listener {
			return map[contractsevent.Event][]contractsevent.Listener{
				&events.TaskAssigned{}: {&listeners.NotifyAssignee{}, &listeners.EmailAssignee{}},
			}
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		WithCallback(func() {
			facades.Schema().Extend(schema.Extension{
				Models: []any{models.User{}, models.Task{}, models.Attachment{}, models.Notification{}, models.Export{}},
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

			taskPolicy := policies.NewTaskPolicy()
			facades.Gate().Define(policies.CreateTask, taskPolicy.Create)
			facades.Gate().Define(policies.UpdateTask, taskPolicy.Update)
			facades.Gate().Define(policies.DeleteTask, taskPolicy.Delete)

			attachmentPolicy := policies.NewAttachmentPolicy()
			facades.Gate().Define(policies.CreateAttachment, attachmentPolicy.Create)
			facades.Gate().Define(policies.ViewAttachment, attachmentPolicy.View)
			facades.Gate().Define(policies.DeleteAttachment, attachmentPolicy.Delete)
		}).
		Create()
}
