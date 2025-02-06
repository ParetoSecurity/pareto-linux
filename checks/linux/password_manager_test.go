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
				"which dpkg":    "/usr/bin/dpkg",
				"sh -c dpkg -l": "ii  1password  1.0  all  Password manager",
				"which snap":    "not found",
				"which yum":     "not found",
				"which flatpak": "not found",
				"which pacman":  "not found",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Bitwarden present via snap",
			mockCommands: map[string]string{
				"which dpkg":      "not found",
				"which snap":      "/usr/bin/snap",
				"sh -c snap list": "bitwarden  1.0  stable  password manager",
				"which yum":       "not found",
				"which flatpak":   "not found",
				"which pacman":    "not found",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Dashlane present via yum",
			mockCommands: map[string]string{
				"which dpkg":               "not found",
				"which snap":               "not found",
				"which yum":                "/usr/bin/yum",
				"sh -c yum list installed": "dashlane  1.0  installed  password manager",
				"which flatpak":            "not found",
				"which pacman":             "not found",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassX present via flatpak",
			mockCommands: map[string]string{
				"which dpkg":         "not found",
				"which snap":         "not found",
				"which yum":          "not found",
				"which flatpak":      "/usr/bin/flatpak",
				"sh -c flatpak list": "keepassx  1.0  stable  password manager",
				"which pacman":       "not found",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "KeePassXC present via apt",
			mockCommands: map[string]string{
				"which dpkg":    "/usr/bin/dpkg",
				"sh -c dpkg -l": "ii  keepassxc  1.0  all  Password manager",
				"which snap":    "not found",
				"which yum":     "not found",
				"which flatpak": "not found",
				"which pacman":  "not found",
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "No password manager present",
			mockCommands: map[string]string{
				"which dpkg":               "/usr/bin/dpkg",
				"sh -c dpkg -l":            "",
				"which snap":               "/usr/bin/snap",
				"sh -c snap list":          "",
				"which yum":                "/usr/bin/yum",
				"sh -c yum list installed": "",
				"which flatpak":            "not found",
				"which pacman":             "not found",
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
			status := pmc.isManagerInstalled()
			assert.Equal(t, tt.expectedPassed, status)
		})
	}
}

func TestPasswordManagerCheck_Run_BrowserExtensions(t *testing.T) {
	tests := []struct {
		name           string
		mockFileSystem map[string]bool
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "1Password extension present in Chrome",
			mockFileSystem: map[string]bool{
				"/home/user/.config/google-chrome/Default/Extensions/1password": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "Bitwarden extension present in Firefox",
			mockFileSystem: map[string]bool{
				"/home/user/.mozilla/firefox/bitwarden": true,
			},
			expectedPassed: true,
			expectedStatus: "Password manager is present",
		},
		{
			name: "No password manager extensions present",
			mockFileSystem: map[string]bool{
				"/home/user/.config/google-chrome/Default/Extensions/1password": false,
				"/home/user/.mozilla/firefox/bitwarden":                         false,
			},
			expectedPassed: false,
			expectedStatus: "No password manager found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.Stat
			osStat = func(name string) (os.FileInfo, error) {
				if tt.mockFileSystem[name] {
					return nil, nil
				}
				return nil, os.ErrNotExist
			}

			pmc := &PasswordManagerCheck{}
			err := pmc.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, pmc.Passed())
			assert.Equal(t, tt.expectedStatus, pmc.Status())
		})
	}
}
