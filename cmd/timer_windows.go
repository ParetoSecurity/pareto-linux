//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/caarlos0/log"
)

func isUserTimerInstalled() bool {
	// Task name
	taskName := "ParetoHourly"

	// Check if the task exists
	command := fmt.Sprintf(`schtasks /Query /TN "%s"`, taskName)

	// Execute the command
	log.Infof("Checking if task exists: %s", taskName)
	log.Infof("Command: %s", command)

	cmd := exec.Command("cmd", "/C", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Infof("Task does not exist: %s", taskName)
		return false
	}

	log.Infof("Task exists: %s", taskName)
	log.Infof(string(output))
	return true
}

func uninstallUserTimer() {
	// Task name
	taskName := "ParetoHourly"

	// Remove the task using schtasks
	command := fmt.Sprintf(`schtasks /Delete /TN "%s" /F`, taskName)

	// Execute the command
	log.Infof("Removing task: %s", taskName)
	log.Infof("Command: %s", command)

	cmd := exec.Command("cmd", "/C", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error removing task: %v\nOutput: %s", err, string(output))
	}

	log.Info("Task removed successfully!")
	log.Info(string(output))
}

func installUserTimer() {
	// Task name
	taskName := "ParetoHourly"

	// Command to execute (Replace this with your Go executable or script path)
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}
	// Create the scheduled task using schtasks
	command := fmt.Sprintf(`schtasks /Create /SC HOURLY /TN "%s" /TR "%s" /F`, taskName, executablePath)

	// Execute the command
	log.Infof("Creating task: %s", taskName)
	log.Infof("Command: %s", command)

	cmd := exec.Command("cmd", "/C", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error creating task: %v\nOutput: %s", err, string(output))
	}

	log.Info("Task created successfully!")
	log.Info(string(output))
}
