package shared

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/caarlos0/log"
	"github.com/google/uuid"
	"github.com/zcalusic/sysinfo"
)

type ReportingDevice struct {
	MachineUUID    string `json:"machineUUID"` // e.g. 123e4567-e89b-12d3-a456-426614174000
	MachineName    string `json:"machineName"` // e.g. MacBook-Pro.local
	Auth           string `json:"auth"`
	LinuxOSVersion string `json:"linuxOSVersion"` // e.g. Ubuntu 20.04
	ModelName      string `json:"modelName"`      // e.g. MacBook Pro
	ModelSerial    string `json:"modelSerial"`    // e.g. C02C1234
}

func CurrentReportingDevice() ReportingDevice {
	device, err := NewLinkingDevice()
	if err != nil {
		log.WithError(err).Fatal("Failed to get device information")
	}

	return ReportingDevice{
		MachineUUID:    device.UUID,
		MachineName:    device.Hostname,
		Auth:           DeviceAuth(),
		LinuxOSVersion: device.OS,
		ModelName: func() string {
			modelName, err := SystemDevice()
			if err != nil {
				return "Unknown"
			}

			return modelName
		}(),
		ModelSerial: func() string {
			serial, err := SystemSerial()
			if err != nil {
				return "Unknown"
			}

			return serial
		}(),
	}
}

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
	systemUUID, err := SystemUUID()
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

	if Config.AuthToken == "" {
		return ""
	}

	payload := Payload{}
	claims := strings.Split(Config.AuthToken, ".")[1]
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
