package validate

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
		fail bool
	}{
		{"head", HEAD, 1, true},
		{"hisotry", History, 1, false},
		{"hisotry-fail", History, 8, true},
		{"repository", Repository, 1, false},
	}

	for _, tc := range testCases {
		RegisterRuleKind(&dummyKind{
			name: tc.name,
			msg:  "foo",
			lvl:  tc.lvl,
			fail: tc.fail,
		})

		cfg := Config{RuleConfigs: []RuleConfig{{
			Kind:  tc.name,
			Short: tc.name,
		}}}

		r, err := NewRunner(&cfg)
		assert.NoError(t, err)

		result, err := r.Run(basic)
		assert.NoError(t, err)
		assert.Len(t, result, tc.len, tc.name)
	}
}
