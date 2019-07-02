package dockerfile

import (
	"strings"
	"testing"

	"github.com/src-d/git-validate/validate"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestRuleCheck_Fail(t *testing.T) {

	parseDockerfile(strings.NewReader("FROM foo"))

	return
	dco, err := (&Kind{}).Rule(&validate.RuleConfig{})
	assert.NoError(t, err)

	result, err := dco.Check(nil, &object.Commit{Message: "foo"})
	assert.NoError(t, err)
	assert.False(t, result[0].Pass)
}
