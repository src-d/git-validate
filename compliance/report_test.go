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
		Pass:    true,
		Message: "bar",
	}

	assert.Equal(t, "PASS      HIGH [FOO] bar", r.String())
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

	assert.Equal(t, "FAIL  CRITICAL [FOO|0001] bar (000000)", r.String())
}
