package shortsubject

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/src-d/git-compliance/compliance"
)

func init() {
	compliance.RegisterRuleKind(&Kind{})
}

var defaultConfig = &compliance.RuleConfig{
	ID:          "short-subject",
	Severity:    compliance.Medium,
	Description: "commit subject are strictly less than 90 (github ellipsis length)",
}

type Kind struct{}

func (*Kind) Name() string {
	return "short-subject"
}

func (*Kind) Rule(cfg *compliance.RuleConfig) (compliance.Rule, error) {
	cfg.Merge(defaultConfig)
	return &Rule{compliance.NewBaseRule(compliance.History, *cfg)}, nil
}

type Rule struct {
	compliance.BaseRule
}

func (*Rule) Description() string {
	return "commit subject are strictly less than 90 (github ellipsis length)"
}

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*compliance.Report, error) {
	res := &compliance.Report{
		Rule:     r,
		Location: &compliance.CommitLocation{Commit: c},
	}

	if c.NumParents() > 1 {
		return nil, nil
	}

	if c.Message == "" {
		res.Message = "commit subject is empty"
		return []*compliance.Report{res}, nil
	}

	lines := strings.SplitN(c.Message, "\n", 2)
	if len(lines[0]) >= 90 {
		res.Message = "commit subject exceeds 90 characters"
		return []*compliance.Report{res}, nil
	}

	return nil, nil
}
