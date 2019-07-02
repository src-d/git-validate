package shortsubject

import (
	"strings"
	"testing"

	"github.com/src-d/git-validate/validate"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestRuleCheck(t *testing.T) {
	short, err := (&Kind{}).Rule(&validate.RuleConfig{})
	assert.NoError(t, err)

	testCases := []struct {
		msg string
		len int
	}{
		{"", 1},
		{"foo", 0},
		{strings.Repeat("0", 90), 1},
	}

	for _, tc := range testCases {
		result, err := short.Check(nil, &object.Commit{Message: tc.msg})
		assert.NoError(t, err)
		assert.Len(t, result, tc.len)
	}
}

func TestKindRule(t *testing.T) {
	dco, err := (&Kind{}).Rule(&validate.RuleConfig{})
	assert.NoError(t, err)

	assert.Equal(t, dco.ID(), "SHORT-SUBJECT")
}
