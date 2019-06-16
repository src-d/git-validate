package file

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/vbatts/git-validation/validate"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

type Kind struct{}

func (*Kind) Name() string {
	return "file"
}

func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	r := &Rule{
		Name: cfg.Name,
	}

	return r, cfg.LoadParamsTo(&r.Config)
}

type RuleConfig struct {
	Present []string
}

type Rule struct {
	Name   string
	Config RuleConfig
}

func (r *Rule) ID() string {
	return r.Name
}

func (r *Rule) Description() string {
	return fmt.Sprintf("file %q shoud be present", r.Config.Present)
}

func (r *Rule) Check(_ *git.Repository, c *object.Commit) (vr validate.Result, err error) {
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
		vr.Msg = fmt.Sprintf("does not have mandatory files %q", r.Config.Present)
	} else {
		vr.Pass = true
		vr.Msg = "has all the mandatory files"
	}

	return vr, nil
}
