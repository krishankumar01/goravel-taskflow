package feature

import (
	"testing"

	"github.com/goravel/framework/support/http"
	"github.com/stretchr/testify/suite"

	"taskflow/app/facades"
	"taskflow/app/models"
	"taskflow/tests"
)

type AuthTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.Require().True(facades.Cache().Flush())
}

func (s *AuthTestSuite) TestRoutesNeedAToken() {
	for _, path := range []string{"/api/v1/tasks", "/api/v1/tasks/1", "/api/v1/users", "/api/v1/me"} {
		response, err := s.Http(s.T()).Get(path)
		s.Require().NoError(err)
		response.AssertUnauthorized()
	}
}

func (s *AuthTestSuite) TestABadTokenIsRefused() {
	response, err := s.Http(s.T()).WithToken("not-a-token").Get("/api/v1/tasks")
	s.Require().NoError(err)
	response.AssertUnauthorized()
}

func (s *AuthTestSuite) TestAWrongPasswordIsRefused() {
	var ana models.User
	s.Require().NoError(facades.Orm().Factory().Create(&ana, map[string]any{"Email": "ana@example.com"}))

	body, err := http.NewBody().SetFields(map[string]any{"email": "ana@example.com", "password": "wrong"}).Build()
	s.Require().NoError(err)
	response, err := s.Http(s.T()).WithHeader("Content-Type", body.ContentType()).Post("/api/v1/login", body.Reader())
	s.Require().NoError(err)
	response.AssertUnauthorized()
}
