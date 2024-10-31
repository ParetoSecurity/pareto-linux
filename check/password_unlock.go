package check

import (
	"os/exec"
	"strings"
)

type PasswordToUnlock struct {
	passed bool
}

// Name returns the name of the check
func (f *PasswordToUnlock) Name() string {
	return "Password is required to unlock the screen"
}

func (f *PasswordToUnlock) checkGnome() bool {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.screensaver", "lock-enabled")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func (f *PasswordToUnlock) checkKDE() bool {
	cmd := exec.Command("kreadconfig5", "--file", "kscreenlockerrc", "--group", "Daemon", "--key", "Autolock")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// Run executes the check
func (f *PasswordToUnlock) Run() error {
	// Check if running GNOME
	if _, err := exec.LookPath("gsettings"); err == nil {
		f.passed = f.checkGnome()
		return nil
	}

	// Check if running KDE
	if _, err := exec.LookPath("kreadconfig5"); err == nil {
		f.passed = f.checkKDE()
		return nil
	}

	// Neither GNOME nor KDE found
	f.passed = false
	return nil
}

// Passed returns the status of the check
func (f *PasswordToUnlock) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *PasswordToUnlock) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *PasswordToUnlock) UUID() string {
	return "37dee029-605b-4aab-96b9-5438e5aa44d8"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *PasswordToUnlock) ReportIfDisabled() bool {
	return false
}

// Status returns the status of the check
func (f *PasswordToUnlock) Status() string {
	if f.Passed() {
		return "Password after sleep or screensaver is on"
	}
	return "Password after sleep or screensaver is off"
}
