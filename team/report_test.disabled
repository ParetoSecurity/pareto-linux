package team

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
	"github.com/stretchr/testify/assert"
)

func TestReportToTeam_Error(t *testing.T) {
	// Mock shared.Config and shared.HTTPTransport
	shared.Config.TeamID = "test-team-id"

	// Call the function
	err := ReportToTeam()

	// Assertions
	assert.NoError(t, err)

}
