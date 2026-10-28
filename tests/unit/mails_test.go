package unit

import (
	"testing"

	"github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"taskflow/app/facades"
	"taskflow/app/mails"
	"taskflow/app/models"
	"taskflow/tests"
)

type MailsTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestMailsTestSuite(t *testing.T) {
	suite.Run(t, new(MailsTestSuite))
}

// render builds a mail's HTML and text the way sending does, without sending anything
func (s *MailsTestSuite) render(mailable mail.Mailable) (string, string) {
	engine := mails.NewPongo(facades.App().ResourcePath("mail"))
	content := mailable.Content()
	html, err := engine.Render(content.HtmlView, content.With)
	s.Require().NoError(err)
	text, err := engine.Render(content.TextView, content.With)
	s.Require().NoError(err)

	return html, text
}

func (s *MailsTestSuite) TestAssignment() {
	ana := models.User{Name: "Ana"}
	ben := models.User{Name: "Ben", Email: "ben@example.com"}
	task := models.Task{Title: "Fix the <login> page", Description: "The button & the logo overlap."}

	s.Run("the envelope goes to the assignee", func() {
		envelope := mails.NewAssignment(task, ben, ana).Envelope()
		s.Equal([]string{"ben@example.com"}, envelope.To)
		s.Equal("Ana assigned you a task", envelope.Subject)
	})
	s.Run("the HTML part escapes the title and sits in the layout", func() {
		html, _ := s.render(mails.NewAssignment(task, ben, ana))
		s.Contains(html, "<strong>Fix the &lt;login&gt; page</strong>")
		s.Contains(html, "You get this email because you have a TaskFlow account.")
	})
	s.Run("the text part keeps the title as typed", func() {
		_, text := s.render(mails.NewAssignment(task, ben, ana))
		s.Contains(text, `Ana assigned you "Fix the <login> page".`)
		s.Contains(text, "The button & the logo overlap.")
	})
	s.Run("a due date adds a line, and no due date leaves it out", func() {
		_, text := s.render(mails.NewAssignment(task, ben, ana))
		s.NotContains(text, "due")

		due := carbon.NewDateTime(carbon.Parse("2026-11-06 09:00:00"))
		task.DueAt = due
		_, text = s.render(mails.NewAssignment(task, ben, ana))
		s.Contains(text, "It's due 2026-11-06 09:00:00.")
	})
}
