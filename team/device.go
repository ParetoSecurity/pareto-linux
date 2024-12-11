package team

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/caarlos0/log"
	"github.com/google/uuid"
	"github.com/zcalusic/sysinfo"
	"paretosecurity.com/auditor/shared"
)

type LinkingDevice struct {
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	OSVersion string `json:"osVersion"`
	Kernel    string `json:"kernel"`
	UUID      string `json:"uuid"`
	Ticket    string `json:"ticket"`
	Version   string `json:"version"`
}

// NewLinkingDevice creates a new instance of LinkingDevice with system information.
// It retrieves the system UUID and device ticket, and populates the LinkingDevice struct
// with the hostname, OS name, OS version, kernel version, UUID, and ticket.
// Returns a pointer to the LinkingDevice and an error if any occurs during the process.
func NewLinkingDevice() (*LinkingDevice, error) {

	sysinfo := sysinfo.SysInfo{}
	sysinfo.GetSysInfo()
	systemUUID, err := shared.SystemUUID()
	if err != nil {
		return nil, err
	}
	ticket, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &LinkingDevice{
		Hostname:  hostname,
		OS:        sysinfo.OS.Name,
		OSVersion: sysinfo.OS.Release,
		Kernel:    sysinfo.Kernel.Release,
		UUID:      systemUUID,
		Ticket:    ticket.String(),
	}, nil
}

// DeviceAuth decodes a JWT token from the shared configuration's AuthToken,
// extracts the payload, and returns the token string from the payload.
// It returns an empty string if there is an error during decoding or unmarshalling.
func DeviceAuth() string {
	type Payload struct {
		Sub    string `json:"sub"`
		TeamID string `json:"teamID"`
		Role   string `json:"role"`
		Iat    int    `json:"iat"`
		Token  string `json:"token"`
	}

	if shared.Config.AuthToken == "" {
		return ""
	}

	payload := Payload{}
	claims := strings.Split(shared.Config.AuthToken, ".")[1]
	token, err := base64.RawURLEncoding.DecodeString(claims)
	if err != nil {
		log.WithError(err).WithField("claims", claims).Warn("failed to decode claims")
		return ""
	}
	if err := json.Unmarshal(token, &payload); err != nil {
		log.WithError(err).WithField("claims", claims).Warn("failed to unmarshal claims")
		return ""
	}
	return payload.Token
}
