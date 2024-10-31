package team

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"paretosecurity.com/auditor/claims"
	"paretosecurity.com/auditor/shared"
)

type ReportingDevice struct {
	MachineUUID  string `json:"machineUUID"`
	MachineName  string `json:"machineName"`
	Auth         string `json:"auth"`
	MacOSVersion string `json:"macOSVersion"`
	ModelName    string `json:"modelName"`
	ModelSerial  string `json:"modelSerial"`
}

func CurrentReportingDevice() ReportingDevice {

	return ReportingDevice{
		MachineUUID: func() string {
			uuid, _ := shared.SystemUUID()
			return uuid
		}(),
		MachineName: func() string {
			name, _ := os.Hostname()
			return name
		}(),
		Auth:         shared.Config.AuthToken,
		MacOSVersion: "Linux",
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
