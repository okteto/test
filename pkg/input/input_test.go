// pkg/input/input_test.go
package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInput(t *testing.T) {
	// Arrange: Create sample arguments and environment variables
	args := []string{"my-name", "my-namespace", "my-file.yaml", "true", "false", "VAR1=value1,VAR2=value2", "60s", "test-suite", "info", "pre-commit,release"}
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
	assert.Equal(t, []string{"pre-commit", "release"}, input.Profiles)
	assert.Equal(t, "my-ca-cert", input.CaCert)
}

func TestToParams(t *testing.T) {
	// Arrange: Create sample input
	args := []string{"my-name", "my-namespace", "my-file.yaml", "true", "false", "VAR1=value1,VAR2=value2,VAR3=\"this is a test\"", "60s", "test-suite unit", "info", ""}
	envVars := map[string]string{}
	input, _ := NewInput(args, envVars)

	// Act: Get the command parameters
	params := input.ToParams()

	// Assert: Validate the params
	expectedParams := []string{
		"test-suite",
		"unit",
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

func TestToParams_WithProfiles(t *testing.T) {
	// Arrange: Create sample input without tests, as they cannot be combined with profiles
	args := []string{"", "", "", "false", "false", "", "5m", "", "", "pre-commit, release"}
	input, _ := NewInput(args, map[string]string{})

	// Act: Get the command parameters
	params := input.ToParams()

	// Assert: Each profile is passed as a --profile flag. The empty tests input is passed
	// as an empty argument, which okteto test ignores
	expectedParams := []string{
		"",
		"--profile=pre-commit",
		"--profile=release",
		"--timeout=5m",
	}
	assert.ElementsMatch(t, expectedParams, params)
}

func TestNewInput_WithMissingArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"OneArgument", []string{"my-name"}},
		{"NoProfiles", []string{"my-name", "my-namespace", "my-file.yaml", "true", "false", "", "60s", "", "info"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewInput(tt.args, map[string]string{})
			assert.ErrorIs(t, err, ErrInsufficientArguments)
		})
	}
}

func TestParseProfiles(t *testing.T) {
	tests := []struct {
		name     string
		profiles string
		expected []string
	}{
		{"Empty", "", []string{}},
		{"OneProfile", "pre-commit", []string{"pre-commit"}},
		{"SeveralProfiles", "pre-commit,release", []string{"pre-commit", "release"}},
		{"SpacesAreTrimmed", " pre-commit , release ", []string{"pre-commit", "release"}},
		{"EmptyProfilesAreIgnored", "pre-commit,, ,release,", []string{"pre-commit", "release"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, parseProfiles(tt.profiles))
		})
	}
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
