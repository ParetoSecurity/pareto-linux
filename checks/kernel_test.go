package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/stretchr/testify/assert"
)

func TestKernelParamsCheck_Run(t *testing.T) {

	tests := []struct {
		name           string
		isRoot         bool
		runCheckHelper bool
		helperError    error
		sysctlValues   map[string]string
		sysctlError    error
		expectedPassed bool
		expectedStatus string
	}{
		{
			name:           "All parameters correct",
			isRoot:         true,
			sysctlValues:   map[string]string{"net.ipv4.tcp_syncookies": "1", "kernel.randomize_va_space": "2", "fs.protected_hardlinks": "1", "fs.protected_symlinks": "1"},
			expectedPassed: true,
		},
		{
			name:           "Some parameters incorrect",
			isRoot:         true,
			sysctlValues:   map[string]string{"net.ipv4.tcp_syncookies": "0", "kernel.randomize_va_space": "2", "fs.protected_hardlinks": "1", "fs.protected_symlinks": "0"},
			expectedPassed: false,
			expectedStatus: "net.ipv4.tcp_syncookies is set to 0 but should be 1. fs.protected_symlinks is set to 0 but should be 1. ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KernelParamsCheck{}
			shared.RunCommandMocks = map[string]string{
				"sysctl -n net.ipv4.tcp_syncookies":   tt.sysctlValues["net.ipv4.tcp_syncookies"],
				"sysctl -n kernel.randomize_va_space": tt.sysctlValues["kernel.randomize_va_space"],
				"sysctl -n fs.protected_hardlinks":    tt.sysctlValues["fs.protected_hardlinks"],
				"sysctl -n fs.protected_symlinks":     tt.sysctlValues["fs.protected_symlinks"],
			}
			err := k.Run()
			if tt.sysctlError != nil || tt.helperError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedPassed, k.Passed())
			assert.NotEmpty(t, k.UUID())
			assert.True(t, k.RequiresRoot())
		})
	}
}
