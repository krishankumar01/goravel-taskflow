package policies

import (
	"context"

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

func (r *TaskPolicy) Delete(ctx context.Context, arguments map[string]any) contractsaccess.Response {
	return access.NewDenyResponse("Only managers can delete tasks.")
}
