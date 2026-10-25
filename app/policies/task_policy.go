package policies

import (
	"context"

	"taskflow/app/http/authuser"
	"taskflow/app/models"

	"github.com/goravel/framework/auth/access"
	contractsaccess "github.com/goravel/framework/contracts/auth/access"
)

type TaskPolicy struct {
}

func NewTaskPolicy() *TaskPolicy {
	return &TaskPolicy{}
}

func (r *TaskPolicy) Create(ctx context.Context, arguments map[string]any) contractsaccess.Response {
	return access.NewDenyResponse("Only managers can create tasks.")
}

func (r *TaskPolicy) Update(ctx context.Context, arguments map[string]any) contractsaccess.Response {
	user := authuser.From(ctx)
	task := arguments["task"].(*models.Task)
	changes := arguments["changes"].(map[string]any)

	if user == nil || task.AssignedTo == nil || *task.AssignedTo != user.ID {
		return access.NewDenyResponse("You can only change tasks assigned to you.")
	}
	for field := range changes {
		if field != "status" {
			return access.NewDenyResponse("You can only change the status of a task.")
		}
	}

	return access.NewAllowResponse()
}

func (r *TaskPolicy) Delete(ctx context.Context, arguments map[string]any) contractsaccess.Response {
	return access.NewDenyResponse("Only managers can delete tasks.")
}
