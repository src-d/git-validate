package dco

import (
	"testing"

	"github.com/src-d/git-compliance/compliance"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestRuleCheck_Fail(t *testing.T) {
	dco, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	result, err := dco.Check(nil, &object.Commit{Message: "foo"})
	assert.NoError(t, err)
	assert.False(t, result[0].Pass)
}

func TestRuleCheck_Ignore(t *testing.T) {
	dco, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	result, err := dco.Check(nil, &object.Commit{ParentHashes: []plumbing.Hash{
		plumbing.ZeroHash, plumbing.ZeroHash,
	}})

	assert.NoError(t, err)
	assert.True(t, result[0].Pass)
}

func TestRuleCheck_Pass(t *testing.T) {
	dco, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	result, err := dco.Check(nil, &object.Commit{Message: "Signed-off-by: Máximo Cuadros <mcuadros@gmail.com>"})
	assert.NoError(t, err)
	assert.True(t, result[0].Pass)
}

func TestKindRule(t *testing.T) {
	dco, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	assert.Equal(t, dco.ID(), "dco")
}
