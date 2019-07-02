package dco

import (
	"testing"

	"github.com/src-d/git-validate/validate"

	"github.com/stretchr/testify/assert"
	fixtures "gopkg.in/src-d/go-git-fixtures.v3"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func TestKindRule(t *testing.T) {
	dco, err := (&Kind{}).Rule(&validate.RuleConfig{})
	assert.NoError(t, err)

	assert.Equal(t, dco.ID(), "STALE-BRANCH")
}

func TestRuleCheck(t *testing.T) {
	err := fixtures.Init()
	assert.NoError(t, err)

	f := fixtures.Basic().One()
	r, err := git.Open(filesystem.NewStorage(f.DotGit(), cache.NewObjectLRUDefault()), nil)
	assert.NoError(t, err)

	df, err := (&Kind{}).Rule(&validate.RuleConfig{})
	assert.NoError(t, err)

	result, err := df.Check(r, nil)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.NotNil(t, result[0].Rule)
	assert.Equal(t, result[0].Location.String(), "refs/heads/branch")
}
