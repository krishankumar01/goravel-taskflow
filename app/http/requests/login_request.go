package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (r *LoginRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *LoginRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"email":    "required|string|email",
		"password": "bail|required|string|max_bytes:72",
	}
}

func (r *LoginRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{
		"email": "trim|lower",
	}
}
