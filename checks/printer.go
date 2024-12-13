package checks

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/log"
)

type Printer struct {
	passed bool
	ports  map[int]string
}

// Name returns the name of the check
func (f *Printer) Name() string {
	return "Sharing printers is off"
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
		// Check all interfaces
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			// Filter out 127.0.0.1
			if ip.IsLoopback() {
				continue
			}

			address := fmt.Sprintf("%s:%d", ip.String(), port)
			log.WithField("check", f.Name()).WithField("address", address).Debug("Checking port")
			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err == nil {
				f.passed = false
				f.ports[port] = service
				conn.Close()
			}
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
	return "Sharing printers is off"
}

// FailedMessage returns the message to return if the check failed
func (f *Printer) FailedMessage() string {
	return "Sharing printers is on"
}

// RequiresRoot returns whether the check requires root access
func (f *Printer) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *Printer) Status() string {
	if !f.Passed() {
		msg := "Printer sharing services found running on ports:"
		for port, service := range f.ports {
			msg += fmt.Sprintf(" %s(%d)", service, port)
		}
		return msg
	}
	return f.PassedMessage()
}
