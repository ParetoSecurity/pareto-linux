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
