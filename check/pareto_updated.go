package check

import (
	"context"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"paretosecurity.com/auditor/shared"
)

type ParetoReleases []struct {
	TagName         string    `json:"tag_name,omitempty"`
	TargetCommitish string    `json:"target_commitish,omitempty"`
	Name            string    `json:"name,omitempty"`
	Draft           bool      `json:"draft,omitempty"`
	Prerelease      bool      `json:"prerelease,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	PublishedAt     time.Time `json:"published_at,omitempty"`
}

type ParetoUpdated struct {
	passed  bool
	details string
}

// Name returns the name of the check
func (f *ParetoUpdated) Name() string {
	return "Pareto Security is up to date"
}

// Run executes the check
func (f *ParetoUpdated) Run() error {
	f.passed = false
	res := ParetoReleases{}
	device := shared.CurrentReportingDevice()

	// uuid=REDACTED&version=1.7.91&os_version=15.1.1&distribution=app-live-setapp"
	err := requests.URL("https://paretosecurity.com/api/updates").
		Param("uuid", device.MachineUUID).
		Param("version", shared.Version).
		Param("os_version", device.LinuxOSVersion).
		Param("platform", "linux").
		Param("app", "auditor").
		Param("distribution", func() string {
			if shared.IsLinked() {
				return "app-live-team"
			}
			return "app-live-opensource"
		}()).
		Transport(shared.HTTPTransport()).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {

		log.WithError(err).
			Warnf("Failed to report to team: %s", shared.Config.TeamID)
		return err
	}

	if len(res) == 0 {
		f.details = "No releases found"
	}

	if res[0].TagName == shared.Version {
		f.passed = true
	}

	return nil
}

// Passed returns the status of the check
func (f *ParetoUpdated) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *ParetoUpdated) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *ParetoUpdated) UUID() string {
	return "05a103af-6031-42b7-8ff4-655f9a27ddf2"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *ParetoUpdated) ReportIfDisabled() bool {
	return false
}

// PassedMessage returns the message to return if the check passed
func (f *ParetoUpdated) PassedMessage() string {
	return "Pareto Security is up to date"
}

// FailedMessage returns the message to return if the check failed
func (f *ParetoUpdated) FailedMessage() string {
	return "Pareto Security is oudated"
}

// RequiresRoot returns whether the check requires root access
func (f *ParetoUpdated) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *ParetoUpdated) Status() string {
	return f.details
}
