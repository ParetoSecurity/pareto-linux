package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
	"github.com/stretchr/testify/assert"
)

func TestPasswordManagerCheck_Run_Linux(t *testing.T) {
	tests := []struct {
		name           string
		mockCommands   map[string]string
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "1Password present via apt",
			mockCommands: map[string]string{
				"sh -c dpkg -l | grep 1password": "ii  1password  1.0  all  Password manager",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Bitwarden present via snap",
			mockCommands: map[string]string{
				"sh -c snap list | grep bitwarden": "bitwarden  1.0  stable  password manager",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Dashlane present via yum",
			mockCommands: map[string]string{
				"sh -c yum list installed | grep dashlane": "dashlane  1.0  installed  password manager",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassX present via flatpak",
			mockCommands: map[string]string{
				"sh -c flatpak list | grep keepassx": "keepassx  1.0  stable  password manager",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassXC present via apt",
			mockCommands: map[string]string{
				"sh -c dpkg -l | grep keepassxc": "ii  keepassxc  1.0  all  Password manager",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "No password manager present",
			mockCommands: map[string]string{
				"sh -c dpkg -l | grep 1password":           "",
				"sh -c snap list | grep bitwarden":         "",
				"sh -c yum list installed | grep dashlane": "",
				"sh -c flatpak list | grep keepassx":       "",
				"sh -c dpkg -l | grep keepassxc":           "",
			},
			expectedPassed: false,
			expectedStatus: "No password manager found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock shared.RunCommand

			shared.RunCommandMocks = tt.mockCommands

			pmc := &PasswordManagerCheck{}
			err := pmc.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, pmc.Passed())
			assert.Equal(t, tt.expectedStatus, pmc.Status())
		})
	}
}
