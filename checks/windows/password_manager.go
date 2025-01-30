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
	// TODO; need real paths
	paths := []string{
		filepath.Join(os.Getenv("PROGRAMFILES"), "1Password", "1Password.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES"), "Bitwarden", "Bitwarden.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES"), "Dashlane", "Dashlane.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES"), "KeePassX", "KeePassX.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES"), "KeePassXC", "KeePassXC.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "1Password", "1Password.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Bitwarden", "Bitwarden.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Dashlane", "Dashlane.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "KeePassX", "KeePassX.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "KeePassXC", "KeePassXC.exe"),
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
