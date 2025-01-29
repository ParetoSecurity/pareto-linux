package team

import (
	"testing"

	"github.com/ParetoSecurity/pareto-core/shared"
	"github.com/stretchr/testify/assert"
)

func TestAddDevice(t *testing.T) {

	shared.Config.TeamID = "test-team-id"

	t.Run("successful device addition", func(t *testing.T) {

		err := AddDevice()
		assert.NoError(t, err)
	})

}
