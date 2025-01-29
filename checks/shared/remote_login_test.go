package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteLogin_Run_NoOpenPorts(t *testing.T) {
	remoteLogin := &RemoteLogin{}

	// Mock checkPort to always return false
	checkPortMock = func(port int, proto string) bool {
		return false
	}

	err := remoteLogin.Run()
	assert.NoError(t, err)
	assert.True(t, remoteLogin.Passed())
	assert.Empty(t, remoteLogin.ports)
}

func TestRemoteLogin_Run_OpenPorts(t *testing.T) {
	remoteLogin := &RemoteLogin{}

	// Mock checkPort to return true for specific ports
	checkPortMock = func(port int, _ string) bool {
		return port == 22 || port == 3389
	}

	err := remoteLogin.Run()
	assert.NoError(t, err)
	assert.False(t, remoteLogin.Passed())
	assert.NotEmpty(t, remoteLogin.ports)
	assert.Contains(t, remoteLogin.ports, 22)
	assert.Contains(t, remoteLogin.ports, 3389)
	assert.NotContains(t, remoteLogin.ports, 3390)
	assert.NotContains(t, remoteLogin.ports, 5900)
	assert.NotEmpty(t, remoteLogin.UUID())
	assert.False(t, remoteLogin.RequiresRoot())
}
