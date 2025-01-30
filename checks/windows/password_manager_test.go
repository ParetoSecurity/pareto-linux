package checks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordManagerCheck_Run(t *testing.T) {
	tests := []struct {
		name           string
		mockFiles      map[string]bool
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "1Password present",
			mockFiles: map[string]bool{
				"C:\\Program Files\\1Password\\1Password.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Bitwarden present",
			mockFiles: map[string]bool{
				"C:\\Program Files\\Bitwarden\\Bitwarden.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Dashlane present",
			mockFiles: map[string]bool{
				"C:\\Program Files\\Dashlane\\Dashlane.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassX present",
			mockFiles: map[string]bool{
				"C:\\Program Files\\KeePassX\\KeePassX.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassXC present",
			mockFiles: map[string]bool{
				"C:\\Program Files\\KeePassXC\\KeePassXC.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "No password manager present",
			mockFiles: map[string]bool{
				"C:\\Program Files\\1Password\\1Password.exe": false,
				"C:\\Program Files\\Bitwarden\\Bitwarden.exe": false,
				"C:\\Program Files\\Dashlane\\Dashlane.exe":   false,
				"C:\\Program Files\\KeePassX\\KeePassX.exe":   false,
				"C:\\Program Files\\KeePassXC\\KeePassXC.exe": false,
			},
			expectedPassed: false,
			expectedStatus: "No password manager found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.Stat
			statMock := func(name string) (os.FileInfo, error) {
				if tt.mockFiles[name] {
					return nil, nil
				}
				return nil, os.ErrNotExist
			}
			osStat = statMock

			pmc := &PasswordManagerCheck{}
			err := pmc.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, pmc.Passed())
			assert.Equal(t, tt.expectedStatus, pmc.Status())
		})
	}
}
