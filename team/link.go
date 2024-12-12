package team

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"github.com/davecgh/go-spew/spew"
	"paretosecurity.com/auditor/shared"
)

const enrollURL = "https://dash.paretosecurity.com/api/v1/team/enroll"

type NewDevice struct {
	MachineName    string `json:"machineName"`
	ModelName      string `json:"modelName"`
	ModelSerial    string `json:"modelSerial"`
	LinuxOSVersion string `json:"LinuxOSVersion"`
	MachineUUID    string `json:"machineUUID"`
	Auth           string `json:"auth"`
}

type LinkingResponse struct {
	Team string `json:"team"`
	Auth string `json:"auth"`
}

type TicketResponse struct {
	URL string `json:"url"`
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
	var wg sync.WaitGroup
	device, err := shared.NewLinkingDevice()
	if err != nil {
		return err
	}
	log.Infof("Linking device with UUID: %s", device.UUID)
	log.Infof("Device hostname: %s", device.Hostname)
	log.Infof("Device OS: %s", device.OS)
	log.Infof("Device OS version: %s", device.OSVersion)
	log.Infof("Device kernel version: %s", device.Kernel)
	log.Infof("Device ticket: %s", device.Ticket)
	log.Info("Please wait while we link your device to a team...")
	var linkResp TicketResponse
	err = requests.
		URL(enrollURL).
		BodyJSON(&device).
		ToJSON(&linkResp).
		Transport(shared.HTTPTransport()).
		Fetch(context.Background())
	if err != nil {
		return err
	}

	log.Infof("Please visit the following URL to enroll your device: %s", linkResp.URL)
	wg.Add(1)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			var linkStatus LinkingResponse
			err := requests.
				URL(enrollURL).
				Param("ticket", device.Ticket).
				BodyJSON(&device).
				ToJSON(&linkStatus).
				Transport(shared.HTTPTransport()).
				Fetch(context.Background())
			if err != nil {
				log.Errorf("Error checking link status: %v", err)
				var httpErr *requests.ResponseError
				if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {

					fmt.Print(".")
					continue
				}
				continue
			}

			if linkStatus.Auth != "" {
				shared.Config.TeamID = linkStatus.Team
				shared.Config.AuthToken = linkStatus.Auth
				err := AddDevice()
				if err != nil {
					log.WithError(err).Fatal("Failed to report device to team")
					os.Exit(1)
				}

				err = shared.SaveConfig()
				if err != nil {
					log.Errorf("Error saving config: %v", err)
					os.Exit(1)
				}

				log.Infof("Device successfully linked to team: %s", linkStatus.Team)
				wg.Done()
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
	wg.Wait()
	return nil
}

// AddDevice reports the device to the team.
func AddDevice() error {
	device := shared.CurrentReportingDevice()

	res := ""
	errRes := ""
	newDevice := NewDevice{
		MachineName:    device.MachineName,
		ModelName:      device.ModelName,
		ModelSerial:    device.ModelSerial,
		LinuxOSVersion: device.LinuxOSVersion,
		MachineUUID:    device.MachineUUID,
		Auth:           device.Auth,
	}
	log.Debug(spew.Sdump(newDevice))
	err := requests.URL(reportURL).
		Pathf("/api/v1/team/%s/device", shared.Config.TeamID).
		Method(http.MethodPut).
		Transport(shared.HTTPTransport()).
		AddValidator(
			requests.ValidatorHandler(
				requests.DefaultValidator,
				requests.ToString(&errRes),
			)).
		BodyJSON(&newDevice).
		ToString(&res).
		Fetch(context.Background())
	if err != nil {

		log.WithField("response", errRes).
			WithError(err).
			Warnf("Failed to report to team: %s", shared.Config.TeamID)
		return err
	}
	log.WithField("response", res).Debug("API Response")
	return nil
}
