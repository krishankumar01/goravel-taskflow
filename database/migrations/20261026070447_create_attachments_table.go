package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"taskflow/app/facades"
)

type M20261026070447CreateAttachmentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20261026070447CreateAttachmentsTable) Signature() string {
	return "20261026070447_create_attachments_table"
}

// Up Run the migrations.
func (r *M20261026070447CreateAttachmentsTable) Up() error {
	if !facades.Schema().HasTable("attachments") {
		return facades.Schema().Create("attachments", func(table schema.Blueprint) {
			table.TimestampTz("created_at").Nullable()
			table.TimestampTz("updated_at").Nullable()
			table.BigIncrements("id")
			table.UnsignedBigInteger("task_id")
			table.UnsignedBigInteger("uploaded_by")
			table.String("name", 255)
			table.String("path", 255)
			table.BigInteger("size")
			table.String("mime_type", 255)

			table.Foreign("task_id").References("id").On("tasks").CascadeOnDelete()
			table.Foreign("uploaded_by").References("id").On("users")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261026070447CreateAttachmentsTable) Down() error {
	return facades.Schema().DropIfExists("attachments")
}
