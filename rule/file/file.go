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
	Short:       "Mandatory files %s are present",
	Description: "Verify the presence of certains files in the HEAD commit",
	Params: map[string]interface{}{
		"present": []string{"README.md", "LICENSE"},
	},
}

// Kind describes a rule kind verifies the presensence of certain files.
type Kind struct{}

// Name it honors the compliance.RuleKind interface.
func (*Kind) Name() string {
	return "file"
}

// Rule it honors the compliance.RuleKind interface.
func (*Kind) Rule(cfg *compliance.RuleConfig) (compliance.Rule, error) {
	cfg.Merge(defaultConfig)

	r := &Rule{
		BaseRule: compliance.NewBaseRule(compliance.HEAD, *cfg),
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
	compliance.BaseRule
	Config RuleConfig
}

// ShortDescription it honors compliance.Rule interface.
func (r *Rule) ShortDescription() string {
	return fmt.Sprintf(r.BaseRule.ShortDescription(), r.Config.Present)
}

// Description it honors compliance.Rule interface.
func (r *Rule) Description() string {
	return fmt.Sprintf(r.BaseRule.Description(), r.Config.Present)
}

// Check it honors the compliance.Rule interface.
func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*compliance.Report, error) {
	var reports []*compliance.Report

	for _, present := range r.Config.Present {
		_, err := c.File(present)
		if err != nil && err != object.ErrFileNotFound {
			continue
		}

		report := &compliance.Report{
			Rule:     r,
			Location: &compliance.CommitLocation{Commit: c},
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
