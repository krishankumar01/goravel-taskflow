package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"taskflow/app/facades"
)

type M20261022221906CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20261022221906CreateUsersTable) Signature() string {
	return "20261022221906_create_users_table"
}

// Up Run the migrations.
func (r *M20261022221906CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.TimestampTz("created_at").Nullable()
			table.TimestampTz("updated_at").Nullable()
			table.BigIncrements("id")
			table.String("name", 255)
			table.String("email", 255)
			table.String("password", 255)
			table.Enum("role", []any{"manager", "member"}).Default("member")

			table.Unique("email")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261022221906CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
