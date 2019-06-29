package compliance

import (
	"fmt"

	"github.com/fatih/color"
)

// Severity describes the severity of a rule.
type Severity int

const (
	_ Severity = iota
	Low
	Medium
	High
	Critical
)

// Color returns a color.Color representing the severity.
func (s Severity) Color() *color.Color {
	switch s {
	case Low:
		return color.New(color.FgCyan)
	case Medium:
		return color.New(color.FgYellow)
	case High:
		return color.New(color.FgRed)
	case Critical:
		return color.New(color.FgHiRed).Add(color.BgWhite)
	default:
		return color.New(color.FgWhite)
	}
}

func (s Severity) String() string {
	switch s {
	case Low:
		return "LOW"
	case Medium:
		return "MEDIUM"
	case High:
		return "HIGH"
	case Critical:
		return "CRITICAL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

//UnmarshalYAML honors the yaml.Unmarshaler interface.
func (s *Severity) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	err := unmarshal(&str)
	if err != nil {
		return err
	}

	switch str {
	case "low":
		*s = Low
		return nil
	case "medium":
		*s = Medium
		return nil
	case "high":
		*s = High
		return nil
	case "critical":
		*s = Critical
		return nil
	default:
		return fmt.Errorf("invalid severity value %q", str)
	}
}
