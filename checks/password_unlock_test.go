package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/stretchr/testify/assert"
)

func TestCheckKDE(t *testing.T) {
	tests := []struct {
		name       string
		commandOut string
		commandErr error
		expected   bool
	}{
		{
			name:       "Autolock enabled",
			commandOut: "true\n",
			commandErr: nil,
			expected:   true,
		},
		{
			name:       "Autolock disabled",
			commandOut: "false\n",
			commandErr: nil,
			expected:   false,
		},
		{
			name:       "Command error",
			commandOut: "",
			commandErr: assert.AnError,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = map[string]string{
				"kreadconfig5 --file kscreenlockerrc --group Daemon --key Autolock": tt.commandOut,
			}

			f := &PasswordToUnlock{}
			result := f.checkKDE()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestCheckGnome(t *testing.T) {
	tests := []struct {
		name       string
		commandOut string
		commandErr error
		expected   bool
	}{
		{
			name:       "Lock enabled",
			commandOut: "true\n",
			commandErr: nil,
			expected:   true,
		},
		{
			name:       "Lock disabled",
			commandOut: "false\n",
			commandErr: nil,
			expected:   false,
		},
		{
			name:       "Command error",
			commandOut: "",
			commandErr: assert.AnError,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = map[string]string{
				"gsettings get org.gnome.desktop.screensaver lock-enabled": tt.commandOut,
			}

			f := &PasswordToUnlock{}
			result := f.checkGnome()
			assert.Equal(t, tt.expected, result)
		})
	}
}
