package check

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/log"
)

type RemoteLogin struct {
	passed bool
	ports  map[int]string
}

// Name returns the name of the check
func (f *RemoteLogin) Name() string {
	return "Remote login is disabled"
}

// checkPort tests if a port is open
func (f *RemoteLogin) checkPort(port int) bool {
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
		log.WithField("address", address).Debug("Checking port")
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			conn.Close()
			return true
		}
	}

	return false
}

// Run executes the check
func (f *RemoteLogin) Run() error {
	f.passed = true
	f.ports = make(map[int]string)

	// Check common remote access ports
	portsToCheck := map[int]string{
		22:   "SSH",
		3389: "RDP",
		3390: "RDP",
		5900: "VNC",
	}

	for port, service := range portsToCheck {
		if f.checkPort(port) {
			f.passed = false
			f.ports[port] = service
		}
	}

	return nil
}

// Passed returns the status of the check
func (f *RemoteLogin) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *RemoteLogin) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *RemoteLogin) UUID() string {
	return "4ced961d-7cfc-4e7b-8f80-195f6379446e"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *RemoteLogin) ReportIfDisabled() bool {
	return false
}

// PassedMessage returns the message to return if the check passed
func (f *RemoteLogin) PassedMessage() string {
	return "Remote access services are found running"
}

// FailedMessage returns the message to return if the check failed
func (f *RemoteLogin) FailedMessage() string {
	return "No remote access services found running"
}

// RequiresRoot returns whether the check requires root access
func (f *RemoteLogin) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *RemoteLogin) Status() string {
	if !f.Passed() {
		msg := "Remote access services found running on ports:"
		for port, service := range f.ports {
			msg += fmt.Sprintf(" %s(%d)", service, port)
		}
		return msg
	}
	return "No remote access services found running"
}
