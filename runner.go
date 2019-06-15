package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/vbatts/git-validation/validate"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Runner is the for processing a set of rules against a range of commits
type Runner struct {
	Root        string
	Repository  *git.Repository
	Rules       []validate.Rule
	Results     validate.Results
	Verbose     bool
	CommitRange string // if this is empty, then it will default to FETCH_HEAD, then HEAD
}

// NewRunner returns an initiallized Runner.
func NewRunner(root string, rules []validate.Rule, commitrange string, verbose bool) (*Runner, error) {
	r, err := git.PlainOpen(root)
	if err != nil {
		return nil, err
	}

	return &Runner{
		Root:        root,
		Repository:  r,
		Rules:       rules,
		CommitRange: commitrange,
		Verbose:     verbose,
	}, nil
}

func shortCommitMessage(c *object.Commit) string {
	lines := strings.SplitN(c.Message, "\n", 2)
	return fmt.Sprintf("%.80s", lines[0])
}

// Run processes the rules for each commit in the range provided
func (r *Runner) Run() error {
	iter, err := r.Repository.Log(&git.LogOptions{})
	if err != nil {
		return err
	}

	return iter.ForEach(func(c *object.Commit) error {
		vr, err := validate.Commit(nil, c, r.Rules)
		if err != nil {
			fmt.Println(err)
		}

		r.Results = append(r.Results, vr...)
		_, fail := vr.PassFail()
		if os.Getenv("QUIET") != "" {
			if fail != 0 {
				for _, res := range vr {
					if !res.Pass {
						fmt.Printf(" %s - FAIL - %s\n", c.Hash.String(), res.Msg)
					}
				}
			}

			// everything else in the loop is printing output.
			// If we're quiet, then just continue
			return nil
		}

		result := color.GreenString("PASS")
		if fail != 0 {
			result = color.RedString("FAIL")
		}

		if os.Getenv("QUIET") == "" {
			fmt.Printf(" * %s, %s %q\n", result, c.Hash.String(), shortCommitMessage(c))
		}

		for _, res := range vr {
			if r.Verbose {
				result := color.GreenString("PASS")
				if !res.Pass {
					result = color.RedString("FAIL")
				}

				fmt.Printf("   └ %s [%s]  %s\n", result, res.Rule.Name, res.Msg)
			} else if !res.Pass {
				fmt.Printf("   └ %s [%s] %s\n", color.RedString("FAIL"), res.Rule.Name, res.Msg)
			}
		}

		return nil
	})
}
