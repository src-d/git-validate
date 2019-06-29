package compliance

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestReportString(t *testing.T) {
	color.NoColor = true
	r := &Report{
		Rule: &dummyRule{BaseRule: NewBaseRule(History, RuleConfig{
			ID:       "foo",
			Severity: High,
		})},
		Message: "bar",
	}

	assert.Equal(t, "      HIGH [FOO] bar", r.String())
}

func TestReportString_WithSeverityCodeAndLocation(t *testing.T) {
	color.NoColor = true
	r := &Report{
		Rule: &dummyRule{BaseRule: NewBaseRule(History, RuleConfig{
			ID:       "foo",
			Severity: High,
		})},
		Code:     "0001",
		Message:  "bar",
		Severity: Critical,
		Location: &CommitLocation{Commit: &object.Commit{}},
	}

	assert.Equal(t, "  CRITICAL [FOO|0001] bar (000000)", r.String())
}
