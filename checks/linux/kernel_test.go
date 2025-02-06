package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
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
		{
			name:           "Helper function error",
			isRoot:         false,
			runCheckHelper: true,
			helperError:    assert.AnError,
			expectedPassed: false,
			expectedStatus: "Failed to run check via helper",
		},
		{
			name:           "Sysctl command error",
			isRoot:         true,
			sysctlError:    assert.AnError,
			expectedPassed: false,
			expectedStatus: "Failed to get sysctl value",
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
			if tt.runCheckHelper {
				shared.RunCheckViaHelperMock = func(uuid string) (bool, error) {
					return false, tt.helperError
				}
			}
			err := k.Run()
			if tt.sysctlError != nil || tt.helperError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedPassed, k.Passed())
			assert.Equal(t, tt.expectedStatus, k.Status())
			assert.NotEmpty(t, k.UUID())
			assert.True(t, k.RequiresRoot())
		})
	}
}

func TestKernelParamsCheck_Name(t *testing.T) {
	k := &KernelParamsCheck{}
	expectedName := "Kernel Parameters are set correctly"
	if k.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, k.Name())
	}
}

func TestKernelParamsCheck_Status(t *testing.T) {
	k := &KernelParamsCheck{}
	expectedStatus := "Critical kernel parameters are correct"
	if k.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, k.Status())
	}
}

func TestKernelParamsCheck_UUID(t *testing.T) {
	k := &KernelParamsCheck{}
	expectedUUID := "cbf2736b-72df-43e3-8789-8eb676ff9014"
	if k.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, k.UUID())
	}
}

func TestKernelParamsCheck_Passed(t *testing.T) {
	k := &KernelParamsCheck{passed: true}
	expectedPassed := true
	if k.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, k.Passed())
	}
}

func TestKernelParamsCheck_FailedMessage(t *testing.T) {
	k := &KernelParamsCheck{}
	expectedFailedMessage := "Critical kernel parameters are not correct"
	if k.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, k.FailedMessage())
	}
}

func TestKernelParamsCheck_PassedMessage(t *testing.T) {
	k := &KernelParamsCheck{}
	expectedPassedMessage := "Critical kernel parameters are correct"
	if k.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, k.PassedMessage())
	}
}
