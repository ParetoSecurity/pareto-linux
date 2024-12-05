package team

import (
	"os"
	"testing"

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
	if testing.Testing() {
		return &LinkingDevice{
			Hostname:  "test-hostname",
			OS:        "test-os",
			OSVersion: "test-os-version",
			Kernel:    "test-kernel",
			UUID:      "test-uuid",
			Ticket:    "test-ticket",
		}, nil
	}

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
