package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Task struct {
	orm.Model
	CreatedBy   uint  `gorm:"not null"`
	Creator     *User `gorm:"foreignKey:CreatedBy"`
	AssignedTo  *uint
	Assignee    *User  `gorm:"foreignKey:AssignedTo"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text;not null"`
	Status      Status `gorm:"type:enum('pending','in-progress','done');default:'pending';not null"`
	DueAt       *carbon.DateTime
}
