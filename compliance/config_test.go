package compliance

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Default(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo"})

	cfg := Config{}

	rules, err := cfg.Rules()
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.NotNil(t, rules[0])
}

func TestConfigRules(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo"})

	cfg := Config{RuleConfigs: []RuleConfig{{
		Kind: "foo",
	}}}

	rules, err := cfg.Rules()
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.NotNil(t, rules[0])
}

func TestConfigRules_NotFound(t *testing.T) {
	RegisterRuleKind(&dummyKind{name: "foo"})

	cfg := Config{RuleConfigs: []RuleConfig{{
		Kind: "bar",
	}}}

	rules, err := cfg.Rules()
	assert.Errorf(t, err, "unable to find")
	assert.Len(t, rules, 0)
}

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
	assert.Len(t, cfg.RuleConfigs, 4)
	assert.Equal(t, cfg.RuleConfigs[0].ID, "bar")
	assert.Equal(t, cfg.RuleConfigs[0].Severity, High)
	assert.Equal(t, cfg.RuleConfigs[0].Params["foo"], 42)
}

const cfgExampleInvalidSeverity = `
rules:
  - severity: foo 
`

func TestConfigDecode_Error(t *testing.T) {
	cfg := &Config{}
	err := cfg.Decode(strings.NewReader(cfgExampleInvalidSeverity))
	assert.Error(t, err)
}

func TestRuleConfigLoadParamsTo(t *testing.T) {
	cfg := &RuleConfig{Params: map[string]interface{}{"foo": "bar"}}

	target := &struct{ Foo string }{}

	err := cfg.LoadParamsTo(target)

	assert.NoError(t, err)
	assert.Equal(t, target.Foo, "bar")
}

func TestRuleConfigMerge(t *testing.T) {
	cfg := &RuleConfig{}
	cfg.Merge(&RuleConfig{
		ID:          "foo",
		Severity:    High,
		Description: "qux",
		Params:      map[string]interface{}{"foo": "bar"},
	})

	assert.Equal(t, cfg.ID, "foo")
	assert.Equal(t, cfg.Severity, High)
	assert.Equal(t, cfg.Description, "qux")
	assert.Len(t, cfg.Params, 1)
}
