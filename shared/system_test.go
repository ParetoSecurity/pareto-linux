package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemDevice_Success(t *testing.T) {
	ReadFileMocks = map[string]string{
		"/sys/devices/virtual/dmi/id/product_name": "TestDeviceName",
	}
	expected := "TestDeviceName"
	deviceName, err := SystemDevice()
	assert.NoError(t, err)
	assert.Equal(t, expected, deviceName)
}

func TestSystemDevice_EmptyContent(t *testing.T) {
	ReadFileMocks = map[string]string{
		"/sys/devices/virtual/dmi/id/product_name": "",
	}

	_, err := SystemDevice()
	assert.Error(t, err)
	assert.Equal(t, "unable to retrieve device name", err.Error())
}
func TestSystemSerial_Success(t *testing.T) {
	ReadFileMocks = map[string]string{
		"/sys/devices/virtual/dmi/id/product_serial": "TestSerialNumber",
	}
	expected := "TestSerialNumber"
	serialNumber, err := SystemSerial()
	assert.NoError(t, err)
	assert.Equal(t, expected, serialNumber)
}

func TestSystemSerial_EmptyContent(t *testing.T) {
	ReadFileMocks = map[string]string{
		"/sys/devices/virtual/dmi/id/product_serial": "",
	}

	_, err := SystemSerial()
	assert.Error(t, err)
	assert.Equal(t, "unable to retrieve serial number", err.Error())
}
