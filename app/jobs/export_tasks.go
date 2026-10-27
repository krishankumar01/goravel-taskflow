package jobs

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/queue"
	frameworkerrors "github.com/goravel/framework/errors"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"
)

type ExportTasks struct {
}

func NewExportTasks(exportID uint) (*ExportTasks, []queue.Arg) {
	return &ExportTasks{}, []queue.Arg{
		{Type: "uint", Value: exportID},
	}
}

// Signature The name and signature of the job.
func (receiver *ExportTasks) Signature() string {
	return "export_tasks"
}

// Handle Execute the job.
func (receiver *ExportTasks) Handle(args ...any) error {
	exportID, ok := args[0].(uint)
	if !ok {
		return errors.New("export tasks: the argument must be the export id")
	}

	var export models.Export
	err := facades.Orm().Query().FindOrFail(&export, exportID)
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	var tasks []models.Task
	if err := facades.Orm().Query().With("Assignee").OrderBy("id").Find(&tasks); err != nil {
		return err
	}

	rows := [][]string{{"id", "title", "status", "assignee", "due_at"}}
	for _, task := range tasks {
		assignee, dueAt := "", ""
		if task.Assignee != nil {
			assignee = task.Assignee.Name
		}
		if task.DueAt != nil {
			dueAt = task.DueAt.ToDateTimeString()
		}
		rows = append(rows, []string{strconv.Itoa(int(task.ID)), task.Title, string(task.Status), assignee, dueAt})
	}

	var file bytes.Buffer
	if err := csv.NewWriter(&file).WriteAll(rows); err != nil {
		return err
	}

	path := fmt.Sprintf("exports/%d.csv", export.ID)
	if err := facades.Storage().Put(path, file.String()); err != nil {
		return err
	}

	export.Status = enums.ExportDone
	export.Path = path

	return facades.Orm().Query().Save(&export)
}

// ShouldRetry gives a failed export two more tries, five seconds apart.
func (receiver *ExportTasks) ShouldRetry(err error, attempt int) (bool, time.Duration) {
	return attempt < 3, 5 * time.Second
}
