package checks

import (
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
	packageManagers := []string{"apt", "snap", "yum", "flatpak"}
	passwordManagers := []string{"1password", "bitwarden", "dashlane", "keepassx", "keepassxc"}

	for _, pkgManager := range packageManagers {
		for _, pwdManager := range passwordManagers {
			if isPackageInstalled(pkgManager, pwdManager) {
				pmc.passed = true
				return nil
			}
		}
	}

	pmc.passed = false
	return nil
}

func isPackageInstalled(pkgManager, pkgName string) bool {
	var cmd string
	switch pkgManager {
	case "apt":
		cmd = "dpkg -l | grep " + pkgName
	case "snap":
		cmd = "snap list | grep " + pkgName
	case "yum":
		cmd = "yum list installed | grep " + pkgName
	case "flatpak":
		cmd = "flatpak list | grep " + pkgName
	default:
		return false
	}

	output, err := shared.RunCommand("sh", "-c", cmd)
	return err == nil && strings.Contains(output, pkgName)
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
