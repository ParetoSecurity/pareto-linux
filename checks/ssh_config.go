package checks

import (
	"os"
	"strings"

	"github.com/caarlos0/log"
	"paretosecurity.com/auditor/shared"
)

type SSHConfigCheck struct {
	passed bool
	status string
}

func (s *SSHConfigCheck) Name() string {
	return "SSH Server Configuration is Secure"
}

func (s *SSHConfigCheck) PassedMessage() string {
	return "SSH configuration is secure."
}

func (s *SSHConfigCheck) FailedMessage() string {
	return "SSH configuration is not secure."
}

func (s *SSHConfigCheck) Run() error {
	if s.RequiresRoot() && !shared.IsRoot() {
		log.Debug("Running check via helper")
		// Run as root
		passed, err := shared.RunCheckViaHelper(s.UUID())
		if err != nil {
			log.WithError(err).Warn("Failed to run check via helper")
			return err
		}
		s.passed = passed
		return nil
	}
	log.Debug("Running check directly")

	s.passed = false
	data, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		return err
	}
	config := string(data)
	if strings.Contains(config, "PasswordAuthentication no") {
		s.passed = true
		s.status = "PasswordAuthentication is enabled"
	}
	if strings.Contains(config, "PermitRootLogin no") {
		s.passed = true
		s.status = "Root login is enabled"
	}
	return nil
}

func (s *SSHConfigCheck) Passed() bool {
	return s.passed
}

func (s *SSHConfigCheck) IsRunnable() bool {
	if _, err := os.Stat("/etc/ssh/sshd_config"); os.IsNotExist(err) {
		return false
	}
	return true
}

func (s *SSHConfigCheck) ReportIfDisabled() bool {
	return true
}

func (s *SSHConfigCheck) UUID() string {
	return "da4edd80-6af0-4fb3-9fc7-f9a0e9d07f3b"
}

func (s *SSHConfigCheck) Status() string {
	if s.Passed() {
		return s.PassedMessage()
	}
	return s.FailedMessage()
}

func (s *SSHConfigCheck) RequiresRoot() bool {
	return true
}
