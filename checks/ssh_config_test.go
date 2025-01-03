package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/stretchr/testify/assert"
)

func TestCheckSSHConfig(t *testing.T) {

	tests := []struct {
		name           string
		setupMocks     map[string]string
		expectedPassed bool
		expectedDetail string
	}{
		{
			name: "All ok",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication no\nPermitRootLogin no",
			},
			expectedPassed: true,
			expectedDetail: "",
		},
		{
			name: "PasswordAuthentication is enabled",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication yes\nPermitRootLogin no",
			},
			expectedPassed: false,
			expectedDetail: "PasswordAuthentication is enabled",
		},
		{
			name: "PermitRootLogin is enabled",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication no\nPermitRootLogin yes",
			},
			expectedPassed: false,
			expectedDetail: "Root login is enabled",
		},

		{
			name: "PermitEmptyPasswords is enabled",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication no\nPermitRootLogin no\nPermitEmptyPasswords yes",
			},
			expectedPassed: false,
			expectedDetail: "Empty passwords are allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = tt.setupMocks
			lookPathMock = func(file string) (string, error) {
				return file, nil
			}
			su := &SSHConfigCheck{}

			err := su.Run()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedPassed, su.passed)
			assert.Equal(t, tt.expectedDetail, su.status)
		})
	}
}
