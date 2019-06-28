package dockerfile

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/src-d/git-compliance/compliance"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/zabio3/godolint/linter/rules"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	compliance.RegisterRuleKind(&Kind{})
}

var defaultConfig = &compliance.RuleConfig{
	ID:          "dockerfile",
	Severity:    compliance.Medium,
	Description: "",
}

type Kind struct{}

func (*Kind) Name() string {
	return "dockerfile"
}

func (*Kind) Rule(cfg *compliance.RuleConfig) (compliance.Rule, error) {
	cfg.Merge(defaultConfig)

	r := &Rule{
		BaseRule: compliance.NewBaseRule(compliance.SingleCommit, *cfg),
	}

	return r, cfg.LoadParamsTo(&r.Config)
}

type RuleConfig struct {
	Present []string
}

type Rule struct {
	compliance.BaseRule
	Config RuleConfig
}

const DockerfilePrefix = "Dockerfile"

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]compliance.Result, error) {
	res := compliance.Result{}
	res.Commit = c
	res.Severity = compliance.Medium

	iter, err := c.Files()
	if err != nil {
		return []compliance.Result{res}, err
	}

	var results []compliance.Result
	err = iter.ForEach(func(f *object.File) error {
		filename := filepath.Base(f.Name)
		if !strings.HasPrefix(filename, DockerfilePrefix) {
			return nil
		}
		df, err := f.Reader()
		if err != nil {
			return err
		}

		defer df.Close()

		result, err := r.validateDockerfile(filename, df)
		if err != nil {
			return err
		}

		results = append(results, result...)
		return nil
	})

	return results, nil
}

func (r *Rule) validateDockerfile(filename string, df io.Reader) ([]compliance.Result, error) {
	ast, err := parser.Parse(df)
	if err != nil {
		return nil, err
	}

	var results []compliance.Result
	for _, rule := range rules.Rules {
		result, _ := rule.ValidateFunc(ast.AST)
		results = append(results, r.toComplianceResult(filename, rule, result)...)
	}

	return results, err
}
func (r *Rule) toComplianceResult(filename string, rule *rules.Rule, results []rules.ValidateResult) []compliance.Result {
	if len(results) == 0 {
		return []compliance.Result{{
			Rule:    r,
			Code:    rule.Code,
			Pass:    true,
			Message: rule.Description,
		}}
	}

	msgs := rules.CreateMessage(rule, results)
	list := make([]compliance.Result, len(results))
	for i, msg := range msgs {
		parts := strings.SplitN(msg, " ", 3)

		list[i] = compliance.Result{
			Rule:     r,
			Code:     rule.Code,
			Pass:     false,
			Location: fmt.Sprintf("%s:%s", filename, parts[0]),
			Message:  strings.Trim(parts[2], "\n"),
		}
	}

	return list
}
