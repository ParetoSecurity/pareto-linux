package check

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

type EncryptingFS struct {
	passed bool
	status string
}

// Name returns the name of the check
func (f *EncryptingFS) Name() string {
	return "Block device encryption"
}

// Passed returns the status of the check
func (f *EncryptingFS) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *EncryptingFS) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *EncryptingFS) UUID() string {
	return "c3aee29a-f16d-4573-a861-b3ba0d860067"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *EncryptingFS) ReportIfDisabled() bool {
	return true
}

// PassedMessage returns the message to return if the check passed
func (f *EncryptingFS) PassedMessage() string {
	return "Block device encryption is enabled"
}

// FailedMessage returns the message to return if the check failed
func (f *EncryptingFS) FailedMessage() string {
	return "Block device encryption is disabled"
}

// Status returns the status of the check
func (f *EncryptingFS) Status() string {
	return f.status
}

// Run executes the check
func (f *EncryptingFS) Run() error {
	// Check if cryptsetup is available
	if _, err := exec.LookPath("cryptsetup"); err != nil {
		f.passed = false
		f.status = "cryptsetup not found"
		return nil
	}

	encryptedDevices := make(map[string]string)
	mountPoints := make(map[string]string)

	// Read crypttab to get encrypted devices
	crypttab, err := os.Open("/etc/crypttab")
	if err == nil {
		scanner := bufio.NewScanner(crypttab)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				encryptedDevices[fields[0]] = fields[1]
			}
		}
		crypttab.Close()
	}

	// Read fstab to get mount points
	fstab, err := os.Open("/etc/fstab")
	if err == nil {
		scanner := bufio.NewScanner(fstab)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if strings.HasPrefix(fields[0], "/dev/mapper/") {
					mountPoints[fields[1]] = fields[0]
				}
			}
		}
		fstab.Close()
	}

	rootEncrypted := false
	homeEncrypted := false

	// Check if root and home are encrypted
	for mp := range mountPoints {
		if mp == "/" {
			rootEncrypted = true
		}
		if mp == "/home" {
			homeEncrypted = true
		}
	}

	if rootEncrypted && homeEncrypted {
		f.passed = true
		f.status = "Both root and home are LUKS encrypted"
	} else if rootEncrypted {
		f.passed = true
		f.status = "Only root is LUKS encrypted"
	} else if homeEncrypted {
		f.passed = true
		f.status = "Only home is LUKS encrypted"
	} else {
		f.passed = false
		f.status = "Neither root nor home are LUKS encrypted"
	}

	return nil
}
