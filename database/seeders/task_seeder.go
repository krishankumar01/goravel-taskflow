package seeders

import (
	"taskflow/app/facades"
	"taskflow/app/models"
)

type TaskSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *TaskSeeder) Signature() string {
	return "TaskSeeder"
}

// Run executes the seeder logic.
func (s *TaskSeeder) Run() error {
	var manager models.User
	err := facades.Orm().Query().Where("role", models.RoleManager).FirstOrFail(&manager)
	if err != nil {
		return err
	}

	var members []models.User
	err = facades.Orm().Query().Where("role", models.RoleMember).Find(&members)
	if err != nil {
		return err
	}

	for _, member := range members {
		var tasks []models.Task
		err = facades.Orm().Factory().Count(2).Create(&tasks, map[string]any{
			"CreatedBy":  manager.ID,
			"AssignedTo": &member.ID,
		})
		if err != nil {
			return err
		}
	}

	var unassigned models.Task
	return facades.Orm().Factory().Create(&unassigned, map[string]any{
		"CreatedBy": manager.ID,
	})
}
