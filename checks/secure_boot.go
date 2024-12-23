package checks

import (
	"os"
	"path/filepath"
)

type SecureBoot struct {
	passed bool
	status string
}

// Name returns the name of the check
func (f *SecureBoot) Name() string {
	return "SecureBoot is enabled"
}

// Run executes the check
func (f *SecureBoot) Run() error {

	// Find and read the SecureBoot EFI variable
	pattern := "/sys/firmware/efi/efivars/SecureBoot-*"
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		f.passed = false
		f.status = "Could not find SecureBoot EFI variable"
		return nil
	}

	data, err := os.ReadFile(matches[0])
	if err != nil {
		f.passed = false
		f.status = "Could not read SecureBoot status"
		return nil
	}

	// The SecureBoot variable has a 5-byte structure
	// First 4 bytes are the attribute flags, last byte is the value
	// Value of 1 means enabled, 0 means disabled
	if len(data) >= 5 && data[4] == 1 {
		f.passed = true
		f.status = f.PassedMessage()
	}
	f.passed = false
	f.status = f.FailedMessage()

	return nil
}

// Passed returns the status of the check
func (f *SecureBoot) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *SecureBoot) IsRunnable() bool {
	if _, err := os.Stat("/sys/firmware/efi/efivars"); os.IsNotExist(err) {
		f.status = "System is not running in UEFI mode"
		return true
	}
	return false
}

// UUID returns the UUID of the check
func (f *SecureBoot) UUID() string {
	return "c96524f2-850b-4bb9-abc7-517051b6c14e"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *SecureBoot) ReportIfDisabled() bool {
	return true
}

// PassedMessage returns the message to return if the check passed
func (f *SecureBoot) PassedMessage() string {
	return "SecureBoot is enabled"
}

// FailedMessage returns the message to return if the check failed
func (f *SecureBoot) FailedMessage() string {
	return "SecureBoot is disabled"
}

// RequiresRoot returns whether the check requires root access
func (f *SecureBoot) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *SecureBoot) Status() string {
	return f.status
}
