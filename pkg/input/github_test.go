package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGithubEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected []string
	}{
		{
			name: "No GitHub env vars",
			envVars: map[string]string{
				"OTHER_VAR":   "value1",
				"ANOTHER_VAR": "value2",
			},
			expected: []string{},
		},
		{
			name: "Some GitHub env vars",
			envVars: map[string]string{
				"GITHUB_TOKEN": "token123",
				"OTHER_VAR":    "value1",
				"GITHUB_REPO":  "repo123",
			},
			expected: []string{
				"GITHUB_TOKEN=token123",
				"GITHUB_REPO=repo123",
			},
		},
		{
			name: "All GitHub env vars",
			envVars: map[string]string{
				"GITHUB_TOKEN": "token123",
				"GITHUB_REPO":  "repo123",
			},
			expected: []string{
				"GITHUB_TOKEN=token123",
				"GITHUB_REPO=repo123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getGithubEnvVars(tt.envVars)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
