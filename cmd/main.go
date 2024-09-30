package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/okteto/test/pkg/cert"
	"github.com/okteto/test/pkg/command"
	"github.com/okteto/test/pkg/input"
	"github.com/spf13/afero"
)

func main() {
	args := os.Args[1:]
	environ := os.Environ()
	envVars := getEnvVars(environ)
	fs := afero.NewOsFs()

	userInput, err := input.NewInput(args, envVars)
	if err != nil {
		logger := slog.New(&slog.TextHandler{})
		logger.Error("Error parsing input", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logLevel := parseLogLevel(userInput.LogLevel)
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	log := slog.New(textHandler)

	log.Info("Starting execution...")

	runner := &command.DefaultRunner{}
	log.Debug("Using default command runner")
	if err := cert.HandleCaCert(userInput.CaCert, runner, fs); err != nil {
		log.Error("Error handling CA certificate: %v", slog.String("error", err.Error()))
		os.Exit(2) // Return error code for certificate handling failure
	}
	log.Debug("CA certificate handled successfully")

	log.Debug("parsing input to params")
	params := userInput.ToParams()
	log.Debug("Input parsed successfully")
	log.Debug("params: ", slog.String("params", strings.Join(params, " ")))
	log.Debug("Preparing and running command...")
	if err := command.PrepareAndRunCommand(params, runner, log); err != nil {
		log.Error("Command execution failed ", slog.String("error", err.Error()))
		os.Exit(3) // Return error code for command execution failure
	}
	log.Info("Execution completed successfully.")
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

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
