package dco

import (
	"regexp"
	"strings"

	"github.com/src-d/git-validate/validate"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

var defaultConfig = &validate.RuleConfig{
	ID:       "dco",
	Severity: validate.Medium,
	Short:    "All commits are signed-off",
	Description: "" +
		"Enforces the [Developer Certificate of Origin](https://developercertificate.org/) " +
		"(DCO) on commits. It requires all commit messages to contain the Signed-off-by " +
		"line with an email address that matches the commit author.",
}

// Kind describes a rule kind that validates all the commits in a repository are
// signed-off.
type Kind struct{}

// Name it honors the validate.RuleKind interface.
func (*Kind) Name() string {
	return "dco"
}

// Rule it honors the validate.RuleKind interface.
func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	cfg.Merge(defaultConfig)
	return &Rule{validate.NewBaseRule(validate.History, *cfg)}, nil
}

// Rule of a dco.Kind
type Rule struct {
	validate.BaseRule
}

// ValidDCO regexp used to validate the commit message.
var ValidDCO = regexp.MustCompile(`^Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>$`)

// Check it honors the validate.Rule interface.
func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*validate.Report, error) {
	var msg string
	if c.NumParents() > 1 {
		return nil, nil
	}

	hasValid := false
	msg = "Commit does not have a valid DCO"
	for _, line := range strings.Split(c.Message, "\n") {
		if ValidDCO.MatchString(line) {
			msg = "Commit has a valid DCO"
			hasValid = true
		}
	}

	return []*validate.Report{{
		Rule:     r,
		Pass:     hasValid,
		Message:  msg,
		Location: &validate.CommitLocation{Commit: c},
	}}, nil
}
