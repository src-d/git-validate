package shortsubject

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/vbatts/git-validation/validate"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

type Kind struct{}

func (*Kind) Name() string {
	return "short-subject"
}

func (*Kind) Rule(*validate.RuleConfig) (validate.Rule, error) {
	return &Rule{}, nil
}

type Rule struct{}

func (*Rule) ID() string {
	return "short-subject"
}

func (*Rule) Context() validate.Context {
	return validate.History
}

func (*Rule) Description() string {
	return "commit subject are strictly less than 90 (github ellipsis length)"
}

func (*Rule) Check(_ *git.Repository, c *object.Commit) (vr validate.Result, err error) {
	if c.NumParents() > 1 {
		vr.Pass = true
		vr.Message = "merge commits do not require length check"
		return vr, nil
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
