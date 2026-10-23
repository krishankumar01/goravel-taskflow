package seeders

import (
	"taskflow/app/facades"
	"taskflow/app/models"
)

type UserSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder logic.
func (s *UserSeeder) Run() error {
	var manager models.User
	err := facades.Orm().Factory().Create(&manager, map[string]any{
		"Name":  "Ana",
		"Email": "ana@example.com",
		"Role":  models.RoleManager,
	})
	if err != nil {
		return err
	}

	var members []models.User
	return facades.Orm().Factory().Count(3).Create(&members)
}
