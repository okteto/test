package cert

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// fakeCommandRunner is a mock implementation of the CommandRunner interface
type fakeCommandRunner struct {
	err error
}

// Run does nothing
func (fcr *fakeCommandRunner) Run(_ *exec.Cmd) error {
	return fcr.err
}

// fakeInfoLogger is a mock implementation of the InfoLogger interface
type fakeInfoLogger struct{}

// Info does nothing
func (fil *fakeInfoLogger) Info(_ string, _ ...interface{}) {}

func TestHandleCaCert(t *testing.T) {
	tests := []struct {
		name        string
		caCert      string
		setupFs     func() afero.Fs
		setupRunner func() *fakeCommandRunner
		expectError error
	}{
		{
			name:   "empty caCert",
			caCert: "",
			setupFs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			setupRunner: func() *fakeCommandRunner {
				return &fakeCommandRunner{}
			},
			expectError: nil,
		},
		{
			name:   "write CA cert failure",
			caCert: "dummy-cert",
			setupFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				return afero.NewReadOnlyFs(fs) // Make the filesystem read-only to simulate write failure
			},
			setupRunner: func() *fakeCommandRunner {
				return &fakeCommandRunner{}
			},
			expectError: ErrWriteFailed,
		},
		{
			name:   "update CA certificates failure",
			caCert: "dummy-cert",
			setupFs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			setupRunner: func() *fakeCommandRunner {
				runner := &fakeCommandRunner{}
				runner.err = errors.New("update failed")
				return runner
			},
			expectError: ErrUpdateFailed,
		},
		{
			name:   "successful case",
			caCert: "dummy-cert",
			setupFs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			setupRunner: func() *fakeCommandRunner {
				runner := &fakeCommandRunner{}
				return runner
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setupFs()
			runner := tt.setupRunner()

			err := HandleCaCert(tt.caCert, runner, fs, &fakeInfoLogger{})

			assert.ErrorIs(t, err, tt.expectError)
		})
	}
}
