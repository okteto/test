// pkg/command/command_test.go
package command

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRunner is a mock implementation of CommandRunner
type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) Run(cmd *exec.Cmd) error {
	args := m.Called(cmd)
	return args.Error(0)
}

type MockLogger struct {
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	fmt.Print(msg)
}

func TestPrepareAndRunCommand(t *testing.T) {
	tests := []struct {
		name        string
		params      []string
		expectedErr error
	}{
		{
			name:        "Success",
			params:      []string{"--name=my-app"},
			expectedErr: nil,
		},
		{
			name:        "Failure",
			params:      []string{"--name=my-app"},
			expectedErr: exec.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := new(MockRunner)
			params := append([]string{"test"}, tt.params...)
			cmd := exec.Command("okteto", params...)

			mockRunner.On("Run", cmd).Return(tt.expectedErr)

			err := PrepareAndRunCommand(tt.params, mockRunner, &MockLogger{})

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
