package dockerfile

import (
	"testing"
	"time"

	"github.com/src-d/git-compliance/compliance"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/util"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func TestKindRule(t *testing.T) {
	dco, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	assert.Equal(t, dco.ID(), "DOCKERFILE")
}

func TestRuleCheck(t *testing.T) {
	r, c, err := CommitWithFile("Dockerfile", "FROM debian")
	assert.NoError(t, err)

	df, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	result, err := df.Check(r, c)
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	report := result[0]
	assert.NotNil(t, report.Rule)
	assert.False(t, report.Pass)
	assert.Equal(t, report.Code, "DL3006")
	assert.Equal(t, report.Message, "Always tag the version of an image explicitly.")
	assert.Equal(t, report.Location.String(), "Dockerfile:1@31d80c")
}

func TestRuleCheck_Empty(t *testing.T) {
	r, c, err := CommitWithFile("Dockerfile", "")
	assert.NoError(t, err)

	df, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	result, err := df.Check(r, c)
	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestRuleCheck_Nested(t *testing.T) {
	r, c, err := CommitWithFile("foo/Dockerfile", "FROM debian")
	assert.NoError(t, err)

	df, err := (&Kind{}).Rule(&compliance.RuleConfig{})
	assert.NoError(t, err)

	result, err := df.Check(r, c)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.NotNil(t, result[0].Rule)
	assert.False(t, result[0].Pass)
}

func TestRuleCheck_Ignore(t *testing.T) {
	r, c, err := CommitWithFile("foo/Dockerfile", "FROM debian")
	assert.NoError(t, err)

	df, err := (&Kind{}).Rule(&compliance.RuleConfig{
		Params: map[string]interface{}{
			"ignored": []string{"DL3006"},
		},
	})
	assert.NoError(t, err)

	result, err := df.Check(r, c)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].Pass)
}

func CommitWithFile(name, content string) (*git.Repository, *object.Commit, error) {
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, nil, err
	}

	err = util.WriteFile(fs, name, []byte(content), 0644)
	if err != nil {
		return nil, nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, nil, err
	}

	_, err = w.Add(name)
	if err != nil {
		return nil, nil, err
	}

	hash, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
		},
	})

	commit, err := r.CommitObject(hash)
	return r, commit, err
}
