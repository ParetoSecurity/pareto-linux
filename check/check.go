package check

import (
	"time"

	"paretosecurity.com/auditor/shared"
)

type Check interface {
	Name() string
	Run() error
	Passed() bool
	IsRunnable() bool
	ReportIfDisabled() bool
	UUID() string
	Status() string
}

func Register(c Check) Check {

	AvailableChecks = +1

	// If the check is already in the checks map, return it
	if found := shared.Config.Checks[c.UUID()]; found != (shared.CheckStatus{}) {
		return c
	}

	// If the checks map is nil, create it
	if shared.Config.Checks == nil {
		shared.Config.Checks = make(map[string]shared.CheckStatus)
	}

	// Add the check to the checks map
	shared.Config.Checks[c.UUID()] = shared.CheckStatus{
		UpdatedAt: time.Now(),
		Passed:    c.Passed(),
		Disabled:  !c.IsRunnable(),
	}
	return c
}

func Update(c Check) Check {
	shared.Config.Checks[c.UUID()] = shared.CheckStatus{
		UpdatedAt: time.Now(),
		Passed:    c.Passed(),
		Disabled:  !c.IsRunnable(),
	}
	return c
}
