package check

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

type Firewall struct {
	passed bool
}

// Name returns the name of the check
func (f *Firewall) Name() string {
	return "Firewall status"
}

func (f *Firewall) isUbuntu() bool {
	if _, err := os.Stat("/etc/lsb-release"); err == nil {
		return true
	}
	return false
}

func (f *Firewall) isFedora() bool {
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		return true
	}
	return false
}

func (f *Firewall) checkUFW() bool {
	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "Status: active")
}

func (f *Firewall) checkFirewalld() bool {
	cmd := exec.Command("systemctl", "is-active", "firewalld")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "active"
}

// Run executes the check
func (f *Firewall) Run() error {
	switch {
	case f.isUbuntu():
		f.passed = f.checkUFW()
	case f.isFedora():
		f.passed = f.checkFirewalld()
	default:
		f.passed = false
	}
	time.Sleep(time.Duration(1 * time.Second))
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

// Status returns the status of the check
func (f *Firewall) Status() string {
	if f.Passed() {
		return f.PassedMessage()
	}
	return f.FailedMessage()
}
