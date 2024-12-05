package team

import (
	"encoding/base64"
	"testing"

	"paretosecurity.com/auditor/shared"
)

func TestDeviceAuth(t *testing.T) {
	// Mock the shared.Config.AuthToken
	shared.Config.AuthToken = "header." + base64.StdEncoding.EncodeToString([]byte(`{"sub":"1234567890","teamID":"team123","role":"admin","iat":1516239022,"token":"test-token"}`)) + ".signature"

	expectedToken := "test-token"
	actualToken := DeviceAuth()

	if actualToken != expectedToken {
		t.Errorf("expected %s, got %s", expectedToken, actualToken)
	}
}

func TestDeviceAuthInvalidToken(t *testing.T) {
	// Mock the shared.Config.AuthToken with an invalid token
	shared.Config.AuthToken = "invalid.token"

	expectedToken := ""
	actualToken := DeviceAuth()

	if actualToken != expectedToken {
		t.Errorf("expected %s, got %s", expectedToken, actualToken)
	}
}

func TestDeviceAuthMalformedPayload(t *testing.T) {
	// Mock the shared.Config.AuthToken with a malformed payload
	shared.Config.AuthToken = "header." + base64.StdEncoding.EncodeToString([]byte(`malformed-payload`)) + ".signature"

	expectedToken := ""
	actualToken := DeviceAuth()

	if actualToken != expectedToken {
		t.Errorf("expected %s, got %s", expectedToken, actualToken)
	}
}
