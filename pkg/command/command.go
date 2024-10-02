// pkg/command/command.go
package command

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	// oktetoCmd is the command to execute
	oktetoCmd = "okteto"

	// subcommand is the subcommand to execute
	subcommand = "test"
)

// CommandRunner is an interface for running commands
type CommandRunner interface {
	Run(cmd *exec.Cmd) error
}

// InfoLogger is an interface for logging information
type InfoLogger interface {
	Info(msg string, keysAndValues ...interface{})
}

// DefaultRunner is the default implementation of CommandRunner
type DefaultRunner struct {
	Environ []string
}

// Run executes the provided command
func (r *DefaultRunner) Run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = r.Environ
	return cmd.Run()
}

// PrepareAndRunCommand prepares and executes the okteto command
func PrepareAndRunCommand(params []string, runner CommandRunner, l InfoLogger) error {
	execArgs := append([]string{subcommand}, params...)
	cmd := exec.Command(oktetoCmd, execArgs...)
	l.Info(fmt.Sprintf("Executing command: %s %v", oktetoCmd, execArgs))
	return runner.Run(cmd)
}
