package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

type StoreTaskRequest struct {
	Title       string           `form:"title" json:"title"`
	Description string           `form:"description" json:"description"`
	AssignedTo  *uint            `form:"assigned_to" json:"assigned_to"`
	DueAt       *carbon.DateTime `form:"due_at" json:"due_at"`
}

func (r *StoreTaskRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *StoreTaskRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"title":       "required|string|max:255",
		"description": "required|string",
		"assigned_to": "bail|nullable|integer|exists:users,id|regex:^[0-9]+$",
		"due_at":      "bail|nullable|date|future_date",
	}
}

func (r *StoreTaskRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"assigned_to": "assignee",
		"due_at":      "due date",
	}
}

func (r *StoreTaskRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{
		"title":       "trim",
		"description": "trim",
	}
}
