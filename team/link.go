package team

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqtest"
	"github.com/zcalusic/sysinfo"
	"paretosecurity.com/auditor/shared"
)

const baseURL = "https://dash.paretosecurity.com/api/v1/team"
const enrollURL = "https://dash.paretosecurity.com/api/v1/team/enroll"

type LinkingDevice struct {
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	OSVersion string `json:"osVersion"`
	Kernel    string `json:"kernel"`
	UUID      string `json:"uuid"`
	Ticket    string `json:"ticket"`
	Version   string `json:"version"`
}

type LinkingResponse struct {
	Team string `json:"team"`
	Auth string `json:"auth"`
}

type TicketResponse struct {
	URL string `json:"url"`
}

func getTransport() http.RoundTripper {
	if testing.Testing() {
		return reqtest.Replay("fixtures")
	}
	return http.DefaultTransport
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
	uuid, err := shared.SystemUUID()
	if err != nil {
		return nil, err
	}
	ticket, err := shared.SystemDevice()
	if err != nil {
		return nil, err
	}
	return &LinkingDevice{
		Hostname:  os.Getenv("HOSTNAME"),
		OS:        sysinfo.OS.Name,
		OSVersion: sysinfo.OS.Release,
		Kernel:    sysinfo.Kernel.Release,
		UUID:      uuid,
		Ticket:    ticket,
	}, nil
}

// LinkAndWaitForTicket initiates the device linking process and waits for the device to be linked to a team.
// It performs the following steps:
// 1. Creates a new linking device.
// 2. Sends a request to enroll the device and retrieves the enrollment URL.
// 3. Logs the enrollment URL for the user to visit.
// 4. Starts a goroutine that periodically checks the link status until the device is linked to a team.
// 5. Updates the configuration with the team ID and authentication token once the device is linked.
// 6. Attempts to open the enrollment URL in the default web browser.
//
// Returns an error if any step in the process fails.
func LinkAndWaitForTicket() error {
	device, err := NewLinkingDevice()
	if err != nil {
		return err
	}

	var linkResp TicketResponse
	err = requests.
		URL(enrollURL).
		BodyJSON(&device).
		ToJSON(&linkResp).
		Transport(getTransport()).
		Fetch(context.Background())
	if err != nil {
		return err
	}

	log.Infof("Please visit the following URL to enroll your device: %s", linkResp.URL)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			var linkStatus LinkingResponse
			err := requests.
				URL(baseURL).
				ToJSON(&linkStatus).
				Transport(getTransport()).
				Fetch(context.Background())
			if err != nil {
				log.Errorf("Error checking link status: %v", err)
				continue
			}

			if linkStatus.Auth != "" {
				shared.Config.TeamID = linkStatus.Team
				shared.Config.AuthToken = linkStatus.Auth
				err := shared.SaveConfig()
				if err != nil {
					log.Errorf("Error saving config: %v", err)
					os.Exit(1)
				}
				log.Infof("Device successfully linked to team: %s", linkStatus.Team)
				break
			} else {
				log.Infof("Waiting for device to be linked to team...")
			}
		}
	}()

	err = exec.Command("open", linkResp.URL).Start()
	if err != nil {
		log.Warnf("Error opening browser: %v", err)
	}

	return nil
}
