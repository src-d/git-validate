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

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]compliance.Result, error) {
	res := compliance.Result{}
	res.Commit = c

	var found int
	for _, present := range r.Config.Present {
		f, err := c.File(present)
		if err != nil {
			if err == object.ErrFileNotFound {
				continue
			}

			return []compliance.Result{res}, err
		}

		if f.Size > 0 {
			found++
		}
	}

	if found != len(r.Config.Present) {
		res.Pass = false
		res.Message = fmt.Sprintf("does not have mandatory files %q", r.Config.Present)
	} else {
		res.Pass = true
		res.Message = "has all the mandatory files"
	}

	return []compliance.Result{res}, nil
}
