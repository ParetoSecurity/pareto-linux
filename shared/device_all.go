package shared

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elastic/go-sysinfo"
	"github.com/google/uuid"
)

func CurrentReportingDevice() ReportingDevice {
	device, err := NewLinkingDevice()
	if err != nil {
		log.WithError(err).Fatal("Failed to get device information")
	}

	osVersion := device.OS
	if runtime.GOOS == "windows" {
		osVersion = strings.ReplaceAll(osVersion, "Microsoft", "")
		osVersion = fmt.Sprintf("%s %s", osVersion, device.OSVersion)
	}

	osVersion = Sanitize(fmt.Sprintf("%s %s", osVersion, device.OSVersion))

	return ReportingDevice{
		MachineUUID: device.UUID,
		MachineName: Sanitize(device.Hostname),
		Auth:        DeviceAuth(),
		OSVersion:   osVersion,
		ModelName: func() string {
			modelName, err := SystemDevice()
			if err != nil {
				return "Unknown"
			}

			return Sanitize(modelName)
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

	hostInfo, err := sysinfo.Host()
	if err != nil {
		log.Warn("Failed to get process information")
		return nil, err
	}
	envInfo := hostInfo.Info()

	systemUUID, err := SystemUUID()
	if err != nil {
		log.Warn("Failed to get system UUID")
		return nil, err
	}
	ticket, err := uuid.NewRandom()
	if err != nil {
		log.Warn("Failed to generate ticket")
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Warn("Failed to get hostname")
		return nil, err
	}

	return &LinkingDevice{
		Hostname:  hostname,
		OS:        envInfo.OS.Name,
		OSVersion: envInfo.OS.Version,
		Kernel:    envInfo.OS.Build,
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
