package input

import (
	"fmt"
	"strings"
)

const (
	GITHUB_ENVVAR_PREFIX = "GITHUB_"
)

func getGithubEnvVars(envVars map[string]string) []string {
	var githubEnvVars []string
	for key, value := range envVars {
		if strings.HasPrefix(key, GITHUB_ENVVAR_PREFIX) {
			githubEnvVars = append(githubEnvVars, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return githubEnvVars
}
