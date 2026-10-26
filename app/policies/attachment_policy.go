package policies

import (
	"context"

	"taskflow/app/http/authuser"
	"taskflow/app/models"

	"github.com/goravel/framework/auth/access"
	contractsaccess "github.com/goravel/framework/contracts/auth/access"
)

type AttachmentPolicy struct {
}

func NewAttachmentPolicy() *AttachmentPolicy {
	return &AttachmentPolicy{}
}

func (r *AttachmentPolicy) Create(ctx context.Context, arguments map[string]any) contractsaccess.Response {
	user := authuser.From(ctx)
	task := arguments["task"].(*models.Task)

	if user == nil || task.AssignedTo == nil || *task.AssignedTo != user.ID {
		return access.NewDenyResponse("You can only attach files to tasks assigned to you.")
	}

	return access.NewAllowResponse()
}

func (r *AttachmentPolicy) View(ctx context.Context, arguments map[string]any) contractsaccess.Response {
	user := authuser.From(ctx)
	task := arguments["task"].(*models.Task)

	if user == nil || task.AssignedTo == nil || *task.AssignedTo != user.ID {
		return access.NewDenyResponse("You can only open files on tasks assigned to you.")
	}

	return access.NewAllowResponse()
}

func (r *AttachmentPolicy) Delete(ctx context.Context, arguments map[string]any) contractsaccess.Response {
	user := authuser.From(ctx)
	attachment := arguments["attachment"].(*models.Attachment)

	if user == nil || attachment.UploadedBy != user.ID {
		return access.NewDenyResponse("You can only delete files you uploaded.")
	}

	return access.NewAllowResponse()
}
