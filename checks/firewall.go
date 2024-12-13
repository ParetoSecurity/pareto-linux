package checks

import (
	"os"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"paretosecurity.com/auditor/shared"
)

type Firewall struct {
	passed bool
}

// Name returns the name of the check
func (f *Firewall) Name() string {
	return "Firewall is on"
}

func (f *Firewall) isUbuntu() bool {
	if _, err := os.Stat("/etc/lsb-release"); err == nil {
		log.WithError(err).Debug("Failed to check for Ubuntu")
		return true
	}
	return false
}

func (f *Firewall) isFedora() bool {
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		log.WithError(err).Debug("Failed to check for Fedora")
		return true
	}
	return false
}

func (f *Firewall) checkUFW() bool {
	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()
	if err != nil {
		log.WithError(err).Warn("Failed to check UFW status")
		return false
	}
	log.WithField("output", string(output)).Debug("UFW status")
	return strings.Contains(string(output), "active")
}

func (f *Firewall) checkFirewalld() bool {
	cmd := exec.Command("systemctl", "is-active", "firewalld")
	output, err := cmd.Output()
	if err != nil {
		log.WithError(err).Warn("Failed to check firewalld status")
		return false
	}
	log.WithField("output", string(output)).Debug("Firewalld status")
	return strings.TrimSpace(string(output)) == "active"
}

// Run executes the check
func (f *Firewall) Run() error {
	if f.RequiresRoot() && !shared.IsRoot() {
		log.Debug("Running check via helper")
		// Run as root
		passed, err := shared.RunCheckViaHelper(f.UUID())
		if err != nil {
			log.WithError(err).Warn("Failed to run check via helper")
			return err
		}
		f.passed = passed
		return nil
	}

	log.Debug("Running check directly")
	switch {
	case f.isUbuntu():
		f.passed = f.checkUFW()
	case f.isFedora():
		f.passed = f.checkFirewalld()
	default:
		f.passed = false
	}
	return nil
}

// Passed returns the status of the check
func (f *Firewall) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *Firewall) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *Firewall) UUID() string {
	return "2e46c89a-5461-4865-a92e-3b799c12034a"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *Firewall) ReportIfDisabled() bool {
	return false
}

// PassedMessage returns the message to return if the check passed
func (f *Firewall) PassedMessage() string {
	return "Firewall is on"
}

// FailedMessage returns the message to return if the check failed
func (f *Firewall) FailedMessage() string {
	return "Firewall is off"
}

// RequiresRoot returns whether the check requires root access
func (f *Firewall) RequiresRoot() bool {
	return true
}

// Status returns the status of the check
func (f *Firewall) Status() string {
	if f.Passed() {
		return f.PassedMessage()
	}
	return f.FailedMessage()
}
