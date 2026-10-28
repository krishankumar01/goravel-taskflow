package listeners

import (
	"errors"

	"github.com/goravel/framework/contracts/event"
	frameworkerrors "github.com/goravel/framework/errors"

	"taskflow/app/facades"
	"taskflow/app/mails"
	"taskflow/app/models"
)

type EmailAssignee struct {
}

func (receiver *EmailAssignee) Signature() string {
	return "email_assignee"
}

func (receiver *EmailAssignee) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *EmailAssignee) Handle(args ...any) error {
	taskID, ok1 := args[0].(uint)
	assignedTo, ok2 := args[1].(uint)
	assignedBy, ok3 := args[2].(uint)
	if !ok1 || !ok2 || !ok3 {
		return errors.New("email assignee: the ids must be uints")
	}
	if assignedTo == assignedBy {
		return nil
	}

	var task models.Task
	err := facades.Orm().Query().FindOrFail(&task, taskID)
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	var assignee, manager models.User
	if err := facades.Orm().Query().FindOrFail(&assignee, assignedTo); err != nil {
		return err
	}
	if err := facades.Orm().Query().FindOrFail(&manager, assignedBy); err != nil {
		return err
	}

	return facades.Mail().Send(mails.NewAssignment(task, assignee, manager))
}
