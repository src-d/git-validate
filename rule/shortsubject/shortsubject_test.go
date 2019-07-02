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
		msg  string
		pass []bool
	}{
		{"", []bool{false}},
		{"foo", []bool{true}},
		{strings.Repeat("0", 90), []bool{false}},
	}

	for _, tc := range testCases {
		result, err := short.Check(nil, &object.Commit{Message: tc.msg})
		assert.NoError(t, err)
		assert.Len(t, result, len(tc.pass))
		for i, pass := range tc.pass {
			assert.Equal(t, pass, result[i].Pass)
		}
	}
}

func TestKindRule(t *testing.T) {
	dco, err := (&Kind{}).Rule(&validate.RuleConfig{})
	assert.NoError(t, err)

	assert.Equal(t, dco.ID(), "SHORT-SUBJECT")
}
