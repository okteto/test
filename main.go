package main

import (
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

func (i *Input) nameToParams() string {
	if i.name != "" {
		return "--name " + i.name
	}
	return ""
}

func (i *Input) namespaceToParams() string {
	if i.namespace != "" {
		return "--namespace " + i.namespace
	}
	return ""
}

func (i *Input) fileToParams() string {
	if i.file != "" {
		return "--file " + i.file
	}
	return ""
}

func (i *Input) deployToParams() string {
	if i.deploy == "true" {
		return "--deploy "
	}
	return ""
}

func (i *Input) noCacheToParams() string {
	if i.no_cache == "true" {
		return "--no-cache "
	}
	return ""
}

func (i *Input) variablesToParams() string {
	if i.variables == "" {
		return ""
	}
	variables := strings.Split(i.variables, ",")
	params := ""
	for _, variable := range variables {
		params += "--variable " + variable + " "
	}
	return strings.Trim(params, " ")
}

func (i *Input) githubEnvVarsToParams() string {
	envVars := os.Environ()
	params := ""
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "GITHUB_") {
			params += "--variable " + envVar + " "
		}
	}
	return strings.Trim(params, " ")
}

func (i *Input) timeoutToParams() string {
	if i.timeout != "" {
		return "--timeout " + i.timeout
	}
	return ""
}

func (i *Input) testsToParams() string {
	if i.tests != "" {
		return i.tests
	}
	return ""
}

func (i *Input) logLevelToParams() string {
	if os.Getenv("RUNNER_DEBUG") == "1" {
		return "--log-level debug"
	}
	if i.log_level != "" {
		return "--log-level " + i.log_level
	}
	return ""
}

func main() {

	input := getInput()
	params := ""

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

	if input.name != "" {
		params += "--name " + input.name + " "
	}
	if input.namespace != "" {
		params += "--namespace " + input.namespace + " "
	}
	if input.file != "" {
		params += "--file " + input.file + " "
	}
	if input.deploy == "true" {
		params += "--deploy "
	}

	if input.no_cache == "true" {
		params += "--no-cache "
	}

	if input.variables != "" {
		params += input.variables + " "
	}

	if githubEnvVars := input.githubEnvVarsToParams(); githubEnvVars != "" {
		params += githubEnvVars + " "
	}

	if input.timeout != "" {
		params += "--timeout " + input.timeout + " "
	}

	if input.log_level != "" {
		params += "--log-level " + input.log_level + " "
	}
	if input.tests != "" {
		params += input.tests + " "
	}

	cmd := exec.Command("okteto", "test", params)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
