package rules

import (
	"context"
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/validation"
)

type MaxBytes struct {
}

// Signature The name of the rule.
func (receiver *MaxBytes) Signature() string {
	return "max_bytes"
}

// Passes Determine if the validation rule passes.
func (receiver *MaxBytes) Passes(ctx context.Context, data validation.Data, val any, options ...any) bool {
	value, ok := val.(string)
	if !ok || len(options) == 0 {
		return false
	}

	limit, err := strconv.Atoi(fmt.Sprint(options[0]))
	if err != nil {
		return false
	}

	return len(value) <= limit
}

// Message Get the validation error message.
func (receiver *MaxBytes) Message(ctx context.Context) string {
	return "The :attribute must not be longer than :option0 bytes."
}
