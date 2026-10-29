package mails

import (
	"github.com/goravel/framework/contracts/mail"

	"taskflow/app/models"
)

type DueSoon struct {
	assignee models.User
	tasks    []models.Task
}

func NewDueSoon(assignee models.User, tasks []models.Task) *DueSoon {
	return &DueSoon{assignee: assignee, tasks: tasks}
}

func (receiver *DueSoon) Headers() map[string]string {
	return map[string]string{}
}

func (receiver *DueSoon) Attachments() []string {
	return []string{}
}

func (receiver *DueSoon) Content() *mail.Content {
	return &mail.Content{
		HtmlView: "due_soon.html",
		TextView: "due_soon.txt",
		With: map[string]any{
			"Name":  receiver.assignee.Name,
			"Tasks": receiver.tasks,
		},
	}
}

func (receiver *DueSoon) Envelope() *mail.Envelope {
	return &mail.Envelope{
		To:      []string{receiver.assignee.Email},
		Subject: "Your tasks due soon",
	}
}

func (receiver *DueSoon) Queue() *mail.Queue {
	return &mail.Queue{}
}
