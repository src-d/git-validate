package dco

import (
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/src-d/git-compliance/compliance"
)

func init() {
	compliance.RegisterRuleKind(&Kind{})
}

var defaultConfig = &compliance.RuleConfig{
	ID:          "dco",
	Severity:    compliance.Medium,
	Description: "makes sure the commits are signed",
}

type Kind struct{}

func (*Kind) Name() string {
	return "dco"
}

func (*Kind) Rule(cfg *compliance.RuleConfig) (compliance.Rule, error) {
	cfg.Merge(defaultConfig)
	return &Rule{compliance.NewBaseRule(compliance.History, *cfg)}, nil
}

type Rule struct {
	compliance.BaseRule
}

var ValidDCO = regexp.MustCompile(`^Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>$`)

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*compliance.Result, error) {
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

	return []*compliance.Result{{
		Rule:     r,
		Message:  "does not have a valid DCO",
		Location: &compliance.CommitLocation{Commit: c},
	}}, nil
}
