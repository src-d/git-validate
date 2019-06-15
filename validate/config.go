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
