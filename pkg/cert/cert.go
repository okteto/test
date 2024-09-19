// pkg/cert/certificate.go
package cert

import (
	"fmt"
	"os/exec"

	"github.com/okteto/test/pkg/command"
	"github.com/spf13/afero"
)

const (
	oktetocaCertPath = "/usr/local/share/ca-certificates/okteto_ca_cert.crt"

	updateCaCertificatesCmd = "update-ca-certificates"
)

var (
	// ErrUpdateFailed is returned when the update of the CA certificates fails
	ErrUpdateFailed = fmt.Errorf("failed to update CA certificates")

	// ErrWriteFailed is returned when the write of the CA certificate fails
	ErrWriteFailed = fmt.Errorf("failed to write CA cert")
)

// HandleCaCert writes the CA certificate to the system and updates certificates
func HandleCaCert(caCert string, runner command.CommandRunner, fs afero.Fs) error {
	if caCert == "" {
		return nil
	}

	err := afero.WriteFile(fs, oktetocaCertPath, []byte(caCert), 0644)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWriteFailed, err)
	}

	cmd := exec.Command(updateCaCertificatesCmd)
	if err := runner.Run(cmd); err != nil {
		return fmt.Errorf("%w: %w", ErrUpdateFailed, err)
	}
	return nil
}
