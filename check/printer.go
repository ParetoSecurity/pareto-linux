package check

import (
	"fmt"
	"net"
	"time"
)

type Printer struct {
	passed bool
	ports  map[int]string
}

// Name returns the name of the check
func (f *Printer) Name() string {
	return "Printer sharing is disabled"
}

// Run executes the check
func (f *Printer) Run() error {
	f.passed = true
	f.ports = make(map[int]string)

	// Samba, NFS and CUPS ports to check
	shareServices := map[int]string{
		631: "CUPS",
	}

	for port, service := range shareServices {
		// Check localhost first
		addr := fmt.Sprintf("0.0.0.0:%d", port)
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err == nil {
			f.passed = false
			f.ports[port] = service
			conn.Close()
		}

	}

	return nil
}

// Passed returns the status of the check
func (f *Printer) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *Printer) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *Printer) UUID() string {
	return "b96524e0-150b-4bb8-abc7-517051b6c14e"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *Printer) ReportIfDisabled() bool {
	return false
}

// PassedMessage returns the message to return if the check passed
func (f *Printer) PassedMessage() string {
	return "Printer sharing is disabled"
}

// FailedMessage returns the message to return if the check failed
func (f *Printer) FailedMessage() string {
	return "Printer sharing is enabled"
}

// RequiresRoot returns whether the check requires root access
func (f *Printer) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *Printer) Status() string {
	if !f.Passed() {
		msg := "Printer/sharing services found running on ports:"
		for port, service := range f.ports {
			msg += fmt.Sprintf(" %s(%d)", service, port)
		}
		return msg
	}
	return "No printer or file sharing services found running"
}
