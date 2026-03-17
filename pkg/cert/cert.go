// pkg/cert/certificate.go
package cert

import (
	"fmt"
	"os/exec"

	"github.com/okteto/test/pkg/command"
	"github.com/spf13/afero"
)

const (
	// oktetoCACertPath is the path to the CA certificate
	oktetoCACertPath = "/etc/ssl/certs/okteto_ca_cert.pem"

	// oktetoCACertDebianPath is the path used when update-ca-certificates is available
	oktetoCACertDebianPath = "/usr/local/share/ca-certificates/okteto_ca_cert.crt"

	// updateCaCertificatesCmd is the command to update the CA certificates
	updateCaCertificatesCmd = "update-ca-certificates"
)

var (
	// ErrUpdateFailed is returned when the update of the CA certificates fails
	ErrUpdateFailed = fmt.Errorf("failed to update CA certificates")

	// ErrWriteFailed is returned when the write of the CA certificate fails
	ErrWriteFailed = fmt.Errorf("failed to write CA cert")
)

// LookPathFunc is the signature for exec.LookPath
type LookPathFunc func(string) (string, error)

// InfoLogger is an interface for logging information
type InfoLogger interface {
	Info(msg string, keysAndValues ...interface{})
}

// HandleCaCert writes the CA certificate to the system and updates certificates
func HandleCaCert(caCert string, runner command.CommandRunner, lookPath LookPathFunc, fs afero.Fs, l InfoLogger) error {
	if caCert == "" {
		l.Info("No CA certificate provided")
		return nil
	}

	if err := afero.WriteFile(fs, oktetoCACertPath, []byte(caCert), 0644); err != nil {
		return fmt.Errorf("%w: %w", ErrWriteFailed, err)
	}
	l.Info("CA certificate written successfully")

	if _, err := lookPath(updateCaCertificatesCmd); err == nil {
		if err := fs.MkdirAll("/usr/local/share/ca-certificates", 0755); err != nil {
			return fmt.Errorf("%w: %w", ErrWriteFailed, err)
		}
		if err := afero.WriteFile(fs, oktetoCACertDebianPath, []byte(caCert), 0644); err != nil {
			return fmt.Errorf("%w: %w", ErrWriteFailed, err)
		}
		cmd := exec.Command(updateCaCertificatesCmd)
		if err := runner.Run(cmd); err != nil {
			return fmt.Errorf("%w: %w", ErrUpdateFailed, err)
		}
		l.Info("CA certificates updated successfully")
	}

	return nil
}
