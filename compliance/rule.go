package compliance

import (
	"strings"
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

// Rule is a rule validator for a given git.Repository and a object.Commit.
type Rule interface {
	// ID short self-explenatory id of the rule.
	ID() string
	// Level represent at what level the rule is checked.
	Level() Level
	// Severity returns the severity of the rule.
	Severity() Severity
	// Description longer description for readability.
	Description() string
	// Check evaluate a repository and a commit againts this rule.
	Check(*git.Repository, *object.Commit) ([]*Result, error)
}

// BaseRule used to avoid code duplication on the creation of new rules and kinds.
type BaseRule struct {
	level  Level
	config RuleConfig
}

// NewBaseRule returns a new base rule for the given context and config.
func NewBaseRule(l Level, cfg RuleConfig) BaseRule {
	return BaseRule{level: l, config: cfg}
}

// ID honors the Rule interface.
func (r *BaseRule) ID() string {
	return strings.ToUpper(r.config.ID)
}

// Level honors the Rule interface.
func (r *BaseRule) Level() Level {
	return r.level
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

// RegisterRuleKind includes the RuleKind in the available set to use.
func RegisterRuleKind(vr RuleKind) {
	registerRuleLock.Lock()
	defer registerRuleLock.Unlock()

	registeredRuleKinds[vr.Name()] = vr
}
