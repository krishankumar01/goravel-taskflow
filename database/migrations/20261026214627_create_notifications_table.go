package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"taskflow/app/facades"
)

type M20261026214627CreateNotificationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20261026214627CreateNotificationsTable) Signature() string {
	return "20261026214627_create_notifications_table"
}

// Up Run the migrations.
func (r *M20261026214627CreateNotificationsTable) Up() error {
	if !facades.Schema().HasTable("notifications") {
		return facades.Schema().Create("notifications", func(table schema.Blueprint) {
			table.TimestampTz("created_at").Nullable()
			table.TimestampTz("updated_at").Nullable()
			table.BigIncrements("id")
			table.UnsignedBigInteger("user_id")
			table.UnsignedBigInteger("task_id")
			table.String("message", 255)
			table.TimestampTz("read_at").Nullable()

			table.Foreign("user_id").References("id").On("users")
			table.Foreign("task_id").References("id").On("tasks").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261026214627CreateNotificationsTable) Down() error {
	return facades.Schema().DropIfExists("notifications")
}
