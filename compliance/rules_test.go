package compliance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type dummyKind struct {
	name string
	ctx  Context
	pass bool
	msg  string
	err  error
}

func (k *dummyKind) Name() string {
	return k.name
}

func (k *dummyKind) Rule(cfg *RuleConfig) (Rule, error) {
	return &dummyRule{
		BaseRule: NewBaseRule(k.ctx, *cfg),
		pass:     k.pass,
		msg:      k.msg,
		err:      k.err,
	}, nil
}

type dummyRule struct {
	BaseRule
	pass bool
	msg  string
	err  error
}

func (r *dummyRule) Check(*git.Repository, *object.Commit) (Result, error) {
	return Result{Pass: r.pass, Message: r.msg}, r.err
}

func TestSeverityString(t *testing.T) {
	assert.Equal(t, Low.String(), "Low")
}

func TestRegisterRuleKind(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo"})

	r, ok := registeredRuleKinds["foo"]
	assert.True(t, ok)
	assert.Equal(t, r.Name(), "foo")

}

func TestRules(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo"})

	cfg := Config{Rules: []RuleConfig{{
		Kind: "foo",
	}}}

	rules, err := Rules(&cfg)
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.NotNil(t, rules[0])
}

func TestRules_NotFound(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo"})

	cfg := Config{Rules: []RuleConfig{{
		Kind: "bar",
	}}}

	rules, err := Rules(&cfg)
	assert.Errorf(t, err, "unable to find")
	assert.Len(t, rules, 0)
}

func TestCommit(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    Context
		isHead bool
		pass   bool
		len    int
	}{
		{"single-head", SingleCommit, true, true, 1},
		{"single-non-head", SingleCommit, false, true, 0},
		{"hisotry-head", History, true, true, 1},
		{"hisotry-non-head", History, true, true, 1},
	}

	for _, tc := range testCases {
		RegisterRuleKind(&dummyKind{
			name: tc.name,
			ctx:  tc.ctx,
			pass: tc.pass,
		})

		cfg := Config{Rules: []RuleConfig{{
			Kind: tc.name,
		}}}

		rules, err := Rules(&cfg)
		assert.NoError(t, err)

		result, err := Commit(rules, nil, nil, tc.isHead)
		assert.NoError(t, err)
		assert.Len(t, result, tc.len)

		if tc.len > 0 {
			assert.True(t, result[0].Pass)
		}
	}
}

func TestResultPassFail(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo", pass: true})
	RegisterRuleKind(&dummyKind{name: "bar", pass: true})
	RegisterRuleKind(&dummyKind{name: "qux", pass: false})

	cfg := Config{Rules: []RuleConfig{
		{Kind: "foo"},
		{Kind: "bar"},
		{Kind: "qux"},
	}}

	rules, err := Rules(&cfg)
	assert.NoError(t, err)
	assert.Len(t, rules, 3)

	result, err := Commit(rules, nil, nil, true)
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	pass, fail := result.PassFail()
	assert.Equal(t, pass, 2)
	assert.Equal(t, fail, 1)
}
