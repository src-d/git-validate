package file

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/src-d/git-validate/validate"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

var defaultConfig = &validate.RuleConfig{
	ID:       "large-blob",
	Severity: validate.Medium,
	Short:    "All blobs are under recommended size of 1MB",
	Description: "" +
		"Keeping repositories small ensures clones are quick for the users. " +
		"This rule verifies blob objects to ensure that are under 1MB.",
}

// Kind describes a rule kind verifies the size of the blobs.
type Kind struct{}

// Name it honors the validate.RuleKind interface.
func (*Kind) Name() string {
	return "large-file"
}

// Rule it honors the validate.RuleKind interface.
func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	cfg.Merge(defaultConfig)
	return &Rule{BaseRule: validate.NewBaseRule(validate.Repository, *cfg)}, nil
}

// Rule of a largeblob.Kind
type Rule struct {
	validate.BaseRule
}

const mb int64 = 1000 * 1000

var limits = map[validate.Severity]int64{
	validate.Low:      mb,
	validate.Medium:   5 * mb,
	validate.High:     50 * mb,
	validate.Critical: 100 * mb,
}

// Check it honors the validate.Rule interface.
func (r *Rule) Check(repository *git.Repository, c *object.Commit) ([]*validate.Report, error) {
	var reports []*validate.Report

	iter, err := repository.BlobObjects()
	if err != nil {
		return nil, err
	}

	err = iter.ForEach(func(b *object.Blob) error {
		for severity, sz := range limits {
			if b.Size > sz {
				reports = append(reports, &validate.Report{
					Rule:     r,
					Pass:     false,
					Severity: severity,
					Message:  fmt.Sprintf("Blob excess recommended size of %s", humanize.Bytes(uint64(sz))),
					Location: &validate.BlobLocation{Blob: b},
				})
			}
		}

		return nil
	})

	if len(reports) == 0 {
		reports = append(reports, &validate.Report{
			Rule:    r,
			Pass:    true,
			Message: "All blobs are under recommended size of 1MB",
		})
	}

	return reports, nil
}
