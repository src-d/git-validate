package compliance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	fixtures "gopkg.in/src-d/go-git-fixtures.v3"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func TestRunnerRun(t *testing.T) {
	err := fixtures.Init()
	assert.NoError(t, err)

	f := fixtures.Basic().One()
	basic, err := git.Open(filesystem.NewStorage(f.DotGit(), cache.NewObjectLRUDefault()), nil)
	assert.NoError(t, err)

	testCases := []struct {
		name string
		lvl  Level
		len  int
	}{
		{"head", HEAD, 1},
		{"hisotry", History, 8},
		{"repository", Repository, 1},
	}

	for _, tc := range testCases {
		RegisterRuleKind(&dummyKind{
			name: tc.name,
			lvl:  tc.lvl,
			fail: true,
		})

		cfg := Config{RuleConfigs: []RuleConfig{{
			Kind: tc.name,
		}}}

		r, err := NewRunner(&cfg)
		assert.NoError(t, err)

		result, err := r.Run(basic)
		assert.NoError(t, err)
		assert.Len(t, result, tc.len)
	}
}
