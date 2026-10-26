package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"taskflow/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20261022221906CreateUsersTable{},
		&migrations.M20261023063944CreateTasksTable{},
		&migrations.M20261026070447CreateAttachmentsTable{},
	}
}
