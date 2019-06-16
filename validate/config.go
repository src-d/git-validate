//go:generate stringer -type=Severity
package validate

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Severity int

const (
	_ Severity = iota
	Low
	Medium
	High
	Critical
)

type Config struct {
	Rules []RuleConfig
}

// Decodes a YAML config from a io.Reader
func (c *Config) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

type RuleConfig struct {
	Name     string
	Kind     string
	Rule     string
	Severity Severity
	Params   map[string]interface{}
}

// LoadParamsTo loads the rule config params into a target.
func (c *RuleConfig) LoadParamsTo(target interface{}) error {
	d, err := yaml.Marshal(c.Params)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(d, target)
}
