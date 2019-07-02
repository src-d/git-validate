package dockerfile

import (
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/src-d/git-validate/validate"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/zabio3/godolint/linter/rules"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

var defaultConfig = &validate.RuleConfig{
	ID:       "dockerfile",
	Severity: validate.Medium,
	Short:    "All the Dockerfiles complies the Docker's best practices guide",
	Description: "" +
		"Enforce to follow the Best Practices for writing Dockerfiles, the " +
		"recommended best practices and methods for building efficient Docker images." +
		"\n" +
		"https://docs.docker.com/develop/develop-images/dockerfile_best-practices/",
}

// Kind describes a rule kind that validates the Dockerfiles contained in the
// HEAD of the repository.
type Kind struct{}

// Name it honors the validate.RuleKind interface.
func (*Kind) Name() string {
	return "dockerfile"
}

// Rule it honors the validate.RuleKind interface.
func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	cfg.Merge(defaultConfig)

	r := &Rule{BaseRule: validate.NewBaseRule(validate.HEAD, *cfg)}
	return r, cfg.LoadParamsTo(&r.Config)
}

// RuleConfig is the specific configuration for this Kind.
type RuleConfig struct {
	// Ignored allow to ignore rules from godolint.Rules
	// https://github.com/zabio3/godolint#rules
	Ignored []string
}

// Rule of a dockerfile.Kind
type Rule struct {
	Config RuleConfig
	validate.BaseRule
}

// DockerfilePrefix prefix used to find Dockerfiles
const DockerfilePrefix = "Dockerfile"

// Check it honors the validate.Rule interface.
func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*validate.Report, error) {
	iter, err := c.Files()
	if err != nil {
		return nil, err
	}

	var results []*validate.Report
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

func (r *Rule) validateDockerfile(c *object.Commit, filename string, df io.Reader) ([]*validate.Report, error) {
	ast, err := parser.Parse(df)
	if err != nil {
		if err.Error() == "file with no instructions." {
			return nil, nil
		}

		return nil, err
	}

	ignored := make(map[string]struct{})
	for _, i := range r.Config.Ignored {
		ignored[i] = struct{}{}
	}

	var results []*validate.Report
	for _, rule := range rules.Rules {
		if _, ok := ignored[rule.Code]; ok {
			continue
		}

		result, _ := rule.ValidateFunc(ast.AST)
		results = append(results, r.toComplianceResult(c, filename, rule, result)...)
	}

	if len(results) == 0 {
		return []*validate.Report{{
			Rule:     r,
			Pass:     true,
			Message:  "Dockerfile complies the Docker's best practices guide",
			Location: &validate.FileLocation{Commit: c, Filename: filename},
		}}, nil
	}

	return results, err
}

func (r *Rule) toComplianceResult(c *object.Commit, filename string, rule *rules.Rule, results []rules.ValidateResult) []*validate.Report {
	if len(results) == 0 {
		return nil
	}

	msgs := rules.CreateMessage(rule, results)
	list := make([]*validate.Report, len(results))
	for i, msg := range msgs {
		parts := strings.SplitN(msg, " ", 3)
		line, _ := strconv.Atoi(strings.Replace(parts[0], "#", "", -1))

		list[i] = &validate.Report{
			Rule:     r,
			Pass:     false,
			Code:     rule.Code,
			Location: &validate.LineLocation{Commit: c, Filename: filename, Line: line},
			Message:  strings.Trim(parts[2], " \n"),
		}
	}

	return list
}
