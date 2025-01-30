package checks

import (
	"strings"

	"github.com/ParetoSecurity/pareto-core/shared"
)

type PasswordManagerCheck struct {
	passed bool
	status string
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
				pmc.status = "Password manager is present"
				return nil
			}
		}
	}

	pmc.passed = false
	pmc.status = "No password manager found"
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
