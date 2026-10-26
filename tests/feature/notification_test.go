package feature

import (
	"fmt"
	"testing"

	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/support/http"
	"github.com/stretchr/testify/suite"

	"taskflow/app/enums"
	"taskflow/app/events"
	"taskflow/app/facades"
	"taskflow/app/models"
	"taskflow/tests"
)

type NotificationTestSuite struct {
	suite.Suite
	tests.TestCase
	ben models.User
	cat models.User
}

func TestNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

func (s *NotificationTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.Require().True(facades.Cache().Flush())

	var ana, ben, cat models.User
	s.Require().NoError(facades.Orm().Factory().Create(&ana, map[string]any{"Name": "Ana", "Email": "ana@example.com", "Role": enums.RoleManager}))
	s.Require().NoError(facades.Orm().Factory().Create(&ben, map[string]any{"Name": "Ben", "Email": "ben@example.com"}))
	s.Require().NoError(facades.Orm().Factory().Create(&cat, map[string]any{"Name": "Cat", "Email": "cat@example.com"}))

	s.ben, s.cat = ben, cat
}

func (s *NotificationTestSuite) login(email string) string {
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

func (s *NotificationTestSuite) createTask(token, title string, assignedTo uint) models.Task {
	body, err := http.NewBody().SetFields(map[string]any{
		"title":       title,
		"description": "Assigned in a test.",
		"assigned_to": assignedTo,
	}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithToken(token).WithHeader("Content-Type", body.ContentType()).Post("/api/v1/tasks", body.Reader())
	s.Require().NoError(err)

	var task models.Task
	s.Require().NoError(response.AssertCreated().Bind(&task))

	return task
}

func (s *NotificationTestSuite) notifications(user models.User) []models.Notification {
	var notifications []models.Notification
	s.Require().NoError(facades.Orm().Query().Where("user_id", user.ID).OrderBy("id").Find(&notifications))

	return notifications
}

func (s *NotificationTestSuite) TestAssigning() {
	ana := s.login("ana@example.com")

	task := s.createTask(ana, "Fix the login page", s.ben.ID)

	patch := func(fields map[string]any) {
		body, err := http.NewBody().SetFields(fields).Build()
		s.Require().NoError(err)
		response, err := s.Http(s.T()).WithToken(ana).WithHeader("Content-Type", body.ContentType()).
			Patch(fmt.Sprintf("/api/v1/tasks/%d", task.ID), body.Reader())
		s.Require().NoError(err)
		response.AssertOk()
	}

	s.Run("a new task notifies its assignee", func() {
		notifications := s.notifications(s.ben)
		s.Require().Len(notifications, 1)
		s.Equal(task.ID, notifications[0].TaskID)
		s.Equal(`Ana assigned you "Fix the login page"`, notifications[0].Message)
	})
	s.Run("a change that keeps the assignee notifies nobody", func() {
		patch(map[string]any{"title": "Fix the login form"})
		patch(map[string]any{"assigned_to": s.ben.ID})
		s.Len(s.notifications(s.ben), 1)
	})
	s.Run("a new assignee is notified", func() {
		patch(map[string]any{"assigned_to": s.cat.ID})
		notifications := s.notifications(s.cat)
		s.Require().Len(notifications, 1)
		s.Equal(`Ana assigned you "Fix the login form"`, notifications[0].Message)
		s.Len(s.notifications(s.ben), 1)
	})
	s.Run("a manager who takes the task is not notified", func() {
		var ana models.User
		s.Require().NoError(facades.Orm().Query().Where("email", "ana@example.com").FirstOrFail(&ana))
		patch(map[string]any{"assigned_to": ana.ID})
		s.Empty(s.notifications(ana))
	})
	s.Run("unassigning notifies nobody", func() {
		patch(map[string]any{"assigned_to": nil})
		s.Len(s.notifications(s.cat), 1)
	})
}

func (s *NotificationTestSuite) TestTheEventChecksItsArguments() {
	err := facades.Event().Job(&events.TaskAssigned{}, []event.Arg{{Type: "uint", Value: uint(1)}}).Dispatch()
	s.EqualError(err, "task assigned: needs the task id, the assignee id and the id of who assigned it")

	err = facades.Event().Job(&events.TaskAssigned{}, []event.Arg{
		{Type: "int", Value: 1}, {Type: "uint", Value: uint(2)}, {Type: "uint", Value: uint(3)},
	}).Dispatch()
	s.EqualError(err, "task assigned: every id must be a uint")
}

func (s *NotificationTestSuite) TestEachUserSeesOnlyTheirOwn() {
	ana := s.login("ana@example.com")
	for _, user := range []models.User{s.ben, s.cat} {
		s.createTask(ana, "Write the API docs", user.ID)
	}

	response, err := s.Http(s.T()).WithToken(s.login("ben@example.com")).Get("/api/v1/notifications")
	s.Require().NoError(err)

	var page struct {
		Data []models.Notification `json:"data"`
		Meta struct {
			Total int64 `json:"total"`
		} `json:"meta"`
	}
	s.Require().NoError(response.AssertOk().Bind(&page))
	s.Require().Len(page.Data, 1)
	s.Equal(s.ben.ID, page.Data[0].UserID)
	s.Equal(int64(1), page.Meta.Total)
}

func (s *NotificationTestSuite) TestDeletingATaskRemovesItsNotifications() {
	ana := s.login("ana@example.com")

	task := s.createTask(ana, "Fix the login page", s.ben.ID)
	s.Require().Len(s.notifications(s.ben), 1)

	response, err := s.Http(s.T()).WithToken(ana).Delete(fmt.Sprintf("/api/v1/tasks/%d", task.ID), nil)
	s.Require().NoError(err)
	response.AssertNoContent()
	s.Empty(s.notifications(s.ben))
}

func (s *NotificationTestSuite) TestMarkingAsRead() {
	s.createTask(s.login("ana@example.com"), "Fix the login page", s.ben.ID)

	notification := s.notifications(s.ben)[0]
	s.Nil(notification.ReadAt)
	read := fmt.Sprintf("/api/v1/notifications/%d/read", notification.ID)

	s.Run("another user can't mark it", func() {
		response, err := s.Http(s.T()).WithToken(s.login("cat@example.com")).Post(read, nil)
		s.Require().NoError(err)
		response.AssertNotFound()
	})
	s.Run("its user marks it, and a second time keeps the first time", func() {
		ben := s.login("ben@example.com")
		response, err := s.Http(s.T()).WithToken(ben).Post(read, nil)
		s.Require().NoError(err)
		var marked models.Notification
		s.Require().NoError(response.AssertOk().Bind(&marked))
		first := s.notifications(s.ben)[0].ReadAt
		s.Require().NotNil(first)
		s.Equal(first.String(), marked.ReadAt.String())

		response, err = s.Http(s.T()).WithToken(ben).Post(read, nil)
		s.Require().NoError(err)
		response.AssertOk()
		s.Equal(first.String(), s.notifications(s.ben)[0].ReadAt.String())
	})
}
