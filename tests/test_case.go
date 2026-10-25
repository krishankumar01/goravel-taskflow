package tests

import (
	"strings"

	"github.com/goravel/framework/testing"

	"taskflow/app/facades"
	"taskflow/bootstrap"
)

func init() {
	bootstrap.Boot()

	if !strings.HasSuffix(facades.Config().GetString("database.connections.postgres.database"), "_test") {
		panic("the tests drop every table, so they only run on a database whose name ends in _test: go test ./tests/... --env=.env.testing")
	}
}

type TestCase struct {
	testing.TestCase
}
