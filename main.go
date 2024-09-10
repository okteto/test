package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Input struct {
	name      string
	namespace string
	file      string
	deploy    string
	no_cache  string
	variables string
	timeout   string
	tests     string
	log_level string
	caCert    string
}

func getInput() Input {
	return Input{
		name:      os.Args[1],
		namespace: os.Args[2],
		file:      os.Args[3],
		deploy:    os.Args[4],
		no_cache:  os.Args[5],
		variables: os.Args[6],
		timeout:   os.Args[7],
		tests:     os.Args[8],
		log_level: os.Args[9],
		caCert:    os.Getenv("OKTETO_CA_CERT"),
	}
}
func (i *Input) nameToParams() []string {
	if i.name != "" {
		return []string{fmt.Sprintf("--name=%s", i.name)}
	}
	return []string{}
}

func (i *Input) namespaceToParams() []string {
	if i.namespace != "" {
		return []string{fmt.Sprintf("--namespace=%s", i.namespace)}
	}
	return []string{}
}

func (i *Input) fileToParams() []string {
	if i.file != "" {
		return []string{fmt.Sprintf("--file=%s", i.file)}
	}
	return []string{}
}

func (i *Input) deployToParams() []string {
	if i.deploy == "true" {
		return []string{"--deploy"}
	}
	return []string{}
}

func (i *Input) noCacheToParams() []string {
	if i.no_cache == "true" {
		return []string{"--no-cache"}
	}
	return []string{}
}

func (i *Input) variablesToParams() []string {
	if i.variables == "" {
		return []string{}
	}
	variables := strings.Split(i.variables, ",")
	params := []string{}
	for _, variable := range variables {
		params = append(params, fmt.Sprintf("--var=%s", variable))
	}
	return params
}

func (i *Input) githubEnvVarsToParams() []string {
	envVars := os.Environ()
	params := []string{}
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "GITHUB_") {
			params = append(params, fmt.Sprintf("--var=%s", envVar))
		}
	}
	return params
}

func (i *Input) timeoutToParams() []string {
	if i.timeout != "" {
		return []string{fmt.Sprintf("--timeout=%s", i.timeout)}
	}
	return []string{}
}

func (i *Input) testsToParams() []string {
	if i.tests != "" {
		return []string{i.tests}
	}
	return []string{}
}

func (i *Input) logLevelToParams() []string {
	if os.Getenv("RUNNER_DEBUG") == "1" {
		return []string{"--log-level=debug"}
	}
	if i.log_level != "" {
		return []string{fmt.Sprintf("--log-level=%s", i.log_level)}
	}
	return []string{}
}

func main() {

	input := getInput()

	if input.caCert != "" {
		err := os.WriteFile("/usr/local/share/ca-certificates/okteto_ca_cert.crt", []byte(input.caCert), 0644)
		if err != nil {
			panic(err)
		}
		cmd := exec.Command("update-ca-certificates")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	execArgs := []string{"test"}
	execArgs = append(execArgs, input.nameToParams()...)
	execArgs = append(execArgs, input.namespaceToParams()...)
	execArgs = append(execArgs, input.fileToParams()...)
	execArgs = append(execArgs, input.deployToParams()...)
	execArgs = append(execArgs, input.noCacheToParams()...)
	execArgs = append(execArgs, input.variablesToParams()...)
	execArgs = append(execArgs, input.githubEnvVarsToParams()...)
	execArgs = append(execArgs, input.timeoutToParams()...)
	execArgs = append(execArgs, input.logLevelToParams()...)
	execArgs = append(execArgs, input.testsToParams()...)

	cmd := exec.Command("okteto", execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
