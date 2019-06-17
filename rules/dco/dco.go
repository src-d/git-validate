package dco

import (
	"regexp"
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
	return "dco"
}

func (*Kind) Context() validate.Context {
	return validate.History
}

func (*Kind) Rule(*validate.RuleConfig) (validate.Rule, error) {
	return &Rule{}, nil
}

type Rule struct{}

func (*Rule) ID() string {
	return "dco"
}

func (*Rule) Context() validate.Context {
	return validate.History
}

func (*Rule) Description() string {
	return "makes sure the commits are signed"
}

var ValidDCO = regexp.MustCompile(`^Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>$`)

func (*Rule) Check(_ *git.Repository, c *object.Commit) (vr validate.Result, err error) {
	vr.Commit = c
	if c.NumParents() > 1 {
		vr.Pass = true
		vr.Message = "merge commits do not require DCO"
		return vr, nil
	}

	hasValid := false
	for _, line := range strings.Split(c.Message, "\n") {
		if ValidDCO.MatchString(line) {
			hasValid = true
		}
	}

	if !hasValid {
		vr.Pass = false
		vr.Message = "does not have a valid DCO"
	} else {
		vr.Pass = true
		vr.Message = "has a valid DCO"
	}

	return vr, nil
}
