package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type Notification struct {
	orm.Model
	UserID  uint             `gorm:"not null" json:"user_id"`
	TaskID  uint             `gorm:"not null" json:"task_id"`
	Message string           `gorm:"size:255;not null" json:"message"`
	ReadAt  *carbon.DateTime `json:"read_at"`
}
