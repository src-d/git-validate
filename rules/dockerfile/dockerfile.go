package dockerfile

import (
	"io"
	"path/filepath"
	"strconv"
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

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*compliance.Result, error) {
	iter, err := c.Files()
	if err != nil {
		return nil, err
	}

	var results []*compliance.Result
	return results, iter.ForEach(func(f *object.File) error {
		filename := filepath.Base(f.Name)
		if !strings.HasPrefix(filename, DockerfilePrefix) {
			return nil
		}

		df, err := f.Reader()
		if err != nil {
			return err
		}

		defer df.Close()

		result, err := r.validateDockerfile(c, filename, df)
		if err != nil {
			return err
		}

		results = append(results, result...)
		return nil
	})
}

func (r *Rule) validateDockerfile(c *object.Commit, filename string, df io.Reader) ([]*compliance.Result, error) {
	ast, err := parser.Parse(df)
	if err != nil {
		return nil, err
	}

	var results []*compliance.Result
	for _, rule := range rules.Rules {
		result, _ := rule.ValidateFunc(ast.AST)
		results = append(results, r.toComplianceResult(c, filename, rule, result)...)
	}

	return results, err
}
func (r *Rule) toComplianceResult(c *object.Commit, filename string, rule *rules.Rule, results []rules.ValidateResult) []*compliance.Result {
	if len(results) == 0 {
		return nil
	}

	msgs := rules.CreateMessage(rule, results)
	list := make([]*compliance.Result, len(results))
	for i, msg := range msgs {
		parts := strings.SplitN(msg, " ", 3)
		line, _ := strconv.Atoi(strings.Replace(parts[0], "", "#", -1))

		list[i] = &compliance.Result{
			Rule:     r,
			Code:     rule.Code,
			Location: &compliance.LineLocation{Commit: c, Filename: filename, Line: line},
			Message:  strings.Trim(parts[2], "\n"),
		}
	}

	return list
}
