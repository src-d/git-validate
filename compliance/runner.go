package compliance

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Level describes what level is checked with a given rule.
type Level int

const (
	_ Level = iota
	// HEAD rules are only checked againsts the HEAD commit.
	HEAD
	// History rules are checked againsts all the commits in the master history.
	History
	// Repository rules are checked againts the repository. Very convinient
	// to validate git references.
	Repository
)

// Runner executes a group of rules based onf a given Config.
type Runner struct {
	rulesbyLevel map[Level][]Rule
}

// NewRunner returns a new Runner based on a given Config.
func NewRunner(cfg *Config) (*Runner, error) {
	r := &Runner{
		rulesbyLevel: make(map[Level][]Rule, 0),
	}

	return r, r.loadConfig(cfg)
}

func (r *Runner) loadConfig(cfg *Config) error {
	rules, err := cfg.Rules()
	if err != nil {
		return nil
	}

	for _, rule := range rules {
		l := rule.Level()
		r.rulesbyLevel[l] = append(r.rulesbyLevel[l], rule)
	}

	return nil
}

// Run executes all the rules against the given repository.
func (r *Runner) Run(repository *git.Repository) ([]*Report, error) {
	var results []*Report

	if err := r.runbyLevel(Repository, repository, nil, &results); err != nil {
		return nil, err
	}

	iter, err := repository.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	isHead := true
	return results, iter.ForEach(func(c *object.Commit) error {
		if err := r.runbyLevel(History, repository, c, &results); err != nil {
			return err
		}

		if !isHead {
			return nil
		}

		isHead = false
		return r.runbyLevel(HEAD, repository, c, &results)
	})
}

func (r *Runner) runbyLevel(l Level, repository *git.Repository, commit *object.Commit, results *[]*Report) error {
	for _, rule := range r.rulesbyLevel[l] {
		res, err := rule.Check(repository, commit)
		if err != nil {
			return err
		}

		*results = append(*results, res...)
	}

	return nil
}
