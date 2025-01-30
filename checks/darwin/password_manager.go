package checks

import (
	"os"
	"path/filepath"
)

type PasswordManagerCheck struct {
	passed bool
	status string
}

func (pmc *PasswordManagerCheck) Name() string {
	return "Password Manager Presence"
}

func (pmc *PasswordManagerCheck) Run() error {
	paths := []string{
		"/Applications/1Password.app",
		"/Applications/1Password 8.app",
		"/Applications/1Password 7.app",
		"/Applications/Bitwarden.app",
		"/Applications/Dashlane.app",
		"/Applications/KeePassXC.app",
		"/Applications/KeePassX.app",
		"/System/Applications/1Password.app",
		"/System/Applications/1Password 8.app",
		"/System/Applications/1Password 7.app",
		"/System/Applications/Bitwarden.app",
		"/System/Applications/Dashlane.app",
		"/System/Applications/KeePassXC.app",
		"/System/Applications/KeePassX.app",
		filepath.Join(os.Getenv("HOME"), "Applications/1Password.app"),
		filepath.Join(os.Getenv("HOME"), "Applications/1Password 8.app"),
		filepath.Join(os.Getenv("HOME"), "Applications/1Password 7.app"),
		filepath.Join(os.Getenv("HOME"), "Applications/Bitwarden.app"),
		filepath.Join(os.Getenv("HOME"), "Applications/Dashlane.app"),
		filepath.Join(os.Getenv("HOME"), "Applications/KeePassXC.app"),
		filepath.Join(os.Getenv("HOME"), "Applications/KeePassX.app"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			pmc.passed = true
			pmc.status = "Password manager is present"
			return nil
		}
	}

	pmc.passed = false
	pmc.status = "No password manager found"
	return nil
}

func (pmc *PasswordManagerCheck) Passed() bool {
	return pmc.passed
}

func (pmc *PasswordManagerCheck) IsRunnable() bool {
	return true
}

func (pmc *PasswordManagerCheck) UUID() string {
	return "d220f1a2-4c5b-0766-9fb5-6bc9963e6721"
}

func (pmc *PasswordManagerCheck) PassedMessage() string {
	return "Password manager is present"
}

func (pmc *PasswordManagerCheck) FailedMessage() string {
	return "No password manager found"
}

func (pmc *PasswordManagerCheck) RequiresRoot() bool {
	return false
}

func (pmc *PasswordManagerCheck) Status() string {
	if pmc.Passed() {
		return pmc.PassedMessage()
	}
	return pmc.FailedMessage()
}
