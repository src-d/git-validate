package dco

import (
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/vbatts/git-validation/compliance"
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

func (*Rule) Check(_ *git.Repository, c *object.Commit) (vr compliance.Result, err error) {
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
