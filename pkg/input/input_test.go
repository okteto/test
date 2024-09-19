// pkg/input/input_test.go
package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInput(t *testing.T) {
	// Arrange: Create sample arguments and environment variables
	args := []string{"my-name", "my-namespace", "my-file.yaml", "true", "false", "VAR1=value1,VAR2=value2", "60s", "test-suite", "info"}
	envVars := map[string]string{
		"OKTETO_CA_CERT": "my-ca-cert",
	}

	// Act: Parse the input
	input, err := NewInput(args, envVars)

	// Assert: Validate the input
	assert.NoError(t, err)
	assert.Equal(t, "my-name", input.Name)
	assert.Equal(t, "my-namespace", input.Namespace)
	assert.Equal(t, "my-file.yaml", input.File)
	assert.True(t, input.Deploy)
	assert.False(t, input.NoCache)
	assert.Equal(t, []string{"VAR1=value1", "VAR2=value2"}, input.Variables)
	assert.Equal(t, "60s", input.Timeout)
	assert.Equal(t, "test-suite", input.Tests)
	assert.Equal(t, "info", input.LogLevel)
	assert.Equal(t, "my-ca-cert", input.CaCert)
}

func TestToParams(t *testing.T) {
	// Arrange: Create sample input
	args := []string{"my-name", "my-namespace", "my-file.yaml", "true", "false", "VAR1=value1,VAR2=value2,VAR3=\"this is a test\"", "60s", "test-suite", "info"}
	envVars := map[string]string{}
	input, _ := NewInput(args, envVars)

	// Act: Get the command parameters
	params := input.ToParams()

	// Assert: Validate the params
	expectedParams := []string{
		"--name=my-name",
		"--namespace=my-namespace",
		"--file=my-file.yaml",
		"--deploy",
		"--var=VAR1=value1",
		"--var=VAR2=value2",
		"--var=VAR3=\"this is a test\"",
		"--timeout=60s",
		"--log-level=info",
	}
	assert.ElementsMatch(t, expectedParams, params)
}

func TestNewInput_WithMissingArguments(t *testing.T) {
	args := []string{"my-name"}
	envVars := map[string]string{}

	_, err := NewInput(args, envVars)
	assert.Error(t, err)
}
func TestLoadBoolOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue bool
		expected     bool
	}{
		{"EmptyValueWithDefaultTrue", "", true, true},
		{"EmptyValueWithDefaultFalse", "", false, false},
		{"TrueValue", "true", false, true},
		{"FalseValue", "false", true, false},
		{"YesValue", "yes", false, true},
		{"NoValue", "no", true, false},
		{"OneValue", "1", false, true},
		{"ZeroValue", "0", true, false},
		{"YValue", "y", false, true},
		{"NValue", "n", true, false},
		{"TValue", "t", false, true},
		{"FValue", "f", true, false},
		{"CaseInsensitiveTrue", "TrUe", false, true},
		{"CaseInsensitiveFalse", "FaLsE", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loadBoolOrDefault(tt.value, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
