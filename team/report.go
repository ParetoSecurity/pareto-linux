package team

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"github.com/davecgh/go-spew/spew"
	"paretosecurity.com/auditor/claims"
	"paretosecurity.com/auditor/shared"
)

const reportURL = "https://dash.paretosecurity.com"

type ReportingDevice struct {
	MachineUUID  string `json:"machineUUID"`
	MachineName  string `json:"machineName"`
	Auth         string `json:"auth"`
	MacOSVersion string `json:"macOSVersion"`
	ModelName    string `json:"modelName"`
	ModelSerial  string `json:"modelSerial"`
}

func CurrentReportingDevice() ReportingDevice {
	device, err := NewLinkingDevice()
	if err != nil {
		log.WithError(err).Fatal("Failed to get device information")
	}

	return ReportingDevice{
		MachineUUID: device.UUID,
		MachineName: device.Hostname,
		Auth:        DeviceAuth(),
		MacOSVersion: func() string {
			if runtime.GOOS == "darwin" {
				version, err := shared.MacOSVersion()
				if err == nil {
					return version
				}
			}
			return fmt.Sprintf("%s %s", device.OS, device.OSVersion)
		}(),
		ModelName: func() string {
			modelName, err := shared.SystemDevice()
			if err != nil {
				return "Unknown"
			}

			return modelName
		}(),
		ModelSerial: func() string {
			serial, err := shared.SystemSerial()
			if err != nil {
				return "Unknown"
			}

			return serial
		}(),
	}
}

type Report struct {
	PassedCount       int               `json:"passedCount"`
	FailedCount       int               `json:"failedCount"`
	DisabledCount     int               `json:"disabledCount"`
	Device            ReportingDevice   `json:"device"`
	Version           string            `json:"version"`
	LastCheck         string            `json:"lastCheck"`
	SignificantChange string            `json:"significantChange"`
	State             map[string]string `json:"state"`
}

func NowReport() Report {
	passed := 0
	failed := 0
	disabled := 0
	disabledSeed, _ := shared.SystemUUID()
	failedSeed, _ := shared.SystemUUID()
	checkStates := make(map[string]string)

	for _, claim := range claims.All {
		if claim.Title != "My Checks" {
			for _, check := range claim.Checks {
				if check.IsRunnable() {
					if check.Passed() {
						passed++
						checkStates[check.UUID()] = "pass"
					} else {
						failed++
						failedSeed += check.UUID()
						checkStates[check.UUID()] = "fail"
					}
				} else {
					if check.ReportIfDisabled() {
						disabled++
						disabledSeed += check.UUID()
						checkStates[check.UUID()] = "off"
					}
				}
			}
		}
	}

	significantChange := sha256.Sum256([]byte(disabledSeed + "." + failedSeed))
	return Report{
		PassedCount:       passed,
		FailedCount:       failed,
		DisabledCount:     disabled,
		Device:            CurrentReportingDevice(),
		Version:           shared.Version,
		LastCheck:         time.Now().Format(time.RFC3339),
		SignificantChange: hex.EncodeToString(significantChange[:]),
		State:             checkStates,
	}
}

// ReportAndSave generates a report and saves it to the configuration file.
func ReportToTeam() error {
	report := NowReport()
	log.Debug(spew.Sdump(report))
	res := ""
	err := requests.URL(reportURL).
		Pathf("/api/v1/team/%s/device", shared.Config.TeamID).
		Method(http.MethodPatch).
		Transport(shared.HTTPTransport()).
		BodyJSON(&report).
		ToString(&res).
		Fetch(context.Background())
	if err != nil {

		log.WithField("response", res).
			WithError(err).
			Warnf("Failed to report to team: %s", shared.Config.TeamID)
		return err
	}
	log.WithField("response", res).Debug("API Response")
	return nil
}
