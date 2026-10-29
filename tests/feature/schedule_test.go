package feature

import (
	"testing"

	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"
	"taskflow/tests"
)

type ScheduleTestSuite struct {
	suite.Suite
	tests.TestCase
	ana models.User
	ben models.User
}

func TestScheduleTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleTestSuite))
}

func (s *ScheduleTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.Require().NoError(facades.Storage().DeleteDirectory("exports"))
	s.Require().NoError(tests.ClearMail())
	carbon.SetTestNow(carbon.Parse("2026-11-02 08:00:00"))

	var ana, ben models.User
	s.Require().NoError(facades.Orm().Factory().Create(&ana, map[string]any{"Name": "Ana", "Email": "ana@example.com", "Role": enums.RoleManager}))
	s.Require().NoError(facades.Orm().Factory().Create(&ben, map[string]any{"Name": "Ben", "Email": "ben@example.com"}))
	s.ana, s.ben = ana, ben
}

func (s *ScheduleTestSuite) TearDownTest() {
	carbon.ClearTestNow()
}

func (s *ScheduleTestSuite) task(title, due string, status enums.Status) {
	dueAt := carbon.NewDateTime(carbon.Parse(due))
	s.Require().NoError(facades.Orm().Factory().Create(&models.Task{}, map[string]any{"Title": title, "Status": status, "CreatedBy": s.ana.ID, "AssignedTo": &s.ben.ID, "DueAt": dueAt}))
}

func (s *ScheduleTestSuite) TestTheReminderListsOpenTasksDueSoon() {
	s.task("Due tomorrow", "2026-11-03 07:00:00", enums.StatusPending)
	s.task("Already done", "2026-11-02 12:00:00", enums.StatusDone)
	s.task("Due next week", "2026-11-09 09:00:00", enums.StatusPending)

	s.Require().NoError(facades.Artisan().Call("tasks:remind"))

	mails, err := tests.MailTo("ben@example.com")
	s.Require().NoError(err)
	s.Require().Len(mails, 1)
	s.Equal("Your tasks due soon", mails[0].Subject)
	s.Contains(mails[0].Text, "Due tomorrow")
	s.NotContains(mails[0].Text, "Already done")
	s.NotContains(mails[0].Text, "Due next week")
}

func (s *ScheduleTestSuite) TestHoursSetsTheWindow() {
	s.task("Due tomorrow", "2026-11-03 07:00:00", enums.StatusPending)

	s.Require().NoError(facades.Artisan().Call("tasks:remind --hours 2"))

	mails, err := tests.MailTo("ben@example.com")
	s.Require().NoError(err)
	s.Empty(mails)
}

func (s *ScheduleTestSuite) TestPruneDeletesOldAndStuckExports() {
	old := models.Export{UserID: s.ana.ID, Status: enums.ExportDone, Path: "exports/1.csv"}
	recent := models.Export{UserID: s.ana.ID, Status: enums.ExportDone, Path: "exports/2.csv"}
	stuck := models.Export{UserID: s.ana.ID, Status: enums.ExportPending}
	running := models.Export{UserID: s.ana.ID, Status: enums.ExportPending}
	for _, export := range []*models.Export{&old, &recent, &stuck, &running} {
		s.Require().NoError(facades.Orm().Query().Create(export))
	}
	s.Require().NoError(facades.Storage().Put(old.Path, "id,title"))
	s.Require().NoError(facades.Storage().Put(recent.Path, "id,title"))
	_, err := facades.Orm().Query().Model(&old).Update("created_at", carbon.Now().SubDays(8).StdTime())
	s.Require().NoError(err)
	_, err = facades.Orm().Query().Model(&stuck).Update("created_at", carbon.Now().SubHours(2).StdTime())
	s.Require().NoError(err)

	s.Require().NoError(facades.Artisan().Call("exports:prune"))

	var left []models.Export
	s.Require().NoError(facades.Orm().Query().OrderBy("id").Find(&left))
	s.Require().Len(left, 2)
	s.Equal(recent.ID, left[0].ID)
	s.Equal(running.ID, left[1].ID)
	s.False(facades.Storage().Exists(old.Path))
	s.True(facades.Storage().Exists(recent.Path))
}

func (s *ScheduleTestSuite) TestTheSchedule() {
	crons := map[string]string{}
	for _, event := range facades.Schedule().Events() {
		crons[event.GetCommand()] = event.GetCron()
	}
	s.Equal(map[string]string{"tasks:remind": "00 09 * * *", "exports:prune": "0 0 * * *"}, crons)
}
