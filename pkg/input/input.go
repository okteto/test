package input

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidInput is returned when the input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrInsufficientArguments is returned when the input arguments are insufficient
	ErrInsufficientArguments = errors.New("insufficient arguments")
)

type Input struct {
	Name      string
	Namespace string
	File      string
	Deploy    bool
	NoCache   bool
	Variables []string
	Timeout   string
	Tests     string
	LogLevel  string
	CaCert    string
}

// NewInput parses the arguments and environment variables into an Input struct
func NewInput(args []string, envVars map[string]string) (*Input, error) {
	if len(args) < 9 {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, ErrInsufficientArguments)
	}

	input := &Input{
		Name:      args[0],
		Namespace: args[1],
		File:      args[2],
		Deploy:    loadBoolOrDefault(args[3], false),
		NoCache:   loadBoolOrDefault(args[4], false),
		Variables: parseVariables(args[5]),
		Timeout:   args[6],
		Tests:     args[7],
		LogLevel:  args[8],
		CaCert:    envVars["OKTETO_CA_CERT"],
	}
	return input, nil
}

func loadBoolOrDefault(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}
	trueValues := map[string]bool{
		"true": true,
		"t":    true,
		"1":    true,
		"yes":  true,
		"y":    true,
	}
	return trueValues[strings.ToLower(value)]
}

func parseVariables(variables string) []string {
	if variables == "" {
		return []string{}
	}
	return strings.Split(variables, ",")
}

func (i *Input) ToParams() []string {
	var params []string

	for _, testName := range strings.Split(i.Tests, " ") {
		params = append(params, testName)
	}

	if i.Name != "" {
		params = append(params, fmt.Sprintf("--name=%s", i.Name))
	}
	if i.Namespace != "" {
		params = append(params, fmt.Sprintf("--namespace=%s", i.Namespace))
	}
	if i.File != "" {
		params = append(params, fmt.Sprintf("--file=%s", i.File))
	}
	if i.Deploy {
		params = append(params, "--deploy")
	}
	if i.NoCache {
		params = append(params, "--no-cache")
	}
	for _, variable := range i.Variables {
		params = append(params, fmt.Sprintf("--var=%s", variable))
	}
	if i.Timeout != "" {
		params = append(params, fmt.Sprintf("--timeout=%s", i.Timeout))
	}
	if i.LogLevel != "" {
		params = append(params, fmt.Sprintf("--log-level=%s", i.LogLevel))
	}
	return params
}
