package compliance

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cfgExample = `
rules:
  - id: bar
    severity: high
    params:
      foo: 42
  - severity: low
  - severity: critical
  - severity: medium
`

func TestConfigDecode(t *testing.T) {
	cfg := &Config{}
	err := cfg.Decode(strings.NewReader(cfgExample))

	assert.NoError(t, err)
	assert.Len(t, cfg.Rules, 4)
	assert.Equal(t, cfg.Rules[0].ID, "bar")
	assert.Equal(t, cfg.Rules[0].Severity, High)
	assert.Equal(t, cfg.Rules[0].Params["foo"], 42)
}

const cfgExample_InvalidSeverity = `
rules:
  - severity: foo 
`

func TestConfigDecode_Error(t *testing.T) {
	cfg := &Config{}
	err := cfg.Decode(strings.NewReader(cfgExample_InvalidSeverity))
	assert.Error(t, err)
}

func TestRuleConfigLoadParamsTo(t *testing.T) {
	cfg := &RuleConfig{Params: map[string]interface{}{"foo": "bar"}}

	target := &struct{ Foo string }{}

	err := cfg.LoadParamsTo(target)

	assert.NoError(t, err)
	assert.Equal(t, target.Foo, "bar")
}
