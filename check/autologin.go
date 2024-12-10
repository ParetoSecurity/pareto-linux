package check

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Autologin struct {
	passed bool
}

// Name returns the name of the check
func (f *Autologin) Name() string {
	return "Autologin status"
}

// Run executes the check
func (f *Autologin) Run() error {
	f.passed = true

	// Check KDE (SDDM) autologin
	sddmFiles, _ := filepath.Glob("/etc/sddm.conf.d/*.conf")
	for _, file := range sddmFiles {
		content, err := os.ReadFile(file)
		if err == nil {
			if strings.Contains(string(content), "Autologin=true") {
				f.passed = false
				return nil
			}
		}
	}

	// Check main SDDM config
	if content, err := os.ReadFile("/etc/sddm.conf"); err == nil {
		if strings.Contains(string(content), "Autologin=true") {
			f.passed = false
			return nil
		}
	}

	// Check GNOME (GDM) autologin
	gdmPaths := []string{"/etc/gdm3/custom.conf", "/etc/gdm/custom.conf"}
	for _, path := range gdmPaths {
		if content, err := os.ReadFile(path); err == nil {
			if strings.Contains(string(content), "AutomaticLoginEnable=true") {
				f.passed = false
				return nil
			}
		}
	}

	time.Sleep(time.Duration(1 * time.Second))
	return nil
}

// Passed returns the status of the check
func (f *Autologin) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *Autologin) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *Autologin) UUID() string {
	return "f962c423-fdf5-428a-a57a-816abc9b253e"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *Autologin) ReportIfDisabled() bool {
	return false
}

// PassedMessage returns the message to return if the check passed
func (f *Autologin) PassedMessage() string {
	return "Automatic login is off"
}

// FailedMessage returns the message to return if the check failed
func (f *Autologin) FailedMessage() string {
	return "Automatic login is on"
}

// RequiresRoot returns whether the check requires root access
func (f *Autologin) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *Autologin) Status() string {
	if !f.Passed() {
		return f.FailedMessage()
	}
	return f.PassedMessage()
}
