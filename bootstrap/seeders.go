package bootstrap

import (
	"github.com/goravel/framework/contracts/database/seeder"

	"taskflow/database/seeders"
)

func Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.UserSeeder{},
		&seeders.TaskSeeder{},
	}
}
