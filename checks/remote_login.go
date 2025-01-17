package checks

import (
	"fmt"
	"net"
	"testing"
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
func (f *RemoteLogin) checkPort(port int, proto string) bool {

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
		if f.checkPort(port, "tcp") {
			log.WithField("check", f.Name()).WithField("port", port).WithField("service", service).Debug("Remote access service found")
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

// PassedMessage returns the message to return if the check passed
func (f *RemoteLogin) PassedMessage() string {
	return "No remote access services found running"
}

// FailedMessage returns the message to return if the check failed
func (f *RemoteLogin) FailedMessage() string {
	return "Remote access services found running"
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
	return f.PassedMessage()
}
