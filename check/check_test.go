package check

import (
	"testing"
	"time"

	"github.com/ParetoSecurity/pareto-linux/shared"
)

type MockCheck struct {
	uuid       string
	passed     bool
	isRunnable bool
}

func (m *MockCheck) Name() string          { return "MockCheck" }
func (m *MockCheck) PassedMessage() string { return "Passed" }
func (m *MockCheck) FailedMessage() string { return "Failed" }
func (m *MockCheck) Run() error            { return nil }
func (m *MockCheck) Passed() bool          { return m.passed }
func (m *MockCheck) IsRunnable() bool      { return m.isRunnable }
func (m *MockCheck) UUID() string          { return m.uuid }
func (m *MockCheck) Status() string        { return "Status" }
func (m *MockCheck) RequiresRoot() bool    { return false }

func TestRegister(t *testing.T) {
	shared.Config.Checks = make(map[string]shared.CheckStatus)

	mockCheck := &MockCheck{
		uuid:       "1234",
		passed:     true,
		isRunnable: true,
	}

	registeredCheck := Register(mockCheck)

	if registeredCheck.UUID() != mockCheck.UUID() {
		t.Errorf("Expected UUID %s, got %s", mockCheck.UUID(), registeredCheck.UUID())
	}

	if len(shared.Config.Checks) != 1 {
		t.Errorf("Expected 1 check in the map, got %d", len(shared.Config.Checks))
	}

	checkStatus, exists := shared.Config.Checks[mockCheck.UUID()]
	if !exists {
		t.Errorf("Check with UUID %s not found in the map", mockCheck.UUID())
	}

	if checkStatus.Passed != mockCheck.Passed() {
		t.Errorf("Expected Passed %v, got %v", mockCheck.Passed(), checkStatus.Passed)
	}

	if checkStatus.Disabled != !mockCheck.IsRunnable() {
		t.Errorf("Expected Disabled %v, got %v", !mockCheck.IsRunnable(), checkStatus.Disabled)
	}

	if time.Since(checkStatus.UpdatedAt) > time.Second {
		t.Errorf("Expected UpdatedAt to be recent, got %v", checkStatus.UpdatedAt)
	}
}
