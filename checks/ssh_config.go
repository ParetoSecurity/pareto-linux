package checks

import (
	"os"
	"os/exec"
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

	s.passed = true

	//run sshd -T to get the sshd config
	configRaw, err := exec.Command("sshd", "-T").CombinedOutput()
	log.WithField("check", s.Name()).Debugf("sshd -T output: %s", configRaw)
	config := strings.ToLower(string(configRaw))
	if err != nil {
		s.passed = false
		s.status = "Failed to get sshd config"
	}

	if strings.Contains(config, "passwordauthentication yes") {
		s.passed = false
		s.status = "PasswordAuthentication is enabled"
	}
	if strings.Contains(config, "permitrootlogin yes") {
		s.passed = false
		s.status = "Root login is enabled"
	}
	if strings.Contains(config, "permitemptypasswords yes") {
		s.passed = false
		s.status = "Empty passwords are allowed"
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
