package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
	"github.com/stretchr/testify/assert"
)

func TestCheckUFW(t *testing.T) {
	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedResult bool
	}{
		{
			name:           "UFW is active",
			mockOutput:     "Status: active",
			mockError:      nil,
			expectedResult: true,
		},
		{
			name:           "UFW is inactive",
			mockOutput:     "Status: inactive",
			mockError:      nil,
			expectedResult: false,
		},
		{
			name:           "UFW command error",
			mockOutput:     "",
			mockError:      assert.AnError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = map[string]string{
				"ufw status": tt.mockOutput,
			}
			f := &Firewall{}
			result := f.checkUFW()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestCheckFirewalld(t *testing.T) {
	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedResult bool
	}{
		{
			name:           "Firewalld is active",
			mockOutput:     "active",
			mockError:      nil,
			expectedResult: true,
		},
		{
			name:           "Firewalld is inactive",
			mockOutput:     "inactive",
			mockError:      nil,
			expectedResult: false,
		},
		{
			name:           "Firewalld command error",
			mockOutput:     "",
			mockError:      assert.AnError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = map[string]string{
				"systemctl is-active firewalld": tt.mockOutput,
			}

			f := &Firewall{}
			result := f.checkFirewalld()
			assert.Equal(t, tt.expectedResult, result)
			assert.NotEmpty(t, f.UUID())
			assert.True(t, f.RequiresRoot())
		})
	}
}

func TestFirewall_Run(t *testing.T) {
	tests := []struct {
		name           string
		mockUFWOutput  string
		mockFirewalldOutput string
		expectedPassed bool
		expectedStatus string
	}{
		{
			name:           "UFW is active",
			mockUFWOutput:  "Status: active",
			mockFirewalldOutput: "",
			expectedPassed: true,
			expectedStatus: "Firewall is on",
		},
		{
			name:           "Firewalld is active",
			mockUFWOutput:  "Status: inactive",
			mockFirewalldOutput: "active",
			expectedPassed: true,
			expectedStatus: "Firewall is on",
		},
		{
			name:           "Both UFW and Firewalld are inactive",
			mockUFWOutput:  "Status: inactive",
			mockFirewalldOutput: "inactive",
			expectedPassed: false,
			expectedStatus: "Firewall is off",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = map[string]string{
				"ufw status": tt.mockUFWOutput,
				"systemctl is-active firewalld": tt.mockFirewalldOutput,
			}

			f := &Firewall{}
			err := f.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, f.Passed())
			assert.Equal(t, tt.expectedStatus, f.Status())
		})
	}
}
