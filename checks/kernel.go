package checks

import (
	"os/exec"
	"strings"
)

type KernelParamsCheck struct {
	passed bool
	status string
}

func (k *KernelParamsCheck) Name() string {
	return "Kernel Parameters are set correctly"
}

func (k *KernelParamsCheck) PassedMessage() string {
	return "All critical kernel parameters are set correctly."
}

func (k *KernelParamsCheck) FailedMessage() string {
	return "Some critical kernel parameters are not set correctly."
}

func (k *KernelParamsCheck) Run() error {
	params := map[string]string{
		"net.ipv4.conf.all.rp_filter": "1",
		"net.ipv4.tcp_syncookies":     "1",
		"kernel.randomize_va_space":   "2",
		"fs.protected_hardlinks":      "1",
		"fs.protected_symlinks":       "1",
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
	return "kernel-params-check"
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
