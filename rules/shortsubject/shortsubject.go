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

func (*Rule) Check(_ *git.Repository, c *object.Commit) (vr compliance.Result, err error) {
	if c.NumParents() > 1 {
		vr.Pass = true
		vr.Message = "merge commits do not require length check"
		return vr, nil
	}

	if c.Message == "" {
		vr.Pass = false
		vr.Message = "commit subject is empty"
		return
	}

	lines := strings.SplitN(c.Message, "\n", 2)
	if len(lines[0]) >= 90 {
		vr.Pass = false
		vr.Message = "commit subject exceeds 90 characters"
		return
	}

	vr.Pass = true
	if len(lines[0]) > 72 {
		vr.Message = "commit subject is under 90 characters, but is still more than 72 chars"
	} else {
		vr.Message = "commit subject is 72 characters or less! *yay*"
	}

	return
}
