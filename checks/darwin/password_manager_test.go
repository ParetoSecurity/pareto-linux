package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckInstalledApplications(t *testing.T) {

	tests := []struct {
		name     string
		appNames []string
		expected bool
	}{
		{
			name:     "Password manager present",
			appNames: []string{"1Password.app", "Bitwarden.app"},
			expected: true,
		},
		{
			name:     "Password manager not present",
			appNames: []string{"NonExistentApp.app"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkInstalledApplications(tt.appNames)
			assert.Equal(t, tt.expected, result)
		})
	}
}
