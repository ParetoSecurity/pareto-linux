package check

import (
	"os/exec"
	"strings"
)

type SoftwareUpdates struct {
	passed  bool
	details string
}

// Name returns the name of the check
func (f *SoftwareUpdates) Name() string {
	return "Apps are up to date"
}

func (f *SoftwareUpdates) checkUpdates() (bool, string) {
	updates := []string{}

	// Check flatpak
	if _, err := exec.LookPath("flatpak"); err == nil {
		cmd := exec.Command("flatpak", "remote-ls", "--updates")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			updates = append(updates, "Flatpak")
		}
	}

	// Check apt-get
	if _, err := exec.LookPath("apt-get"); err == nil {
		cmd := exec.Command("apt-get", "-s", "upgrade")
		output, err := cmd.Output()
		if err == nil && !strings.Contains(string(output), "0 upgraded, 0 newly installed") {
			updates = append(updates, "APT")
		}
	}

	// Check dnf
	if _, err := exec.LookPath("dnf"); err == nil {
		cmd := exec.Command("dnf", "check-update", "--quiet")
		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 100 {
				updates = append(updates, "DNF")
			}
		}
	}

	// Check pacman
	if _, err := exec.LookPath("pacman"); err == nil {
		cmd := exec.Command("pacman", "-Qu")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			updates = append(updates, "Pacman")
		}
	}

	if len(updates) == 0 {
		return true, "All packages are up to date"
	}
	return false, "Updates available for: " + strings.Join(updates, ", ")
}

// Run executes the check
func (f *SoftwareUpdates) Run() error {
	var ok bool
	ok, f.details = f.checkUpdates()
	f.passed = ok
	return nil
}

// Passed returns the status of the check
func (f *SoftwareUpdates) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *SoftwareUpdates) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *SoftwareUpdates) UUID() string {
	return "940e7a88-2dd4-4a50-bf9c-3d842e0a2c94"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *SoftwareUpdates) ReportIfDisabled() bool {
	return false
}

// Status returns the status of the check
func (f *SoftwareUpdates) Status() string {
	return f.details
}
