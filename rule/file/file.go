package file

import (
	"fmt"

	"github.com/src-d/git-validate/validate"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

var defaultConfig = &validate.RuleConfig{
	ID:          "file",
	Severity:    validate.Medium,
	Short:       "Mandatory files %s are present",
	Description: "Verify the presence of certains files in the HEAD commit",
	Params: map[string]interface{}{
		"present": []string{"README.md", "LICENSE"},
	},
}

// Kind describes a rule kind verifies the presensence of certain files.
type Kind struct{}

// Name it honors the validate.RuleKind interface.
func (*Kind) Name() string {
	return "file"
}

// Rule it honors the validate.RuleKind interface.
func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	cfg.Merge(defaultConfig)

	r := &Rule{
		BaseRule: validate.NewBaseRule(validate.HEAD, *cfg),
	}

	return r, cfg.LoadParamsTo(&r.Config)
}

// RuleConfig specific configuration of the rule.
type RuleConfig struct {
	// Present list of file to search on the commit.
	Present []string
}

// Rule of a file.Kind
type Rule struct {
	validate.BaseRule
	Config RuleConfig
}

// ShortDescription it honors validate.Rule interface.
func (r *Rule) ShortDescription() string {
	return fmt.Sprintf(r.BaseRule.ShortDescription(), r.Config.Present)
}

// Description it honors validate.Rule interface.
func (r *Rule) Description() string {
	return fmt.Sprintf(r.BaseRule.Description(), r.Config.Present)
}

// Check it honors the validate.Rule interface.
func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*validate.Report, error) {
	var reports []*validate.Report

	for _, present := range r.Config.Present {
		_, err := c.File(present)
		if err != nil && err != object.ErrFileNotFound {
			continue
		}

		report := &validate.Report{
			Rule:     r,
			Location: &validate.CommitLocation{Commit: c},
		}

		if err == object.ErrFileNotFound {
			report.Pass = false
			report.Message = fmt.Sprintf("unable to find mandatory file %q", present)
		} else {
			report.Pass = true
			report.Message = fmt.Sprintf("mandatory file %q was found", present)
		}

		reports = append(reports, report)
	}

	return reports, nil
}
