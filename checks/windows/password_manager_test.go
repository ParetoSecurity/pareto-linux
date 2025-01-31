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
				"C:\\Users\\TestUser/AppData/Local/1Password/app/8/1Password.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Bitwarden present",
			mockFiles: map[string]bool{
				"C:\\Users\\TestUser/AppData/Local/Programs/Bitwarden/Bitwarden.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePass present",
			mockFiles: map[string]bool{
				"C:\\Program Files (x86)/KeePass Password Safe 2/KeePass.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassXC present",
			mockFiles: map[string]bool{
				"C:\\Program Files/KeePassXC/KeePassXC.exe": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name:           "No password manager present",
			mockFiles:      map[string]bool{},
			expectedPassed: false,
			expectedStatus: "No password manager found",
		},
	}

	for _, tt := range tests {
		os.Setenv("USERPROFILE", "C:\\Users\\TestUser")
		os.Setenv("PROGRAMFILES", "C:\\Program Files")
		os.Setenv("PROGRAMFILES(X86)", "C:\\Program Files (x86)")
		t.Run(tt.name, func(t *testing.T) {
			osStatMock = tt.mockFiles
			pmc := &PasswordManagerCheck{}
			err := pmc.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, pmc.Passed())
			assert.Equal(t, tt.expectedStatus, pmc.Status())
		})
		os.Unsetenv("USERPROFILE")
		os.Unsetenv("PROGRAMFILES")
		os.Unsetenv("PROGRAMFILES(X86)")
	}
}
