package feature

import (
	"testing"

	"github.com/goravel/framework/support/http"
	"github.com/stretchr/testify/suite"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"
	"taskflow/tests"
)

type MailTestSuite struct {
	suite.Suite
	tests.TestCase
	ana models.User
	ben models.User
}

func TestMailTestSuite(t *testing.T) {
	suite.Run(t, new(MailTestSuite))
}

func (s *MailTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.Require().True(facades.Cache().Flush())
	s.Require().NoError(facades.Storage().DeleteDirectory("exports"))
	s.Require().NoError(tests.ClearMail())

	var ana, ben models.User
	s.Require().NoError(facades.Orm().Factory().Create(&ana, map[string]any{"Name": "Ana", "Email": "ana@example.com", "Role": enums.RoleManager}))
	s.Require().NoError(facades.Orm().Factory().Create(&ben, map[string]any{"Name": "Ben", "Email": "ben@example.com"}))

	s.ana, s.ben = ana, ben
}

func (s *MailTestSuite) login(email string) string {
	body, err := http.NewBody().SetFields(map[string]any{"email": email, "password": "password"}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithHeader("Content-Type", body.ContentType()).Post("/api/v1/login", body.Reader())
	s.Require().NoError(err)
	response.AssertOk()

	var login struct {
		Token string `json:"token"`
	}
	s.Require().NoError(response.Bind(&login))

	return login.Token
}

func (s *MailTestSuite) createTask(token string, assignedTo uint) {
	body, err := http.NewBody().SetFields(map[string]any{
		"title":       "Fix the login page",
		"description": "Assigned in a test.",
		"assigned_to": assignedTo,
	}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithToken(token).WithHeader("Content-Type", body.ContentType()).Post("/api/v1/tasks", body.Reader())
	s.Require().NoError(err)
	response.AssertCreated()
}

func (s *MailTestSuite) TestAssigningEmailsTheAssignee() {
	ana := s.login("ana@example.com")

	s.Run("the assignee gets an email", func() {
		s.createTask(ana, s.ben.ID)

		mails, err := tests.MailTo("ben@example.com")
		s.Require().NoError(err)
		s.Require().Len(mails, 1)
		s.Equal("Ana assigned you a task", mails[0].Subject)
		s.Contains(mails[0].Text, `Ana assigned you "Fix the login page".`)
	})
	s.Run("a manager who takes a task themselves gets no email", func() {
		s.createTask(ana, s.ana.ID)

		mails, err := tests.MailTo("ana@example.com")
		s.Require().NoError(err)
		s.Empty(mails)
	})
}

func (s *MailTestSuite) TestAFinishedExportEmailsItsOwner() {
	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).Post("/api/v1/exports", nil)
	s.Require().NoError(err)
	response.AssertStatus(202)

	mails, err := tests.MailTo("ana@example.com")
	s.Require().NoError(err)
	s.Require().Len(mails, 1)
	s.Equal("Your task export is ready", mails[0].Subject)
	s.Contains(mails[0].Text, "GET /api/v1/exports/1/download")
}
