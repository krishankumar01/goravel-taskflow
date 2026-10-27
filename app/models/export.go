package models

import (
	"github.com/goravel/framework/database/orm"

	"taskflow/app/enums"
)

type Export struct {
	orm.Model
	UserID uint               `gorm:"not null" json:"user_id"`
	Status enums.ExportStatus `gorm:"type:enum('pending','done');default:'pending';not null" json:"status"`
	Path   string             `gorm:"size:255" json:"-"`
}
