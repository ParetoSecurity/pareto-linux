//go:build linux
// +build linux

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caarlos0/log"
)

const timerContent = `[Unit]
Description=Timer for pareto-core hourly execution

[Timer]
OnCalendar=hourly
Persistent=true

[Install]
WantedBy=timers.target`

const localServiceContent = `[Unit]
Description=Service for pareto-core

[Service]
Type=oneshot
ExecStart=/usr/bin/paretosecurity check
StandardInput=null

[Install]
WantedBy=timers.target`

func isUserTimerInstalled() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory:", err)
		return false
	}

	systemdPath := filepath.Join(homeDir, ".config", "systemd", "user")
	if _, err := os.Stat(filepath.Join(systemdPath, "pareto-coretimer")); err == nil {
		return true
	}
	return false
}

func installUserTimer() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory:", err)
		return
	}

	systemdPath := filepath.Join(homeDir, ".config", "systemd", "user")
	if err := os.MkdirAll(systemdPath, 0755); err != nil {
		log.Fatalf("Failed to create systemd user directory:", err)
		return
	}

	// Create timer file
	timerPath := filepath.Join(systemdPath, "pareto-coretimer")
	if err := os.WriteFile(timerPath, []byte(timerContent), 0644); err != nil {
		log.Fatalf("Failed to create timer file:", err)
		return
	}

	// Create service file
	servicePath := filepath.Join(systemdPath, "pareto-coreservice")
	if err := os.WriteFile(servicePath, []byte(localServiceContent), 0644); err != nil {
		log.Fatalf("Failed to create service file:", err)
		return
	}

	// Execute commands
	if err := exec.Command("systemctl", "--user", "daemon-reload").Run(); err != nil {
		log.Fatalf("Failed to reload systemd:", err)
		return
	}
	if err := exec.Command("systemctl", "--user", "enable", "--now", "pareto-coretimer").Run(); err != nil {
		log.Fatalf("Failed to enable and start timer:", err)
		return
	}

	log.Info("Timer installed successfully, to enable it run:")
	log.Infof("sudo loginctl enable-linger %s", os.Getenv("USER"))
}

func uninstallUserTimer() {
	// Logic to uninstall the user timer
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory:", err)
		return
	}

	systemdPath := filepath.Join(homeDir, ".config", "systemd", "user")
	if err := os.MkdirAll(systemdPath, 0755); err != nil {
		log.Fatalf("Failed to create systemd user directory:", err)
		return
	}

	// Create timer file
	timerPath := filepath.Join(systemdPath, "pareto-coretimer")
	if err := os.Remove(timerPath); err != nil {
		log.Fatalf("Failed to remove timer file:", err)
		return
	}

	// Create service file
	servicePath := filepath.Join(systemdPath, "pareto-coreservice")
	if err := os.Remove(servicePath); err != nil {
		log.Fatalf("Failed to remove service file:", err)
		return
	}
	// Execute commands
	if err := exec.Command("systemctl", "--user", "daemon-reload").Run(); err != nil {
		log.Fatalf("Failed to reload systemd:", err)
		return
	}
}
