package unit

import (
	"context"
	"strings"
	"testing"

	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"taskflow/app/rules"
	"taskflow/tests"
)

type RulesTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestRulesTestSuite(t *testing.T) {
	suite.Run(t, new(RulesTestSuite))
}

func (s *RulesTestSuite) TestMaxBytes() {
	cases := []struct {
		name   string
		value  any
		passes bool
	}{
		{"72 bytes", strings.Repeat("a", 72), true},
		{"73 bytes", strings.Repeat("a", 73), false},
		{"18 emoji are 72 bytes", strings.Repeat("🔒", 18), true},
		{"19 emoji are 76 bytes", strings.Repeat("🔒", 19), false},
		{"not a string", 12345, false},
	}

	rule := &rules.MaxBytes{}
	for _, c := range cases {
		s.Equal(c.passes, rule.Passes(context.Background(), nil, c.value, "72"), c.name)
	}
}

func (s *RulesTestSuite) TestFutureDate() {
	cases := []struct {
		name   string
		value  any
		passes bool
	}{
		{"tomorrow", carbon.Now().AddDay().ToDateTimeString(), true},
		{"yesterday", carbon.Now().SubDay().ToDateTimeString(), false},
		{"not a date", "soon", false},
		{"not a string", 12345, false},
	}

	rule := &rules.FutureDate{}
	for _, c := range cases {
		s.Equal(c.passes, rule.Passes(context.Background(), nil, c.value), c.name)
	}
}
