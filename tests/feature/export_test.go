package feature

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goravel/framework/support/http"
	"github.com/stretchr/testify/suite"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"
	"taskflow/tests"
)

type ExportTestSuite struct {
	suite.Suite
	tests.TestCase
	ana models.User
}

func TestExportTestSuite(t *testing.T) {
	suite.Run(t, new(ExportTestSuite))
}

func (s *ExportTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.Require().True(facades.Cache().Flush())
	s.Require().NoError(facades.Storage().DeleteDirectory("exports"))

	var ana, ben models.User
	s.Require().NoError(facades.Orm().Factory().Create(&ana, map[string]any{"Name": "Ana", "Email": "ana@example.com", "Role": enums.RoleManager}))
	s.Require().NoError(facades.Orm().Factory().Create(&ben, map[string]any{"Name": "Ben", "Email": "ben@example.com"}))
	s.Require().NoError(facades.Orm().Factory().Create(&models.Task{}, map[string]any{"Title": "Fix the login page", "Status": enums.StatusPending, "CreatedBy": ana.ID, "AssignedTo": &ben.ID}))

	s.ana = ana
}

func (s *ExportTestSuite) login(email string) string {
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

func (s *ExportTestSuite) TestExportingTasks() {
	ana := s.login("ana@example.com")

	response, err := s.Http(s.T()).WithToken(ana).Post("/api/v1/exports", nil)
	s.Require().NoError(err)

	var export models.Export
	s.Require().NoError(response.AssertStatus(202).Bind(&export))

	s.Run("the job has run, so the export is done", func() {
		response, err := s.Http(s.T()).WithToken(ana).Get(fmt.Sprintf("/api/v1/exports/%d", export.ID))
		s.Require().NoError(err)
		var shown models.Export
		s.Require().NoError(response.AssertOk().Bind(&shown))
		s.Equal(enums.ExportDone, shown.Status)
	})
	s.Run("the download is a CSV of every task", func() {
		response, err := s.Http(s.T()).WithToken(ana).Get(fmt.Sprintf("/api/v1/exports/%d/download", export.ID))
		s.Require().NoError(err)
		response.AssertOk().
			AssertHeader("Content-Type", "text/csv").
			AssertHeader("Content-Disposition", fmt.Sprintf("attachment; filename=tasks-%d.csv", export.ID))
		content, err := response.Content()
		s.Require().NoError(err)
		lines := strings.Split(strings.TrimSpace(content), "\n")
		s.Require().Len(lines, 2)
		s.Equal("id,title,status,assignee,due_at", lines[0])
		s.Contains(lines[1], "Fix the login page,pending,Ben,")
	})
	s.Run("another user can't see it", func() {
		response, err := s.Http(s.T()).WithToken(s.login("ben@example.com")).Get(fmt.Sprintf("/api/v1/exports/%d", export.ID))
		s.Require().NoError(err)
		response.AssertNotFound()
	})
}

func (s *ExportTestSuite) TestAnExportThatIsNotReadyCannotBeDownloaded() {
	export := models.Export{UserID: s.ana.ID, Status: enums.ExportPending}
	s.Require().NoError(facades.Orm().Query().Create(&export))

	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).Get(fmt.Sprintf("/api/v1/exports/%d/download", export.ID))
	s.Require().NoError(err)
	response.AssertStatus(409)
}
