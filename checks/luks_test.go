package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/stretchr/testify/assert"
)

func TestMaybeCryptoViaKernel(t *testing.T) {
	tests := []struct {
		name     string
		cmdline  string
		expected bool
	}{
		{
			name:     "cryptdevice present and root",
			cmdline:  "BOOT_IMAGE=/vmlinuz-linux cryptdevice=UUID=1234-5678-90AB-CDEF:cryptroot:root root=/dev/mapper/cryptroot",
			expected: true,
		},
		{
			name:     "cryptdevice present but not root",
			cmdline:  "BOOT_IMAGE=/vmlinuz-linux cryptdevice=UUID=1234-5678-90AB-CDEF:cryptroot:other root=/dev/mapper/cryptroot",
			expected: false,
		},
		{
			name:     "cryptdevice not present",
			cmdline:  "BOOT_IMAGE=/vmlinuz-linux root=/dev/sda1",
			expected: false,
		},
		{
			name:     "empty cmdline",
			cmdline:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the shared.ReadFile function
			shared.ReadFileMocks = map[string]string{
				"/proc/cmdline": tt.cmdline,
			}

			result := maybeCryptoViaKernel()
			assert.Equal(t, tt.expected, result)
		})
	}
}
