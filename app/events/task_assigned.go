package events

import (
	"errors"

	"github.com/goravel/framework/contracts/event"
)

type TaskAssigned struct {
}

func NewTaskAssigned(taskID, assignedTo, assignedBy uint) (*TaskAssigned, []event.Arg) {
	return &TaskAssigned{}, []event.Arg{
		{Type: "uint", Value: taskID},
		{Type: "uint", Value: assignedTo},
		{Type: "uint", Value: assignedBy},
	}
}

func (receiver *TaskAssigned) Handle(args []event.Arg) ([]event.Arg, error) {
	if len(args) != 3 {
		return nil, errors.New("task assigned: needs the task id, the assignee id and the id of who assigned it")
	}
	for _, arg := range args {
		if _, ok := arg.Value.(uint); !ok {
			return nil, errors.New("task assigned: every id must be a uint")
		}
	}

	return args, nil
}
