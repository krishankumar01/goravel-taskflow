package feature

import (
	"fmt"
	"testing"

	contractshttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/goravel/framework/support/http"
	"github.com/stretchr/testify/suite"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"
	"taskflow/tests"
)

type TaskTestSuite struct {
	suite.Suite
	tests.TestCase
	ana     models.User
	ben     models.User
	benTask models.Task
	catTask models.Task
}

func TestTaskTestSuite(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}

func (s *TaskTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.Require().True(facades.Cache().Flush())

	var ana, ben, cat models.User
	s.Require().NoError(facades.Orm().Factory().Create(&ana, map[string]any{"Email": "ana@example.com", "Role": enums.RoleManager}))
	s.Require().NoError(facades.Orm().Factory().Create(&ben, map[string]any{"Email": "ben@example.com"}))
	s.Require().NoError(facades.Orm().Factory().Create(&cat, map[string]any{"Email": "cat@example.com"}))

	var benTask, catTask models.Task
	s.Require().NoError(facades.Orm().Factory().Create(&benTask, map[string]any{"CreatedBy": ana.ID, "AssignedTo": &ben.ID}))
	s.Require().NoError(facades.Orm().Factory().Create(&catTask, map[string]any{"CreatedBy": ana.ID, "AssignedTo": &cat.ID}))

	s.ana, s.ben, s.benTask, s.catTask = ana, ben, benTask, catTask
}

func (s *TaskTestSuite) login(email string) string {
	body, err := http.NewBody().SetFields(map[string]any{"email": email, "password": "password"}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithHeader("Content-Type", body.ContentType()).Post("/api/v1/login", body.Reader())
	s.Require().NoError(err)
	response.AssertOk()

	var login struct {
		Token string `json:"token"`
	}
	s.Require().NoError(response.Bind(&login))
	s.Require().NotEmpty(login.Token)

	return login.Token
}

func (s *TaskTestSuite) TestRoles() {
	ana := s.login("ana@example.com")
	ben := s.login("ben@example.com")
	benTask := fmt.Sprintf("/api/v1/tasks/%d", s.benTask.ID)
	catTask := fmt.Sprintf("/api/v1/tasks/%d", s.catTask.ID)
	newTask := map[string]any{"title": "Write the tests", "description": "Cover every role rule."}

	cases := []struct {
		name   string
		token  string
		method string
		path   string
		body   map[string]any
		status int
	}{
		{"a member reads tasks", ben, "GET", "/api/v1/tasks", nil, 200},
		{"a member can't list users", ben, "GET", "/api/v1/users", nil, 403},
		{"a manager lists users", ana, "GET", "/api/v1/users", nil, 200},
		{"a member can't create a task", ben, "POST", "/api/v1/tasks", newTask, 403},
		{"a manager creates a task", ana, "POST", "/api/v1/tasks", newTask, 201},
		{"a member changes the status of his task", ben, "PATCH", benTask, map[string]any{"status": "done"}, 200},
		{"a member can't rename his task", ben, "PATCH", benTask, map[string]any{"title": "Renamed"}, 403},
		{"a member can't change someone else's task", ben, "PATCH", catTask, map[string]any{"status": "done"}, 403},
		{"a member reads a task", ben, "GET", benTask, nil, 200},
		{"a manager reads a task", ana, "GET", benTask, nil, 200},
		{"a member can't see a user", ben, "GET", fmt.Sprintf("/api/v1/users/%d", s.ana.ID), nil, 403},
		{"a manager sees a user", ana, "GET", fmt.Sprintf("/api/v1/users/%d", s.ben.ID), nil, 200},
		{"a member can't reassign his task", ben, "PATCH", benTask, map[string]any{"assigned_to": s.ana.ID}, 403},
		{"a manager changes any task", ana, "PATCH", catTask, map[string]any{"title": "Renamed"}, 200},
		{"a member can't delete a task", ben, "DELETE", catTask, nil, 403},
		{"a manager deletes a task", ana, "DELETE", catTask, nil, 204},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			body, err := http.NewBody().SetFields(c.body).Build()
			s.Require().NoError(err)
			request := s.Http(s.T()).WithToken(c.token).WithHeader("Content-Type", body.ContentType())

			var response contractshttp.Response
			switch c.method {
			case "GET":
				response, err = request.Get(c.path)
			case "POST":
				response, err = request.Post(c.path, body.Reader())
			case "PATCH":
				response, err = request.Patch(c.path, body.Reader())
			case "DELETE":
				response, err = request.Delete(c.path, nil)
			}
			s.Require().NoError(err)
			response.AssertStatus(c.status)
		})
	}
}

func (s *TaskTestSuite) TestCreateSetsTheCreator() {
	body, err := http.NewBody().SetFields(map[string]any{
		"title":       "Write the tests",
		"description": "Cover every role rule.",
		"assigned_to": s.ben.ID,
	}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).WithHeader("Content-Type", body.ContentType()).Post("/api/v1/tasks", body.Reader())
	s.Require().NoError(err)

	response.AssertCreated()

	var task models.Task
	s.Require().NoError(response.Bind(&task))
	s.Equal("Write the tests", task.Title)
	s.Equal(s.ana.ID, task.CreatedBy)
	s.Equal(s.ben.ID, *task.AssignedTo)
	s.Equal(enums.StatusPending, task.Status)
}

func (s *TaskTestSuite) TestAMissingTaskIsNotFound() {
	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).Get("/api/v1/tasks/999")
	s.Require().NoError(err)
	response.AssertNotFound()
}

func (s *TaskTestSuite) TestBadInputIsRefused() {
	body, err := http.NewBody().SetFields(map[string]any{"description": "No title.", "due_at": "yesterday"}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).WithHeader("Content-Type", body.ContentType()).Post("/api/v1/tasks", body.Reader())
	s.Require().NoError(err)

	response.AssertUnprocessableEntity()

	var refused struct {
		Errors map[string]any `json:"errors"`
	}
	s.Require().NoError(response.Bind(&refused))
	s.Contains(refused.Errors, "title")
	s.Contains(refused.Errors, "due_at")
}

func (s *TaskTestSuite) TestUpdateChangesOnlyTheSentFields() {
	path := fmt.Sprintf("/api/v1/tasks/%d", s.benTask.ID)
	body, err := http.NewBody().SetFields(map[string]any{"status": "done"}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).WithHeader("Content-Type", body.ContentType()).Patch(path, body.Reader())
	s.Require().NoError(err)

	response.AssertOk().AssertFluentJson(func(json contractshttp.AssertableJSON) {
		json.Where("status", "done").
			Where("title", s.benTask.Title).
			Where("description", s.benTask.Description)
	})
}

func (s *TaskTestSuite) TestDeleteRemovesTheTask() {
	path := fmt.Sprintf("/api/v1/tasks/%d", s.catTask.ID)
	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).Delete(path, nil)
	s.Require().NoError(err)
	response.AssertNoContent()

	count, err := facades.Orm().Query().Model(&models.Task{}).Where("id", s.catTask.ID).Count()
	s.Require().NoError(err)
	s.Equal(int64(0), count)
}

func (s *TaskTestSuite) TestLogoutEndsTheToken() {
	token := s.login("cat@example.com")

	response, err := s.Http(s.T()).WithToken(token).Post("/api/v1/logout", nil)
	s.Require().NoError(err)
	response.AssertNoContent()

	response, err = s.Http(s.T()).WithToken(token).Get("/api/v1/me")
	s.Require().NoError(err)
	response.AssertUnauthorized()
}
