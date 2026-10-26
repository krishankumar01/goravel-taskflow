package listeners

import (
	"errors"
	"fmt"

	"github.com/goravel/framework/contracts/event"
	frameworkerrors "github.com/goravel/framework/errors"

	"taskflow/app/facades"
	"taskflow/app/models"
)

type NotifyAssignee struct {
}

func (receiver *NotifyAssignee) Signature() string {
	return "notify_assignee"
}

func (receiver *NotifyAssignee) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *NotifyAssignee) Handle(args ...any) error {
	taskID, ok1 := args[0].(uint)
	assignedTo, ok2 := args[1].(uint)
	assignedBy, ok3 := args[2].(uint)
	if !ok1 || !ok2 || !ok3 {
		return errors.New("notify assignee: the ids must be uints")
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

	var manager models.User
	if err := facades.Orm().Query().FindOrFail(&manager, assignedBy); err != nil {
		return err
	}

	return facades.Orm().Query().Create(&models.Notification{
		UserID:  assignedTo,
		TaskID:  task.ID,
		Message: fmt.Sprintf("%s assigned you %q", manager.Name, task.Title),
	})
}
