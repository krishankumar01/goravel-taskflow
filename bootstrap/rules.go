package bootstrap

import (
	"github.com/goravel/framework/contracts/validation"

	"taskflow/app/rules"
)

func Rules() []validation.Rule {
	return []validation.Rule{
		&rules.FutureDate{},
	}
}
