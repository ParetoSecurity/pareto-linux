//go:build darwin
// +build darwin

package cmd

import (
	"github.com/caarlos0/log"
)

func isUserTimerInstalled() bool {
	return false
}

func uninstallUserTimer() {
	log.Info("Removing user timer is not supported on macOS")
}

func installUserTimer() {
	log.Info("Installing user timer is not supported on macOS")
}
