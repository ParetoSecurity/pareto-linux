package checks

import (
	"strings"

	"github.com/caarlos0/log"

	"github.com/ParetoSecurity/pareto-linux/shared"
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
	configRaw, err := shared.RunCommand("sshd", "-T")
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

	// Check if sshd service is running via systemd
	sshdStatus, err := shared.RunCommand("systemctl", "is-active", "sshd")
	if err != nil || strings.TrimSpace(string(sshdStatus)) == "active" {
		return true
	}

	// Check if ssh service is running via systemd
	sshStatus, err := shared.RunCommand("systemctl", "is-active", "ssh")
	if err != nil || strings.TrimSpace(string(sshStatus)) == "active" {
		return true
	}
	// Check if ssh socket service is enabled via systemd
	sshSocketStatus, err := shared.RunCommand("systemctl", "is-enabled", "sshd.socket")
	if err != nil || strings.TrimSpace(string(sshSocketStatus)) == "enabled" {
		return true
	}

	// Check if ssh socket service is enabled via systemd
	sshSocketStatus, err = shared.RunCommand("systemctl", "is-enabled", "ssh.socket")
	if err != nil || strings.TrimSpace(string(sshSocketStatus)) == "enabled" {
		return true
	}

	return false
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
	return s.status
}

func (s *SSHConfigCheck) RequiresRoot() bool {
	return true
}
