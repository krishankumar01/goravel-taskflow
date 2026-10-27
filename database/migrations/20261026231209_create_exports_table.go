package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"taskflow/app/facades"
)

type M20261026231209CreateExportsTable struct{}

// Signature The unique signature for the migration.
func (r *M20261026231209CreateExportsTable) Signature() string {
	return "20261026231209_create_exports_table"
}

// Up Run the migrations.
func (r *M20261026231209CreateExportsTable) Up() error {
	if !facades.Schema().HasTable("exports") {
		return facades.Schema().Create("exports", func(table schema.Blueprint) {
			table.TimestampTz("created_at").Nullable()
			table.TimestampTz("updated_at").Nullable()
			table.BigIncrements("id")
			table.UnsignedBigInteger("user_id")
			table.Enum("status", []any{"pending", "done"}).Default("pending")
			table.String("path", 255)

			table.Foreign("user_id").References("id").On("users")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261026231209CreateExportsTable) Down() error {
	return facades.Schema().DropIfExists("exports")
}
