package checks

import (
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"paretosecurity.com/auditor/shared"
)

type KernelParamsCheck struct {
	passed bool
	status string
}

func (k *KernelParamsCheck) Name() string {
	return "Kernel Parameters are set correctly"
}

func (k *KernelParamsCheck) PassedMessage() string {
	return "Critical kernel parameters are correct"
}

func (k *KernelParamsCheck) FailedMessage() string {
	return "Critical kernel parameters are not correct"
}

func (k *KernelParamsCheck) Run() error {

	if k.RequiresRoot() && !shared.IsRoot() {
		log.Debug("Running check via helper")
		// Run as root
		passed, err := shared.RunCheckViaHelper(k.UUID())
		if err != nil {
			log.WithError(err).Warn("Failed to run check via helper")
			return err
		}
		k.passed = passed
		return nil
	}

	params := map[string]string{
		"net.ipv4.tcp_syncookies":   "1",
		"kernel.randomize_va_space": "2",
		"fs.protected_hardlinks":    "1",
		"fs.protected_symlinks":     "1",
	}
	k.passed = true
	for param, expected := range params {
		value, err := getSysctlValue(param)
		if err != nil {
			return err
		}
		if value != expected {
			k.passed = false
			k.status += param + " is set to " + value + " but should be " + expected + ". "
		}
	}
	return nil
}

func (k *KernelParamsCheck) Passed() bool {
	return k.passed
}

func (k *KernelParamsCheck) IsRunnable() bool {
	return true
}

func (k *KernelParamsCheck) ReportIfDisabled() bool {
	return true
}

func (k *KernelParamsCheck) UUID() string {
	return "cbf2736b-72df-43e3-8789-8eb676ff9014"
}

func (k *KernelParamsCheck) Status() string {
	if k.Passed() {
		return k.PassedMessage()
	}
	return k.status
}

func (k *KernelParamsCheck) RequiresRoot() bool {
	return true
}

func getSysctlValue(param string) (string, error) {
	out, err := exec.Command("sysctl", "-n", param).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
