package main

import (
	"os"
	"strings"

	"github.com/okteto/test/pkg/cert"
	"github.com/okteto/test/pkg/command"
	"github.com/okteto/test/pkg/input"
	"github.com/okteto/test/pkg/logger"
	"github.com/spf13/afero"
)

func main() {
	args := os.Args[1:]
	environ := os.Environ()
	envVars := getEnvVars(environ)
	fs := afero.NewOsFs()

	userInput, err := input.NewInput(args, envVars)
	if err != nil {
		log := logger.NewLogger(logger.LevelError)
		log.LogError("Error parsing input", err)
		os.Exit(1)
	}

	logLevel := logger.ParseLogLevel(userInput.LogLevel)
	log := logger.NewLogger(logLevel)

	log.LogInfo("Starting execution...")

	runner := &command.DefaultRunner{}
	log.LogDebug("Using default command runner")
	if err := cert.HandleCaCert(userInput.CaCert, runner, fs); err != nil {
		log.LogError("Error handling CA certificate", err)
		os.Exit(2) // Return error code for certificate handling failure
	}
	log.LogDebug("CA certificate handled successfully")

	log.LogDebug("parsing input to params")
	params := userInput.ToParams()
	log.LogDebug("Input parsed successfully")
	log.LogDebug("params: ", params)
	log.LogDebug("Preparing and running command...")
	if err := command.PrepareAndRunCommand(params, runner, log); err != nil {
		log.LogError("Command execution failed", err)
		os.Exit(3) // Return error code for command execution failure
	}
	log.LogInfo("Execution completed successfully.")
}

func getEnvVars(environList []string) map[string]string {
	envVars := make(map[string]string)
	for _, env := range environList {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		envVars[parts[0]] = parts[1]
	}
	return envVars
}
