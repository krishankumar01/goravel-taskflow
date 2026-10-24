package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type UpdateTaskRequest struct {
}

func (r *UpdateTaskRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateTaskRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{
		"title":       "trim",
		"description": "trim",
	}
}

func (r *UpdateTaskRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"title":       "sometimes|required|string|max:255",
		"description": "sometimes|required|string",
		"assigned_to": "bail|sometimes|nullable|integer|exists:users,id|regex:^[0-9]+$",
		"due_at":      "bail|sometimes|nullable|date|future_date",
		"status":      "sometimes|required|in:pending,in-progress,done",
	}
}

func (r *UpdateTaskRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"assigned_to": "assignee",
		"due_at":      "due date",
	}
}
