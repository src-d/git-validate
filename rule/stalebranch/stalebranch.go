package dco

import (
	"fmt"
	"path"
	"time"

	"github.com/src-d/git-validate/validate"

	"github.com/dustin/go-humanize"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

var defaultConfig = &validate.RuleConfig{
	ID:       "stale-branch",
	Severity: validate.Medium,
	Description: "" +
		"Branch management is an important part of the Git workflow. After some " +
		"time, your list of branches may grow, so it's a good idea to delete " +
		"stale branches. A stale branches is a branches witch last commit was " +
		"done more than 3 months ago",
}

// Kind describes a rule kind that validates the age of the branches.
type Kind struct{}

// Name it honors the validate.RuleKind interface.
func (*Kind) Name() string {
	return "stale-branch"
}

// Rule it honors the validate.RuleKind interface.
func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	cfg.Merge(defaultConfig)
	return &Rule{validate.NewBaseRule(validate.Repository, *cfg)}, nil
}

// Rule of a stalebranch.Kind
type Rule struct {
	validate.BaseRule
}

// Check it honors the validate.Rule interface.
func (r *Rule) Check(repository *git.Repository, _ *object.Commit) ([]*validate.Report, error) {
	iter, err := repository.References()
	if err != nil {
		return nil, err
	}

	head, err := repository.Reference(plumbing.HEAD, false)
	if err != nil {
		return nil, err
	}

	var reports []*validate.Report
	return reports, iter.ForEach(func(ref *plumbing.Reference) error {
		ok, err := r.isValidBranch(head.Target(), ref)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		report, err := r.checkReference(repository, ref)
		if err != nil {
			return err
		}

		if report == nil {
			return nil
		}

		reports = append(reports, report)
		return nil
	})
}

func (r *Rule) isValidBranch(head plumbing.ReferenceName, ref *plumbing.Reference) (bool, error) {
	if !ref.Name().IsRemote() && !ref.Name().IsBranch() {
		return false, nil
	}

	if path.Base(ref.Name().String()) == head.Short() {
		return false, nil
	}

	if ref.Type() == plumbing.SymbolicReference {
		return false, nil
	}

	return true, nil
}

const defaultAge = time.Hour * 24 * 30 * 6

func (r *Rule) checkReference(repository *git.Repository, ref *plumbing.Reference) (*validate.Report, error) {
	c, err := repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	if time.Since(c.Committer.When) > defaultAge {
		return &validate.Report{
			Rule:     r,
			Location: &validate.ReferenceLocation{Reference: ref},
			Severity: validate.Low,
			Message:  fmt.Sprintf("stalled branch, last commit was done %s, consider delete it", humanize.Time(c.Author.When)),
		}, nil
	}

	return nil, nil
}
