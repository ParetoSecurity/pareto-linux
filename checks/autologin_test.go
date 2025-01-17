package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/stretchr/testify/assert"
)

func TestAutologin_Run(t *testing.T) {
	tests := []struct {
		name           string
		mockFiles      map[string]string
		mockCommand    string
		mockCommandOut string
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "SDDM autologin enabled in conf.d",
			mockFiles: map[string]string{
				"/etc/sddm.conf.d/test.conf": "Autologin=true",
				"/etc/sddm.conf":             "Autologin=true",
			},
			expectedPassed: false,
			expectedStatus: "Autologin=true in SDDM is enabled",
		},
		{
			name: "SDDM autologin enabled in main config",
			mockFiles: map[string]string{
				"/etc/sddm.conf": "Autologin=true",
			},
			expectedPassed: false,
			expectedStatus: "Autologin=true in SDDM is enabled",
		},
		{
			name: "GDM autologin enabled in custom.conf",
			mockFiles: map[string]string{
				"/etc/gdm3/custom.conf": "AutomaticLoginEnable=true",
			},
			expectedPassed: false,
			expectedStatus: "AutomaticLoginEnable=true in GDM is enabled",
		},
		{
			name:           "GDM autologin enabled in dconf",
			mockCommand:    "dconf read /org/gnome/login-screen/enable-automatic-login",
			mockCommandOut: "false",
			expectedPassed: true,
			expectedStatus: "Automatic login is off",
		},
		{
			name:           "No autologin enabled",
			expectedPassed: true,
			expectedStatus: "Automatic login is off",
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

			a := &Autologin{}
			err := a.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, a.Passed())
			assert.Equal(t, tt.expectedStatus, a.Status())
			assert.NotEmpty(t, a.UUID())
			assert.False(t, a.RequiresRoot())
		})
	}
}
