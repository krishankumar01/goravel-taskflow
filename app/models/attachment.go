package models

import (
	"github.com/goravel/framework/database/orm"
)

type Attachment struct {
	orm.Model
	TaskID     uint   `gorm:"not null" json:"task_id"`
	UploadedBy uint   `gorm:"not null" json:"uploaded_by"`
	Name       string `gorm:"size:255;not null" json:"name"`
	Path       string `gorm:"size:255;not null" json:"-"`
	Size       int64  `gorm:"not null" json:"size"`
	MimeType   string `gorm:"size:255;not null" json:"mime_type"`
}
