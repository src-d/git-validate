package compliance

import (
	"fmt"
	"sync"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// RuleKind is creates rules of a kind.
type RuleKind interface {
	// Name short self-explenatory name of the kind.
	Name() string
	// Rule returns a new rule based on the given config, or an error if the
	// configuration is wrong.
	Rule(*RuleConfig) (Rule, error)
}

// Severity describes the severity of a rule.
type Severity int

const (
	_ Severity = iota
	Low
	Medium
	High
	Critical
)

func (s Severity) String() string {
	switch s {
	case Low:
		return "LOW"
	case Medium:
		return "MEDIUM"
	case High:
		return "HIGH"
	case Critical:
		return "CRITITCAL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

//UnmarshalYAML honors the yaml.Unmarshaler interface.
func (s *Severity) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	err := unmarshal(&str)
	if err != nil {
		return err
	}

	switch str {
	case "low":
		*s = Low
		return nil
	case "medium":
		*s = Medium
		return nil
	case "high":
		*s = High
		return nil
	case "critical":
		*s = Critical
		return nil
	default:
		return fmt.Errorf("invalid severity value %q", str)
	}
}

type Context int

const (
	_ Context = iota
	SingleCommit
	History
)

// Rule is a rule validator for a given git.Repository and a object.Commit.
type Rule interface {
	// ID short self-explenatory id of the rule.
	ID() string
	// Context represent the context where this rule is checked.
	Context() Context
	// Severity returns the severity of the rule.
	Severity() Severity
	// Description longer description for readability.
	Description() string
	// Check evaluate a repository and a commit againts this rule.
	Check(*git.Repository, *object.Commit) (Result, error)
}

// BaseRule used to avoid code duplication on the creation of new rules and kinds.
type BaseRule struct {
	context Context
	config  RuleConfig
}

// NewBaseRule returns a new base rule for the given context and config.
func NewBaseRule(ctx Context, cfg RuleConfig) BaseRule {
	return BaseRule{context: ctx, config: cfg}
}

// ID honors the Rule interface.
func (r *BaseRule) ID() string {
	return r.config.ID
}

// Context honors the Rule interface.
func (r *BaseRule) Context() Context {
	return r.context
}

// Severity honors the Rule interface.
func (r *BaseRule) Severity() Severity {
	return r.config.Severity
}

// Description honors the Rule interface.
func (r *BaseRule) Description() string {
	return r.config.Description
}

var (
	registeredRuleKinds = make(map[string]RuleKind, 0)
	registerRuleLock    = sync.Mutex{}
)

// RegisterRuleKind includes the RuleKind in the avaible set to use
func RegisterRuleKind(vr RuleKind) {
	registerRuleLock.Lock()
	defer registerRuleLock.Unlock()

	registeredRuleKinds[vr.Name()] = vr
}

//Rules generates the rules based on a given config.
func Rules(cfg *Config) ([]Rule, error) {
	rules := make([]Rule, len(cfg.Rules))
	for i, rc := range cfg.Rules {
		var err error
		rules[i], err = rule(&rc)
		if err != nil {
			return nil, err
		}
	}

	return rules, nil
}

func rule(cfg *RuleConfig) (Rule, error) {
	k, ok := registeredRuleKinds[cfg.Kind]
	if !ok {
		return nil, fmt.Errorf("unable to find %q kind", cfg.Kind)
	}

	return k.Rule(cfg)
}

// Commit processes the given rules on the provided commit, and returns the
// result set.
func Commit(rules []Rule, r *git.Repository, c *object.Commit, isHead bool) (Results, error) {
	results := Results{}
	for _, rule := range rules {
		if !isHead && rule.Context() != History {
			continue
		}

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
	Rule    Rule
	Pass    bool
	Message string

	Commit *object.Commit
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
