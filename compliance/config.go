package compliance

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Config is a group of rule configurations.
type Config struct {
	// Rules rule configurations.
	Rules []RuleConfig
}

// Decodes a YAML config from a io.Reader
func (c *Config) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
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
