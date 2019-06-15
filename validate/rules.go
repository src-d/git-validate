package validate

import (
	"fmt"
	"sync"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type RuleKind interface {
	// Name short self-explenatory name of the kind.
	Name() string
	// Rule returns a new rule based on the given config, or an error if the
	// configuration is wrong.
	Rule(*RuleConfig) (Rule, error)
}

type Rule interface {
	// ID short self-explenatory name of the rule.
	ID() string
	// Description longer description for readability.
	Description() string
	// Check evaluate a repository and a commit againts this rule.
	Check(*git.Repository, *object.Commit) (Result, error)
}

var (
	RegisteredRuleKinds = make(map[string]RuleKind, 0)
	registerRuleLock    = sync.Mutex{}
)

// RegisterRuleKind includes the RuleKind in the avaible set to use
func RegisterRuleKind(vr RuleKind) {
	registerRuleLock.Lock()
	defer registerRuleLock.Unlock()

	RegisteredRuleKinds[vr.Name()] = vr
}

func Rules(cfg *Config) ([]Rule, error) {
	rules := make([]Rule, len(cfg.Rules))
	for i, rc := range cfg.Rules {
		k, ok := RegisteredRuleKinds[rc.Kind]
		if !ok {
			return nil, fmt.Errorf("unable to find %q kind", rc.Kind)
		}

		var err error
		rules[i], err = k.Rule(&rc)
		if err != nil {
			return nil, err
		}
	}

	return rules, nil
}

// Commit processes the given rules on the provided commit, and returns the
// result set.
func Commit(r *git.Repository, c *object.Commit, rules []Rule) (Results, error) {
	results := Results{}
	for _, rule := range rules {
		result, err := rule.Check(r, c)
		if err != nil {
			return results, err
		}

		result.Rule = rule
		results = append(results, result)
	}

	return results, nil
}

// Result is the result for a single validation of a commit.
type Result struct {
	Rule   Rule
	Commit *object.Commit
	Pass   bool
	Msg    string
}

// Results is a set of results. This is type makes it easy for the following function.
type Results []Result

// PassFail gives a quick over/under of passes and failures of the results in this set
func (vr Results) PassFail() (pass int, fail int) {
	for _, res := range vr {
		if res.Pass {
			pass++
		} else {
			fail++
		}
	}

	return pass, fail
}
