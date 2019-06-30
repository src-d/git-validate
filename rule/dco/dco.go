package dco

import (
	"regexp"
	"strings"

	"github.com/src-d/git-compliance/compliance"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	compliance.RegisterRuleKind(&Kind{})
}

var defaultConfig = &compliance.RuleConfig{
	ID:       "dco",
	Severity: compliance.Medium,
	Description: "" +
		"Enforces the [Developer Certificate of Origin](https://developercertificate.org/) " +
		"(DCO) on commits. It requires all commit messages to contain the Signed-off-by " +
		"line with an email address that matches the commit author.",
}

// Kind describe a rule kind that validates all the commits in a repository are
// signed-off.
type Kind struct{}

// Name it honors the compliance.RuleKind interface.
func (*Kind) Name() string {
	return "dco"
}

// Rule it honors the compliance.RuleKind interface.
func (*Kind) Rule(cfg *compliance.RuleConfig) (compliance.Rule, error) {
	cfg.Merge(defaultConfig)
	return &Rule{compliance.NewBaseRule(compliance.History, *cfg)}, nil
}

// Rule of a dco.Kind
type Rule struct {
	compliance.BaseRule
}

// ValidDCO regexp used to validate the commit message.
var ValidDCO = regexp.MustCompile(`^Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>$`)

// Check it honors the compliance.Rule interface.
func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*compliance.Report, error) {
	if c.NumParents() > 1 {
		return nil, nil
	}

	hasValid := false
	for _, line := range strings.Split(c.Message, "\n") {
		if ValidDCO.MatchString(line) {
			hasValid = true
		}
	}

	if hasValid {
		return nil, nil
	}

	return []*compliance.Report{{
		Rule:     r,
		Message:  "does not have a valid DCO",
		Location: &compliance.CommitLocation{Commit: c},
	}}, nil
}
