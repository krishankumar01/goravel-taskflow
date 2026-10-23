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
	Name     string `gorm:"size:255;not null" json:"name"`
	Email    string `gorm:"size:255;unique;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     Role   `gorm:"type:enum('manager','member');default:'member';not null" json:"role"`
}
