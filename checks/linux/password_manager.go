package checks

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ParetoSecurity/pareto-core/shared"
)

type PasswordManagerCheck struct {
	passed bool
}

func (pmc *PasswordManagerCheck) Name() string {
	return "Password Manager Presence"
}

func (pmc *PasswordManagerCheck) Run() error {

	// Check for password managers installed via package managers

	passwordManagers := []string{"1password", "bitwarden", "dashlane", "keepassx", "keepassxc"}

	for _, pwdManager := range passwordManagers {
		if isPackageInstalled(pwdManager) {
			pmc.passed = true
			return nil
		}
	}

	pmc.passed = checkForBrowserExtensions()
	return nil
}

func checkForBrowserExtensions() bool {
	home := os.Getenv("HOME")
	extensionPaths := map[string]string{
		"Google Chrome":  filepath.Join(home, ".config", "google-chrome", "Default", "Extensions"),
		"Firefox":        filepath.Join(home, ".mozilla", "firefox"),
		"Microsoft Edge": filepath.Join(home, ".config", "microsoft-edge", "Default", "Extensions"),
		"Brave Browser":  filepath.Join(home, ".config", "BraveSoftware", "Brave-Browser", "Default", "Extensions"),
	}

	browserExtensions := []string{
		"LastPass",
		"ProtonPass",
		"NordPass",
		"Bitwarden",
		"1Password",
		"KeePass",
		"Dashlane",
	}

	for _, extPath := range extensionPaths {
		if _, err := os.Stat(extPath); err == nil {
			entries, err := os.ReadDir(extPath)
			if err == nil {
				for _, entry := range entries {
					name := strings.ToLower(entry.Name())
					for _, ext := range browserExtensions {
						if strings.Contains(name, strings.ToLower(ext)) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func isPackageInstalled(pkgName string) bool {
	pkgManagers := make(map[string]string)

	// Check which package managers are available
	if _, err := shared.RunCommand("which", "dpkg"); err == nil {
		pkgManagers["apt"] = "dpkg -l"
	}
	if _, err := shared.RunCommand("which", "snap"); err == nil {
		pkgManagers["snap"] = "snap list"
	}
	if _, err := shared.RunCommand("which", "yum"); err == nil {
		pkgManagers["yum"] = "yum list installed"
	}
	if _, err := shared.RunCommand("which", "flatpak"); err == nil {
		pkgManagers["flatpak"] = "flatpak list"
	}
	if _, err := shared.RunCommand("which", "pacman"); err == nil {
		pkgManagers["pacman"] = "pacman -Q"
	}

	for pkgManager, baseCmd := range pkgManagers {
		// Use cache or get fresh data
		cacheKey := "pkg_" + pkgManager
		cached, ok := shared.GetCache(cacheKey)
		if !ok {
			var err error
			cached, err = shared.RunCommand("sh", "-c", baseCmd)
			if err != nil {
				continue
			}
			shared.SetCache(cacheKey, cached, 10) // Cache for 10 seconds
		}

		if strings.Contains(cached, pkgName) {
			return true
		}
	}

	return false
}

func (pmc *PasswordManagerCheck) Passed() bool {
	return pmc.passed
}

func (pmc *PasswordManagerCheck) IsRunnable() bool {
	return true
}

func (pmc *PasswordManagerCheck) UUID() string {
	return "f962c423-fdf5-428a-a57a-827abc9b253e"
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
