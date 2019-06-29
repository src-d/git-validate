package file

import (
	"fmt"

	"github.com/src-d/git-compliance/compliance"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	compliance.RegisterRuleKind(&Kind{})
}

var defaultConfig = &compliance.RuleConfig{
	ID:          "file",
	Severity:    compliance.Medium,
	Description: "file(s) %q should be present",
}

type Kind struct{}

func (*Kind) Name() string {
	return "file"
}

func (*Kind) Rule(cfg *compliance.RuleConfig) (compliance.Rule, error) {
	cfg.Merge(defaultConfig)

	r := &Rule{
		BaseRule: compliance.NewBaseRule(compliance.SingleCommit, *cfg),
	}

	return r, cfg.LoadParamsTo(&r.Config)
}

type RuleConfig struct {
	Present []string
}

type Rule struct {
	compliance.BaseRule
	Config RuleConfig
}

func (r *Rule) Description() string {
	return fmt.Sprintf("file %q shoud be present", r.Config.Present)
}

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*compliance.Result, error) {
	var results []*compliance.Result

	for _, present := range r.Config.Present {
		_, err := c.File(present)
		if err == nil {
			continue
		}

		if err == object.ErrFileNotFound {
			results = append(results, &compliance.Result{
				Rule:     r,
				Message:  fmt.Sprintf("unable to find mandatory file %q", present),
				Location: &compliance.CommitLocation{Commit: c},
			})
		}
	}

	return results, nil
}
