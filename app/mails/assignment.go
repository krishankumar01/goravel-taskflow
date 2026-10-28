package mails

import (
	"fmt"

	"github.com/goravel/framework/contracts/mail"

	"taskflow/app/models"
)

type Assignment struct {
	task     models.Task
	assignee models.User
	by       models.User
}

func NewAssignment(task models.Task, assignee, by models.User) *Assignment {
	return &Assignment{task: task, assignee: assignee, by: by}
}

func (receiver *Assignment) Headers() map[string]string {
	return map[string]string{}
}

func (receiver *Assignment) Attachments() []string {
	return []string{}
}

func (receiver *Assignment) Content() *mail.Content {
	return &mail.Content{
		HtmlView: "assignment.html",
		TextView: "assignment.txt",
		With: map[string]any{
			"Name": receiver.assignee.Name,
			"By":   receiver.by.Name,
			"Task": receiver.task,
		},
	}
}

func (receiver *Assignment) Envelope() *mail.Envelope {
	return &mail.Envelope{
		To:      []string{receiver.assignee.Email},
		Subject: fmt.Sprintf("%s assigned you a task", receiver.by.Name),
	}
}

func (receiver *Assignment) Queue() *mail.Queue {
	return &mail.Queue{}
}
