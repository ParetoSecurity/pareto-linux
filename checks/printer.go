package checks

import (
	"fmt"
	"net"
	"testing"
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

// checkPort tests if a port is open
func (f *Printer) checkPort(port int, proto string) bool {

	if testing.Testing() {
		return checkPortMock(port, proto)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
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
		conn, err := net.DialTimeout(proto, address, 1*time.Second)
		if err == nil {
			defer conn.Close()
			log.WithField("check", f.Name()).WithField("address", address).WithField("state", true).Debug("Checking port")
			return true
		}
	}

	return false
}

// Run executes the check
func (f *Printer) Run() error {
	f.passed = true
	f.ports = make(map[int]string)

	// Samba, NFS and CUPS ports to check
	printService := map[int]string{
		631: "CUPS",
	}

	for port, service := range printService {
		if f.checkPort(port, "tcp") {
			log.WithField("check", f.Name()).WithField("port", port).WithField("service", service).Debug("Port open")
			f.passed = false
			f.ports[port] = service
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
