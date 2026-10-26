package feature

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	contractshttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/goravel/framework/support/http"
	"github.com/stretchr/testify/suite"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"
	"taskflow/tests"
)

type AttachmentTestSuite struct {
	suite.Suite
	tests.TestCase
	benTask models.Task
	catTask models.Task
	png     []byte
}

func TestAttachmentTestSuite(t *testing.T) {
	suite.Run(t, new(AttachmentTestSuite))
}

func (s *AttachmentTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.Require().True(facades.Cache().Flush())
	s.Require().NoError(facades.Storage().DeleteDirectory("attachments"))

	var ana, ben, cat models.User
	s.Require().NoError(facades.Orm().Factory().Create(&ana, map[string]any{"Email": "ana@example.com", "Role": enums.RoleManager}))
	s.Require().NoError(facades.Orm().Factory().Create(&ben, map[string]any{"Email": "ben@example.com"}))
	s.Require().NoError(facades.Orm().Factory().Create(&cat, map[string]any{"Email": "cat@example.com"}))

	var benTask, catTask models.Task
	s.Require().NoError(facades.Orm().Factory().Create(&benTask, map[string]any{"CreatedBy": ana.ID, "AssignedTo": &ben.ID}))
	s.Require().NoError(facades.Orm().Factory().Create(&catTask, map[string]any{"CreatedBy": ana.ID, "AssignedTo": &cat.ID}))

	var picture bytes.Buffer
	s.Require().NoError(png.Encode(&picture, image.NewRGBA(image.Rect(0, 0, 4, 4))))

	s.benTask, s.catTask, s.png = benTask, catTask, picture.Bytes()
}

func (s *AttachmentTestSuite) login(email string) string {
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

func (s *AttachmentTestSuite) upload(token string, task models.Task, name string, content []byte) contractshttp.Response {
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	part, err := form.CreateFormFile("file", name)
	s.Require().NoError(err)
	_, err = part.Write(content)
	s.Require().NoError(err)
	s.Require().NoError(form.Close())

	response, err := s.Http(s.T()).WithToken(token).
		WithHeader("Content-Type", form.FormDataContentType()).
		Post(fmt.Sprintf("/api/v1/tasks/%d/attachments", task.ID), &body)
	s.Require().NoError(err)

	return response
}

func (s *AttachmentTestSuite) TestUploadAndDownload() {
	ben := s.login("ben@example.com")

	var attachment models.Attachment
	s.Require().NoError(s.upload(ben, s.benTask, "login error.png", s.png).AssertCreated().Bind(&attachment))
	s.Equal("login error.png", attachment.Name)
	s.Equal(int64(len(s.png)), attachment.Size)
	s.Equal("image/png", attachment.MimeType)

	var saved models.Attachment
	s.Require().NoError(facades.Orm().Query().Find(&saved, attachment.ID))
	s.True(strings.HasPrefix(saved.Path, fmt.Sprintf("attachments/%d/", s.benTask.ID)))
	s.True(facades.Storage().Exists(saved.Path))

	response, err := s.Http(s.T()).WithToken(ben).Get(fmt.Sprintf("/api/v1/tasks/%d/attachments/%d", s.benTask.ID, attachment.ID))
	s.Require().NoError(err)
	response.AssertOk().
		AssertHeader("Content-Type", "image/png").
		AssertHeader("Content-Disposition", `attachment; filename="login error.png"`)
	content, err := response.Content()
	s.Require().NoError(err)
	s.Equal(string(s.png), content)
}

func (s *AttachmentTestSuite) TestRules() {
	ana := s.login("ana@example.com")
	ben := s.login("ben@example.com")
	cat := s.login("cat@example.com")

	var attachment models.Attachment
	s.Require().NoError(s.upload(ben, s.benTask, "login-error.png", s.png).AssertCreated().Bind(&attachment))
	file := fmt.Sprintf("/api/v1/tasks/%d/attachments/%d", s.benTask.ID, attachment.ID)

	s.Run("a member can't attach a file to someone else's task", func() {
		s.upload(ben, s.catTask, "login-error.png", s.png).AssertForbidden()
	})
	s.Run("a manager attaches a file to any task", func() {
		s.upload(ana, s.catTask, "login-error.png", s.png).AssertCreated()
	})
	s.Run("a script named .png is refused", func() {
		s.upload(ben, s.benTask, "photo.png", []byte("#!/bin/sh\necho hi\n")).AssertUnprocessableEntity()
	})
	s.Run("a file over 5 MB is refused", func() {
		big := append([]byte("%PDF-1.4\n"), bytes.Repeat([]byte("0"), 6<<20)...)
		s.upload(ben, s.benTask, "spec.pdf", big).AssertUnprocessableEntity()
	})
	s.Run("a member can't open a file on someone else's task", func() {
		response, err := s.Http(s.T()).WithToken(cat).Get(file)
		s.Require().NoError(err)
		response.AssertForbidden()
	})
	s.Run("a file under another task is not found", func() {
		response, err := s.Http(s.T()).WithToken(ana).Get(fmt.Sprintf("/api/v1/tasks/%d/attachments/%d", s.catTask.ID, attachment.ID))
		s.Require().NoError(err)
		response.AssertNotFound()
	})
	s.Run("a member can't delete a file someone else uploaded", func() {
		response, err := s.Http(s.T()).WithToken(cat).Delete(file, nil)
		s.Require().NoError(err)
		response.AssertForbidden()
	})
	s.Run("the uploader deletes the file", func() {
		response, err := s.Http(s.T()).WithToken(ben).Delete(file, nil)
		s.Require().NoError(err)
		response.AssertNoContent()
	})
}

func (s *AttachmentTestSuite) TestUploadLeavesNoTempFile() {
	pattern := filepath.Join(os.TempDir(), facades.Config().GetString("app.name")+"_*")
	before, err := filepath.Glob(pattern)
	s.Require().NoError(err)

	s.upload(s.login("ben@example.com"), s.benTask, "login-error.png", s.png).AssertCreated()

	after, err := filepath.Glob(pattern)
	s.Require().NoError(err)
	s.Len(after, len(before))
}

func (s *AttachmentTestSuite) TestDeletingATaskRemovesItsFiles() {
	var attachment models.Attachment
	s.Require().NoError(s.upload(s.login("ben@example.com"), s.benTask, "login-error.png", s.png).AssertCreated().Bind(&attachment))

	var saved models.Attachment
	s.Require().NoError(facades.Orm().Query().Find(&saved, attachment.ID))

	response, err := s.Http(s.T()).WithToken(s.login("ana@example.com")).Delete(fmt.Sprintf("/api/v1/tasks/%d", s.benTask.ID), nil)
	s.Require().NoError(err)
	response.AssertNoContent()

	count, err := facades.Orm().Query().Model(&models.Attachment{}).Where("task_id", s.benTask.ID).Count()
	s.Require().NoError(err)
	s.Equal(int64(0), count)
	s.True(facades.Storage().Missing(saved.Path))
}
