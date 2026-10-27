package bootstrap

import (
	"github.com/goravel/framework/contracts/queue"

	"taskflow/app/jobs"
)

func Jobs() []queue.Job {
	return []queue.Job{
		&jobs.ExportTasks{},
	}
}
