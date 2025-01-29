package claims

import (
	"github.com/ParetoSecurity/pareto-core/check"
	shared "github.com/ParetoSecurity/pareto-core/checks/shared"
)

var All = []Claim{
	{"Access Security", []check.Check{
		check.Register(&shared.SSHKeys{}),
		check.Register(&shared.SSHKeysAlgo{}),
	}},
	{"Software Updates", []check.Check{
		check.Register(&shared.ParetoUpdated{}),
	}},
	{"Firewall & Sharing", []check.Check{
		check.Register(&shared.RemoteLogin{}),
	}},
	{"System Integrity", []check.Check{}},
}
