package file

import (
	"testing"

	"github.com/src-d/git-compliance/compliance"

	"github.com/stretchr/testify/assert"
	fixtures "gopkg.in/src-d/go-git-fixtures.v3"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func TestKindRule(t *testing.T) {
	dco, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	assert.Equal(t, dco.ID(), "file")
}

func TestRuleCheck_Pass(t *testing.T) {
	err := fixtures.Init()
	assert.NoError(t, err)

	f := fixtures.Basic().One()
	r, err := git.Open(filesystem.NewStorage(f.DotGit(), cache.NewObjectLRUDefault()), nil)
	assert.NoError(t, err)

	c, err := r.CommitObject(f.Head)

	testCases := []struct {
		files []string
		pass  bool
	}{
		{[]string{}, true},
		{[]string{"LICENSE"}, true},
		{[]string{"LICENSE", "not-present"}, false},
		{[]string{"not-present"}, false},
	}

	for _, tc := range testCases {
		dco, err := (&Kind{}).Rule(&compliance.RuleConfig{
			Params: map[string]interface{}{
				"present": tc.files,
			},
		})

		assert.NoError(t, err)

		result, err := dco.Check(r, c)
		assert.NoError(t, err)
		assert.Equal(t, result[0].Pass, tc.pass)
	}
}
