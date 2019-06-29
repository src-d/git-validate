package compliance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type dummyKind struct {
	name string
	lvl  Level
	fail bool
	msg  string
	err  error
}

func (k *dummyKind) Name() string {
	return k.name
}

func (k *dummyKind) Rule(cfg *RuleConfig) (Rule, error) {
	return &dummyRule{
		BaseRule: NewBaseRule(k.lvl, *cfg),
		fail:     k.fail,
		msg:      k.msg,
		err:      k.err,
	}, nil
}

type dummyRule struct {
	BaseRule
	fail bool
	msg  string
	err  error
}

func (r *dummyRule) Check(*git.Repository, *object.Commit) ([]*Report, error) {
	if !r.fail {
		return nil, nil
	}

	return []*Report{{Message: r.msg}}, r.err
}

func TestRegisterRuleKind(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo"})

	r, ok := registeredRuleKinds["foo"]
	assert.True(t, ok)
	assert.Equal(t, r.Name(), "foo")

}
