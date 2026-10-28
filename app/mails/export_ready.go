package mails

import (
	"github.com/goravel/framework/contracts/mail"

	"taskflow/app/models"
)

type ExportReady struct {
	export models.Export
	user   models.User
}

func NewExportReady(export models.Export, user models.User) *ExportReady {
	return &ExportReady{export: export, user: user}
}

func (receiver *ExportReady) Headers() map[string]string {
	return map[string]string{}
}

func (receiver *ExportReady) Attachments() []string {
	return []string{}
}

func (receiver *ExportReady) Content() *mail.Content {
	return &mail.Content{
		HtmlView: "export_ready.html",
		TextView: "export_ready.txt",
		With: map[string]any{
			"Name":   receiver.user.Name,
			"Export": receiver.export,
		},
	}
}

func (receiver *ExportReady) Envelope() *mail.Envelope {
	return &mail.Envelope{
		To:      []string{receiver.user.Email},
		Subject: "Your task export is ready",
	}
}

func (receiver *ExportReady) Queue() *mail.Queue {
	return &mail.Queue{}
}
