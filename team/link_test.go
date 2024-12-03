package team

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkAndWaitForTicket(t *testing.T) {

	// Running the test
	err := LinkAndWaitForTicket()
	assert.NoError(t, err)

}

func TestLinkAndWaitForTicket_NewLinkingDeviceError(t *testing.T) {

	// Running the test
	err := LinkAndWaitForTicket()
	assert.Error(t, err)
	assert.Equal(t, "new linking device error", err.Error())
}

func TestLinkAndWaitForTicket_EnrollError(t *testing.T) {

	// Running the test
	err := LinkAndWaitForTicket()
	assert.Error(t, err)
	assert.Equal(t, "enroll error", err.Error())

}
