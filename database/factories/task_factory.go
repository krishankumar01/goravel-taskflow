package factories

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/goravel/framework/support/carbon"

	"taskflow/app/enums"
)

type TaskFactory struct {
}

// Definition Define the model's default state.
func (f *TaskFactory) Definition() map[string]any {
	title := gofakeit.RandomString([]string{
		"Fix the login bug",
		"Write the API docs",
		"Review the pull request",
		"Plan the next release",
		"Update the dependencies",
		"Add tests for tasks",
		"Clean up old logs",
		"Reply to support emails",
	})
	statuses := []enums.Status{enums.StatusPending, enums.StatusInProgress, enums.StatusDone}
	status := statuses[gofakeit.Number(0, len(statuses)-1)]
	due := carbon.Now().AddDays(gofakeit.Number(1, 30)).StartOfDay().SetHour(9)

	return map[string]any{
		"Title":       title,
		"Description": gofakeit.HackerPhrase(),
		"Status":      status,
		"DueAt":       carbon.NewDateTime(due),
	}
}
