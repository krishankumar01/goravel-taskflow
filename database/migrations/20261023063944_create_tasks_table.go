package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"taskflow/app/facades"
)

type M20261023063944CreateTasksTable struct{}

// Signature The unique signature for the migration.
func (r *M20261023063944CreateTasksTable) Signature() string {
	return "20261023063944_create_tasks_table"
}

// Up Run the migrations.
func (r *M20261023063944CreateTasksTable) Up() error {
	if !facades.Schema().HasTable("tasks") {
		return facades.Schema().Create("tasks", func(table schema.Blueprint) {
			table.TimestampTz("created_at").Nullable()
			table.TimestampTz("updated_at").Nullable()
			table.BigIncrements("id")
			table.UnsignedBigInteger("created_by")
			table.UnsignedBigInteger("assigned_to").Nullable()
			table.String("title", 255)
			table.Text("description")
			table.Enum("status", []any{"pending", "in-progress", "done"}).Default("pending")
			table.TimestampTz("due_at").Nullable()

			table.Foreign("created_by").References("id").On("users")
			table.Foreign("assigned_to").References("id").On("users").NullOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261023063944CreateTasksTable) Down() error {
	return facades.Schema().DropIfExists("tasks")
}
