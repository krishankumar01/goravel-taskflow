package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type RegisterRequest struct {
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (r *RegisterRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *RegisterRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"name":     "required|string|max:255",
		"email":    "bail|required|string|email|max:255|unique:users,email",
		"password": "bail|required|string|min:8|max_bytes:72",
	}
}

func (r *RegisterRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{
		"name":  "trim",
		"email": "trim|lower",
	}
}
