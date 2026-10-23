package factories

import (
	"github.com/brianvoe/gofakeit/v7"

	"taskflow/app/facades"
)

type UserFactory struct {
}

// Definition Define the model's default state.
func (f *UserFactory) Definition() map[string]any {
	password, _ := facades.Hash().Make("password")

	return map[string]any{
		"Name":     gofakeit.Name(),
		"Email":    gofakeit.Email(),
		"Password": password,
		"Role":     "member",
	}
}
