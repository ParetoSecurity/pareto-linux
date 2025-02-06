package checks

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
	"github.com/stretchr/testify/assert"
)

func TestDockerAccess_Run(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		expectedPassed bool
		expectedStatus string
	}{
		{
			name:           "Docker info command fails",
			commandOutput:  "",
			expectedPassed: false,
			expectedStatus: "Failed to get Docker info",
		},
		{
			name:           "Docker not running in rootless mode",
			commandOutput:  "seccomp",
			expectedPassed: false,
			expectedStatus: "Docker is not running in rootless mode",
		},
		{
			name:           "Docker running in rootless mode",
			commandOutput:  "rootless",
			expectedPassed: true,
			expectedStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = map[string]string{
				"docker version": "1.0.0",
				"docker info --format {{.SecurityOptions}}": tt.commandOutput,
			}
			dockerAccess := &DockerAccess{}
			err := dockerAccess.Run()

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, dockerAccess.passed)
			assert.Equal(t, tt.expectedStatus, dockerAccess.status)
			assert.NotEmpty(t, dockerAccess.UUID())
			assert.False(t, dockerAccess.RequiresRoot())
		})
	}
}

func TestDockerAccess_IsRunnable(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		expectedResult bool
		expectedStatus string
	}{
		{
			name:           "Docker is installed",
			commandOutput:  "Docker version 20.10.7, build f0df350",
			expectedResult: true,
			expectedStatus: "",
		},
		{
			name:           "Docker is not installed",
			commandOutput:  "",
			expectedResult: false,
			expectedStatus: "Docker is not installed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = map[string]string{
				"docker version": tt.commandOutput,
			}
			dockerAccess := &DockerAccess{}
			result := dockerAccess.IsRunnable()

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedStatus, dockerAccess.status)
		})
	}
}
