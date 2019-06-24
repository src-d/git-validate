package file

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/vbatts/git-validation/compliance"
)

func init() {
	compliance.RegisterRuleKind(&Kind{})
}

type Kind struct{}

func (*Kind) Name() string {
	return "file"
}

func (*Kind) Rule(cfg *compliance.RuleConfig) (compliance.Rule, error) {
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

func (r *Rule) Check(_ *git.Repository, c *object.Commit) (vr compliance.Result, err error) {
	vr.Commit = c

	var found int
	for _, present := range r.Config.Present {
		f, err := c.File(present)
		if err != nil {
			if err == object.ErrFileNotFound {
				continue
			}

			return vr, err
		}

		if f.Size > 0 {
			found++
		}
	}

	if found != len(r.Config.Present) {
		vr.Pass = false
		vr.Message = fmt.Sprintf("does not have mandatory files %q", r.Config.Present)
	} else {
		vr.Pass = true
		vr.Message = "has all the mandatory files"
	}

	return vr, nil
}
