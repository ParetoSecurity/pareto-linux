package team

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"github.com/davecgh/go-spew/spew"

	"github.com/ParetoSecurity/pareto-linux/claims"
	"github.com/ParetoSecurity/pareto-linux/shared"
)

const reportURL = "https://dash.paretosecurity.com"

type Report struct {
	PassedCount       int                    `json:"passedCount"`
	FailedCount       int                    `json:"failedCount"`
	DisabledCount     int                    `json:"disabledCount"`
	Device            shared.ReportingDevice `json:"device"`
	Version           string                 `json:"version"`
	LastCheck         string                 `json:"lastCheck"`
	SignificantChange string                 `json:"significantChange"`
	State             map[string]string      `json:"state"`
}

func currentTime() string {
	if testing.Testing() {
		return "2025-01-17T20:29:40+01:00"
	}
	return time.Now().Format(time.RFC3339)
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
		Device:            shared.CurrentReportingDevice(),
		Version:           shared.Version,
		LastCheck:         currentTime(),
		SignificantChange: hex.EncodeToString(significantChange[:]),
		State:             checkStates,
	}
}

// ReportAndSave generates a report and saves it to the configuration file.
func ReportToTeam() error {
	report := NowReport()
	log.Debug(spew.Sdump(report))
	res := ""
	errRes := ""
	err := requests.URL(reportURL).
		Pathf("/api/v1/team/%s/device", shared.Config.TeamID).
		Method(http.MethodPatch).
		Transport(shared.HTTPTransport()).
		BodyJSON(&report).
		ToString(&res).
		AddValidator(
			requests.ValidatorHandler(
				requests.DefaultValidator,
				requests.ToString(&errRes),
			)).
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
