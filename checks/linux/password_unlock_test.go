package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
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
			assert.NotEmpty(t, f.UUID())
			assert.False(t, f.RequiresRoot())
		})
	}
}

func TestPasswordToUnlock_Run(t *testing.T) {
	tests := []struct {
		name           string
		mockCommands   map[string]string
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "GNOME lock enabled",
			mockCommands: map[string]string{
				"gsettings get org.gnome.desktop.screensaver lock-enabled": "true\n",
			},
			expectedPassed: true,
			expectedStatus: "Password after sleep or screensaver is on",
		},
		{
			name: "GNOME lock disabled",
			mockCommands: map[string]string{
				"gsettings get org.gnome.desktop.screensaver lock-enabled": "false\n",
			},
			expectedPassed: false,
			expectedStatus: "Password after sleep or screensaver is off",
		},
		{
			name: "KDE autolock enabled",
			mockCommands: map[string]string{
				"kreadconfig5 --file kscreenlockerrc --group Daemon --key Autolock": "true\n",
			},
			expectedPassed: true,
			expectedStatus: "Password after sleep or screensaver is on",
		},
		{
			name: "KDE autolock disabled",
			mockCommands: map[string]string{
				"kreadconfig5 --file kscreenlockerrc --group Daemon --key Autolock": "false\n",
			},
			expectedPassed: false,
			expectedStatus: "Password after sleep or screensaver is off",
		},
		{
			name: "Neither GNOME nor KDE found",
			mockCommands: map[string]string{},
			expectedPassed: false,
			expectedStatus: "Password after sleep or screensaver is off",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = tt.mockCommands

			f := &PasswordToUnlock{}
			err := f.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, f.Passed())
			assert.Equal(t, tt.expectedStatus, f.Status())
		})
	}
}
