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

	if facades.Config().GetString("filesystems.default") != "testing" {
		panic("the tests delete uploaded files, so they only run on the testing disk: go test ./tests/... --env=.env.testing")
	}

	if facades.Config().GetString("queue.default") != "sync" {
		panic("the tests check what queued listeners wrote, so they only run on the sync queue: go test ./tests/... --env=.env.testing")
	}
}

type TestCase struct {
	testing.TestCase
}
