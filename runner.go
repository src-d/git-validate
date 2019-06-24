package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vbatts/git-validation/compliance"

	"github.com/fatih/color"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Runner is the for processing a set of rules against a range of commits
type Runner struct {
	Repository *git.Repository
	Config     compliance.Config
	Verbose    bool
}

// NewRunner returns an initiallized Runner.
func NewRunner(root string, config string, verbose bool) (*Runner, error) {
	runner := &Runner{}

	var err error
	runner.Repository, err = git.PlainOpen(root)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(config)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	if err := runner.Config.Decode(file); err != nil {
		return nil, err
	}

	fmt.Println(err, runner.Config)
	return runner, nil
}

func shortCommitMessage(c *object.Commit) string {
	lines := strings.SplitN(c.Message, "\n", 2)
	return fmt.Sprintf("%.80s", lines[0])
}

// Run processes the rules for each commit in the range provided
func (r *Runner) Run() (compliance.Results, error) {
	rules, err := compliance.Rules(&r.Config)
	if err != nil {
		return nil, err
	}

	iter, err := r.Repository.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	isHead := true
	results := make(compliance.Results, 0)
	return results, iter.ForEach(func(c *object.Commit) error {
		vr, err := compliance.Commit(rules, r.Repository, c, isHead)
		if err != nil {
			return err
		}

		isHead = false
		results = append(results, vr...)

		_, fail := vr.PassFail()
		if os.Getenv("QUIET") != "" {
			if fail != 0 {
				for _, res := range vr {
					if !res.Pass {
						fmt.Printf(" %s - FAIL - %s\n", c.Hash.String(), res.Message)
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

				fmt.Printf("   └ %s %s [%s]  %s\n", res.Rule.Severity(), result, res.Rule.ID(), res.Message)
			} else if !res.Pass {
				fmt.Printf("   └ %s %s [%s] %s\n", res.Rule.Severity(), color.RedString("FAIL"), res.Rule.ID(), res.Message)
			}
		}

		return nil
	})
}
