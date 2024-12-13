package claims

import (
	"paretosecurity.com/auditor/check"
	"paretosecurity.com/auditor/checks"
)

var All = []Claim{
	{"Access Security", []check.Check{
		check.Register(&checks.Autologin{}),
		check.Register(&checks.DockerAccess{}),
		check.Register(&checks.PasswordToUnlock{}),
		check.Register(&checks.SSHKeys{}),
		check.Register(&checks.SSHKeysAlgo{}),
		check.Register(&checks.SSHConfigCheck{}),
	}},
	{"System Updates", []check.Check{
		check.Register(&checks.SoftwareUpdates{}),
		check.Register(&checks.ParetoUpdated{}),
	}},
	{"Firewall & Sharing", []check.Check{
		check.Register(&checks.Firewall{}),
		check.Register(&checks.Printer{}),
		check.Register(&checks.RemoteLogin{}),
		check.Register(&checks.Sharing{}),
	}},
	{"System Integrity", []check.Check{
		check.Register(&checks.SecureBoot{}),
		check.Register(&checks.EncryptingFS{}),
	}},
}
