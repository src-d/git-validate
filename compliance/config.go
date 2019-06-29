package compliance

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Config is a group of rule configurations.
type Config struct {
	// RuleConfigs rule configurations.
	RuleConfigs []RuleConfig `yaml:"rules"`
}

// Decode a YAML config from a io.Reader
func (c *Config) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

// Rules generates the rules based on a given config.
func (c *Config) Rules() ([]Rule, error) {
	rules := make([]Rule, len(c.RuleConfigs))
	for i, rc := range c.RuleConfigs {
		var err error
		rules[i], err = c.rule(&rc)
		if err != nil {
			return nil, err
		}
	}

	return rules, nil
}

func (c *Config) rule(cfg *RuleConfig) (Rule, error) {
	k, ok := registeredRuleKinds[cfg.Kind]
	if !ok {
		return nil, fmt.Errorf("unable to find %q kind", cfg.Kind)
	}

	return k.Rule(cfg)
}

// RuleConfig contains the configuration for a rule.
type RuleConfig struct {
	// Kind kind of the rule, from the list of supported rule kinds.
	Kind string
	// ID short self-explenatory id of the rule.
	ID string
	// Description longer description for readability.
	Description string
	// Severity of the rule.
	Severity Severity
	// Params is a map of params to pass as configuration to the kind.
	Params map[string]interface{}
}

// LoadParamsTo loads the rule config params into a target.
func (c *RuleConfig) LoadParamsTo(target interface{}) error {
	d, err := yaml.Marshal(c.Params)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(d, target)
}

// Merge merges the given rule config with the receiver.
func (c *RuleConfig) Merge(cfg *RuleConfig) {
	if c.ID == "" {
		c.ID = cfg.ID
	}

	if c.Description == "" {
		c.Description = cfg.Description
	}

	if c.Severity == 0 {
		c.Severity = cfg.Severity
	}

	if len(c.Params) == 0 {
		c.Params = cfg.Params
	}
}
