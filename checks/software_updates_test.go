package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/stretchr/testify/assert"
)

func TestCheckUpdates(t *testing.T) {

	tests := []struct {
		name           string
		setupMocks     map[string]string
		expectedPassed bool
		expectedDetail string
	}{
		{
			name: "All up to date",
			setupMocks: map[string]string{
				"flatpak remote-ls --updates": "",
				"apt list --upgradable":       "",
				"dnf check-update --quiet":    "",
				"pacman -Qu":                  "",
				"snap refresh --list":         "",
			},
			expectedPassed: true,
			expectedDetail: "All packages are up to date",
		},
		{
			name: "Updates available",
			setupMocks: map[string]string{
				"flatpak remote-ls --updates": "some updates",
				"apt list --upgradable":       "upgradable, upgradable",
				"dnf check-update --quiet":    "some updates",
				"pacman -Qu":                  "some updates",
				"snap refresh --list":         "some updates",
			},
			expectedPassed: false,
			expectedDetail: "Updates available for: Flatpak, APT, Pacman, Snap",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = tt.setupMocks
			lookPathMock = func(file string) (string, error) {
				return file, nil
			}
			su := &SoftwareUpdates{}
			passed, detail := su.checkUpdates()
			assert.Equal(t, tt.expectedPassed, passed)
			assert.Equal(t, tt.expectedDetail, detail)
			assert.NotEmpty(t, su.UUID())
			assert.False(t, su.RequiresRoot())
		})
	}
}
