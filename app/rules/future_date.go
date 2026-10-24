package rules

import (
	"context"

	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/support/carbon"
)

type FutureDate struct {
}

// Signature The name of the rule.
func (receiver *FutureDate) Signature() string {
	return "future_date"
}

// Passes Determine if the validation rule passes.
func (receiver *FutureDate) Passes(ctx context.Context, data validation.Data, val any, options ...any) bool {
	value, ok := val.(string)
	if !ok {
		return false
	}

	return carbon.Parse(value).Gt(carbon.Now())
}

// Message Get the validation error message.
func (receiver *FutureDate) Message(ctx context.Context) string {
	return "The :attribute must be a date in the future."
}
