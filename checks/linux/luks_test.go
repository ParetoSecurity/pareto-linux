package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
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

func TestEncryptingFS_Run(t *testing.T) {
	tests := []struct {
		name           string
		mockFiles      map[string]string
		mockCommand    string
		mockCommandOut string
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "Encrypted device found in crypttab and blkid",
			mockFiles: map[string]string{
				"/etc/crypttab": "cryptroot UUID=1234-5678-90AB-CDEF none luks",
			},
			mockCommand:    "blkid",
			mockCommandOut: "/dev/sda1: UUID=\"1234-5678-90AB-CDEF\" TYPE=\"crypto_LUKS\"",
			expectedPassed: true,
			expectedStatus: "Block device encryption is enabled",
		},
		{
			name: "No encrypted device found in crypttab",
			mockFiles: map[string]string{
				"/etc/crypttab": "",
			},
			mockCommand:    "blkid",
			mockCommandOut: "/dev/sda1: UUID=\"1234-5678-90AB-CDEF\" TYPE=\"crypto_LUKS\"",
			expectedPassed: false,
			expectedStatus: "Block device encryption is disabled",
		},
		{
			name: "No encrypted device found in blkid",
			mockFiles: map[string]string{
				"/etc/crypttab": "cryptroot UUID=1234-5678-90AB-CDEF none luks",
			},
			mockCommand:    "blkid",
			mockCommandOut: "/dev/sda1: UUID=\"5678-90AB-CDEF-1234\" TYPE=\"ext4\"",
			expectedPassed: false,
			expectedStatus: "Block device encryption is disabled",
		},
		{
			name: "Encrypted device found via kernel parameters",
			mockFiles: map[string]string{
				"/etc/crypttab": "",
				"/proc/cmdline": "BOOT_IMAGE=/vmlinuz-linux cryptdevice=UUID=1234-5678-90AB-CDEF:cryptroot:root root=/dev/mapper/cryptroot",
			},
			mockCommand:    "blkid",
			mockCommandOut: "",
			expectedPassed: true,
			expectedStatus: "Block device encryption is enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock shared.ReadFile
			shared.ReadFileMocks = tt.mockFiles
			// Mock shared.RunCommand
			shared.RunCommandMocks = map[string]string{
				tt.mockCommand: tt.mockCommandOut,
			}

			e := &EncryptingFS{}
			err := e.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, e.Passed())
			assert.Equal(t, tt.expectedStatus, e.Status())
			assert.NotEmpty(t, e.UUID())
			assert.True(t, e.RequiresRoot())
		})
	}
}
