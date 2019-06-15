package validate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cfgExample = `
rules:
  - name: bar
    rule: foo
    params:
      foo: 42
`

func TestConfigDecode(t *testing.T) {
	cfg := &Config{}
	err := cfg.Decode(strings.NewReader(cfgExample))

	assert.NoError(t, err)
	assert.Len(t, cfg.Rules, 1)
	assert.Equal(t, cfg.Rules[0].Name, "bar")
	assert.Equal(t, cfg.Rules[0].Rule, "foo")
	assert.Equal(t, cfg.Rules[0].Params["foo"], 42)
}
