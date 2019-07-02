package file

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

	assert.Equal(t, dco.ID(), "FILE")
}

func TestRuleCheck(t *testing.T) {
	err := fixtures.Init()
	assert.NoError(t, err)

	f := fixtures.Basic().One()
	r, err := git.Open(filesystem.NewStorage(f.DotGit(), cache.NewObjectLRUDefault()), nil)
	assert.NoError(t, err)

	c, err := r.CommitObject(f.Head)

	testCases := []struct {
		files []string
		pass  []bool
	}{
		{[]string{}, []bool{}},
		{[]string{"LICENSE"}, []bool{true}},
		{[]string{"LICENSE", "not-present"}, []bool{true, false}},
		{[]string{"not-present"}, []bool{false}},
	}

	for _, tc := range testCases {
		dco, err := (&Kind{}).Rule(&validate.RuleConfig{
			Params: map[string]interface{}{
				"present": tc.files,
			},
		})

		assert.NoError(t, err)

		result, err := dco.Check(r, c)
		assert.NoError(t, err)
		assert.Len(t, result, len(tc.pass))

		for i, expected := range tc.pass {
			assert.Equal(t, expected, result[i].Pass)
		}
	}
}
