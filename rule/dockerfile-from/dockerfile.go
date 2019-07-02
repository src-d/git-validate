package dockerfile

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/src-d/git-validate/validate"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func init() {
	validate.RegisterRuleKind(&Kind{})
}

var defaultConfig = &validate.RuleConfig{
	ID:          "dockerfile",
	Severity:    validate.Medium,
	Description: "",
}

type Kind struct{}

func (*Kind) Name() string {
	return "dockerfile"
}

func (*Kind) Rule(cfg *validate.RuleConfig) (validate.Rule, error) {
	cfg.Merge(defaultConfig)

	r := &Rule{
		BaseRule: validate.NewBaseRule(validate.HEAD, *cfg),
	}

	return r, cfg.LoadParamsTo(&r.Config)
}

type RuleConfig struct {
	Present []string
}

type Rule struct {
	validate.BaseRule
	Config RuleConfig
}

const DockerfilePrefix = "Dockerfile"

func (r *Rule) Check(_ *git.Repository, c *object.Commit) ([]*validate.Report, error) {
	res := &validate.Report{}

	iter, err := c.Files()
	if err != nil {
		return []*validate.Report{res}, err
	}

	err = iter.ForEach(func(f *object.File) error {
		filename := filepath.Base(f.Name)
		if !strings.HasPrefix(filename, DockerfilePrefix) {
			return nil
		}
		r, err := f.Reader()
		if err != nil {
			return err
		}

		r.Close()
		return nil
	})

	return []*validate.Report{res}, nil
}

func parseDockerfile(r io.Reader) error {
	Report, err := parser.Parse(r)
	if err != nil {
		return err
	}

	nodes := []*parser.Node{Report.AST}
	if Report.AST.Children != nil {
		nodes = append(nodes, Report.AST.Children...)
	}

	images := map[string]int{}
	for _, n := range nodes {
		images = nodeSearch("from", n, images)
	}

	fmt.Println(images)
	return nil
}

func nodeSearch(search string, n *parser.Node, a map[string]int) map[string]int {
	if n.Value == search {
		i := strings.Trim(n.Next.Value, "\"")
		if v, ok := a[i]; ok {
			a[i] = v + 1
		} else {
			a[i] = 1

		}
	}
	return a
}
