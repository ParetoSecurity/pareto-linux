package shared

import (
	"errors"
	"testing"

	sharedG "github.com/ParetoSecurity/pareto-core/shared"
)

func TestIsKeyStrong(t *testing.T) {
	// Save original RunCommand and restore at end.

	tests := []struct {
		name       string
		output     string
		err        error
		expectPass bool
	}{
		{
			name:       "RSA meets requirement",
			output:     "2048 abc dummy RSA",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "RSA below requirement",
			output:     "2047 abc dummy RSA",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "DSA meets requirement",
			output:     "8192 abc dummy DSA",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "DSA below requirement",
			output:     "8191 abc dummy DSA",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "ECDSA meets requirement",
			output:     "521 abc dummy ECDSA",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "ECDSA below requirement",
			output:     "520 abc dummy ECDSA",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "Ed25519 meets requirement",
			output:     "256 abc dummy ED25519",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "Ed25519 below requirement",
			output:     "255 abc dummy ED25519",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "Unknown key type",
			output:     "1024 abc dummy UNKNOWN",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "RunCommand returns error",
			output:     "",
			err:        errors.New("command failed"),
			expectPass: false,
		},
		{
			name:       "Malformed output (less than 4 fields)",
			output:     "2048 abc",
			err:        nil,
			expectPass: false,
		},
	}

	algo := SSHKeysAlgo{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Override RunCommand to simulate different outputs.
			sharedG.RunCommandMocks = map[string]string{
				"ssh-keygen -l -f dummy/path": tc.output,
			}
			result := algo.isKeyStrong("dummy/path")
			if result != tc.expectPass {
				t.Errorf("expected %v, got %v for output %q", tc.expectPass, result, tc.output)
			}
		})
	}
}
