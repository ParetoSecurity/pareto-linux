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
				"/snap/bin/1password": true,
				"/usr/bin/1password": true,
				"/usr/local/bin/1password": true,
				"/opt/1password/1password": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Bitwarden present",
			mockFiles: map[string]bool{
				"/snap/bin/bitwarden": true,
				"/usr/bin/bitwarden": true,
				"/usr/local/bin/bitwarden": true,
				"/opt/bitwarden/bitwarden": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Dashlane present",
			mockFiles: map[string]bool{
				"/snap/bin/dashlane": true,
				"/usr/bin/dashlane": true,
				"/usr/local/bin/dashlane": true,
				"/opt/dashlane/dashlane": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassX present",
			mockFiles: map[string]bool{
				"/snap/bin/keepassx": true,
				"/usr/bin/keepassx": true,
				"/usr/local/bin/keepassx": true,
				"/opt/keepassx/keepassx": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassXC present",
			mockFiles: map[string]bool{
				"/snap/bin/keepassxc": true,
				"/usr/bin/keepassxc": true,
				"/usr/local/bin/keepassxc": true,
				"/opt/keepassxc/keepassxc": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "No password manager present",
			mockFiles: map[string]bool{
				"/snap/bin/1password": false,
				"/snap/bin/bitwarden": false,
				"/snap/bin/dashlane":  false,
				"/snap/bin/keepassx":  false,
				"/snap/bin/keepassxc": false,
				"/usr/bin/1password": false,
				"/usr/bin/bitwarden": false,
				"/usr/bin/dashlane":  false,
				"/usr/bin/keepassx":  false,
				"/usr/bin/keepassxc": false,
				"/usr/local/bin/1password": false,
				"/usr/local/bin/bitwarden": false,
				"/usr/local/bin/dashlane":  false,
				"/usr/local/bin/keepassx":  false,
				"/usr/local/bin/keepassxc": false,
				"/opt/1password/1password": false,
				"/opt/bitwarden/bitwarden": false,
				"/opt/dashlane/dashlane":  false,
				"/opt/keepassx/keepassx":  false,
				"/opt/keepassxc/keepassxc": false,
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
