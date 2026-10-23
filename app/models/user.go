package models

import (
	"github.com/goravel/framework/contracts/database/factory"
	"github.com/goravel/framework/database/orm"

	"taskflow/database/factories"
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

	CreatedTasks  []Task `gorm:"foreignKey:CreatedBy" json:"created_tasks"`
	AssignedTasks []Task `gorm:"foreignKey:AssignedTo" json:"assigned_tasks"`
}

func (u *User) Factory() factory.Factory {
	return &factories.UserFactory{}
}
