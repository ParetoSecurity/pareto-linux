package check

import "os/exec"

type DockerAccess struct {
	passed bool
}

// Name returns the name of the check
func (f *DockerAccess) Name() string {
	return "Access to Docker is restricted"
}

// Run executes the check
func (f *DockerAccess) Run() error {
	cmd := exec.Command("docker", "run", "--rm", "hello-world")
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				f.passed = true
				return nil
			}
		}
	}
	f.passed = false
	return nil
}

// Passed returns the status of the check
func (f *DockerAccess) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *DockerAccess) IsRunnable() bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	return err == nil
}

// UUID returns the UUID of the check
func (f *DockerAccess) UUID() string {
	return "25443ceb-c1ec-408c-b4f3-2328ea0c84e1"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *DockerAccess) ReportIfDisabled() bool {
	return false
}

// PassedMessage returns the message to return if the check passed
func (f *DockerAccess) PassedMessage() string {
	return "Access to Docker is restricted"
}

// FailedMessage returns the message to return if the check failed
func (f *DockerAccess) FailedMessage() string {
	return "Access to Docker is not restricted"
}

// RequiresRoot returns whether the check requires root access
func (f *DockerAccess) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *DockerAccess) Status() string {
	if !f.Passed() {
		return f.FailedMessage()
	}
	return f.PassedMessage()
}
