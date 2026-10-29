package bootstrap

import (
	"github.com/goravel/framework/contracts/console"

	"taskflow/app/console/commands"
)

func Commands() []console.Command {
	return []console.Command{
		&commands.RemindDueTasks{},
		&commands.PruneExports{},
	}
}
