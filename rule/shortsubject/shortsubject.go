package shortsubject

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/src-d/git-validate/validate"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

var defaultConfig = &validate.RuleConfig{
	ID:          "short-subject",
	Severity:    validate.Low,
	Short:       "Commit subject are strictly less than 90 characters",
	Description: "Commit subject are strictly less than 90 characters",
}

type Kind struct{}

func (*Kind) Name() string {
	return "short-subject"
}

func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	cfg.Merge(defaultConfig)
	return &Rule{validate.NewBaseRule(validate.History, *cfg)}, nil
}

type Rule struct {
	validate.BaseRule
}

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*validate.Report, error) {
	res := &validate.Report{
		Rule:     r,
		Location: &validate.CommitLocation{Commit: c},
	}

	if c.NumParents() > 1 {
		return nil, nil
	}

	if c.Message == "" {
		res.Pass = false
		res.Message = "commit subject is empty"
		return []*validate.Report{res}, nil
	}

	lines := strings.SplitN(c.Message, "\n", 2)
	if len(lines[0]) >= 90 {
		res.Pass = false
		res.Message = "commit subject exceeds 90 characters"
		return []*validate.Report{res}, nil
	}

	res.Pass = true
	res.Message = "commit subject is strictly less then 90 characters"
	return []*validate.Report{res}, nil
}
