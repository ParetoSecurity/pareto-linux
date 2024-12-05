package team

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"github.com/davecgh/go-spew/spew"
	"paretosecurity.com/auditor/claims"
	"paretosecurity.com/auditor/shared"
)

const reportURL = "https://dash.paretosecurity.com/api/v1/team"

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
		panic(err)
	}
	return ReportingDevice{
		MachineUUID:  device.UUID,
		MachineName:  device.Hostname,
		Auth:         DeviceAuth(),
		MacOSVersion: fmt.Sprintf("%s %s", device.OS, device.OSVersion),
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
func ReportToTeam() {
	report := NowReport()
	log.Debug(spew.Sdump(report))
	err := requests.URL(reportURL).
		Pathf("/%s/device", shared.Config.TeamID).
		Transport(shared.HTTPTransport()).
		BodyJSON(&report).
		Fetch(context.Background())
	if err != nil {
		log.WithError(err).Warnf("Failed to report to team: %s", shared.Config.TeamID)
	}
}
