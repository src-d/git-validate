package validate

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestSeverityColor(t *testing.T) {
	for i := Low - 1; i <= Critical; i++ {
		assert.IsType(t, &color.Color{}, i.Color())
	}
}

func TestSeverityString(t *testing.T) {
	for i := Low - 1; i <= Critical; i++ {
		assert.IsType(t, "", i.String())
	}
}

const severityYAML = `
low: low
high: high
medium: medium
critical: critical
`

func TestSeverityUnmarshalYAML(t *testing.T) {
	out := map[string]Severity{}
	err := yaml.Unmarshal([]byte(severityYAML), out)
	assert.NoError(t, err)

	assert.Len(t, out, 4)
	assert.Equal(t, Low, out["low"])
	assert.Equal(t, Medium, out["medium"])
	assert.Equal(t, High, out["high"])
	assert.Equal(t, Critical, out["critical"])
}

func TestSeverityUnmarshalYAML_Error(t *testing.T) {
	out := map[string]Severity{}
	err := yaml.Unmarshal([]byte("err: no-defined"), out)
	assert.Error(t, err)
}
