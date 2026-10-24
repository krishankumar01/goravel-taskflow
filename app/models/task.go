package models

import (
	"github.com/goravel/framework/contracts/database/factory"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"

	"taskflow/app/enums"
	"taskflow/database/factories"
)

type Task struct {
	orm.Model
	CreatedBy   uint             `gorm:"not null" json:"created_by"`
	Creator     *User            `gorm:"foreignKey:CreatedBy" json:"creator"`
	AssignedTo  *uint            `json:"assigned_to"`
	Assignee    *User            `gorm:"foreignKey:AssignedTo" json:"assignee"`
	Title       string           `gorm:"size:255;not null" json:"title"`
	Description string           `gorm:"type:text;not null" json:"description"`
	Status      enums.Status     `gorm:"type:enum('pending','in-progress','done');default:'pending';not null" json:"status"`
	DueAt       *carbon.DateTime `json:"due_at"`
}

func (t *Task) Factory() factory.Factory {
	return &factories.TaskFactory{}
}
