package models

import (
	"github.com/goravel/framework/database/orm"
)

type Role string

const (
	RoleManager Role = "manager"
	RoleMember  Role = "member"
)

type User struct {
	orm.Model
	Name     string `gorm:"size:255;not null"`
	Email    string `gorm:"size:255;unique;not null"`
	Password string `gorm:"size:255;not null"`
	Role     Role   `gorm:"type:enum('manager','member');default:'member';not null"`
}
