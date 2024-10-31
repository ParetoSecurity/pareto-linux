package check

import (
	"fmt"
	"net"
	"time"
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
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Run executes the check
func (f *RemoteLogin) Run() error {
	f.passed = true
	f.ports = make(map[int]string)

	// Check common remote access ports
	portsToCheck := map[int]string{
		22:   "SSH",
		3389: "RDP",
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
