// pkg/cert/certificate.go
package cert

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/okteto/test/pkg/command"
	"github.com/spf13/afero"
)

const (
	// oktetocaCertPath is the path to the CA certificate
	oktetocaCertPath = "/etc/ssl/certs/okteto_ca_cert.pem"

	// updateCaCertificatesCmd is the command to update the CA certificates
	updateCaCertificatesCmd = "update-ca-certificates"
)

var (
	// ErrUpdateFailed is returned when the update of the CA certificates fails
	ErrUpdateFailed = fmt.Errorf("failed to update CA certificates")

	// ErrWriteFailed is returned when the write of the CA certificate fails
	ErrWriteFailed = fmt.Errorf("failed to write CA cert")
)

// InfoLogger is an interface for logging information
type InfoLogger interface {
	Info(msg string, keysAndValues ...interface{})
}

// HandleCaCert writes the CA certificate to the system and updates certificates
func HandleCaCert(caCert string, runner command.CommandRunner, fs afero.Fs, l InfoLogger) error {
	if caCert == "" {
		l.Info("No CA certificate provided")
		return nil
	}

	err := afero.WriteFile(fs, oktetocaCertPath, []byte(caCert), 0644)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWriteFailed, err)
	}
	l.Info("CA certificate written successfully")

	cmd := exec.Command(updateCaCertificatesCmd)
	if err := runner.Run(cmd); err != nil {
		var execErr *exec.Error
		if errors.As(err, &execErr) && errors.Is(execErr.Err, exec.ErrNotFound) {
			l.Info("update-ca-certificates not found, skipping")
			return nil
		}
		return fmt.Errorf("%w: %w", ErrUpdateFailed, err)
	}
	l.Info("CA certificates updated successfully")
	return nil
}
