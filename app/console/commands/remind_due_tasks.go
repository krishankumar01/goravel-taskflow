package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support/carbon"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/mails"
	"taskflow/app/models"
)

type RemindDueTasks struct {
}

func (r *RemindDueTasks) Signature() string {
	return "tasks:remind"
}

func (r *RemindDueTasks) Description() string {
	return "Email each assignee their open tasks that are due soon"
}

func (r *RemindDueTasks) Extend() command.Extend {
	return command.Extend{
		Category: "tasks",
		Flags: []command.Flag{
			&command.IntFlag{
				Name:  "hours",
				Value: 24,
				Usage: "remind about tasks due within this many hours",
			},
		},
	}
}

func (r *RemindDueTasks) Handle(ctx console.Context) error {
	hours := ctx.OptionInt("hours")

	var tasks []models.Task
	err := facades.Orm().Query().With("Assignee").
		Where("assigned_to IS NOT NULL").
		Where("status <> ?", enums.StatusDone).
		Where("due_at BETWEEN ? AND ?", carbon.Now().StdTime(), carbon.Now().AddHours(hours).StdTime()).
		OrderBy("due_at").
		Find(&tasks)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		ctx.Info("No tasks are due soon.")
		return nil
	}

	var assignees []uint
	due := map[uint][]models.Task{}
	for _, task := range tasks {
		id := *task.AssignedTo
		if _, ok := due[id]; !ok {
			assignees = append(assignees, id)
		}
		due[id] = append(due[id], task)
	}

	var rows [][]string
	for _, id := range assignees {
		assignee := *due[id][0].Assignee
		if err := facades.Mail().Send(mails.NewDueSoon(assignee, due[id])); err != nil {
			return err
		}
		rows = append(rows, []string{assignee.Email, fmt.Sprint(len(due[id]))})
	}

	ctx.Table([]string{"Assignee", "Tasks"}, rows)

	return nil
}
