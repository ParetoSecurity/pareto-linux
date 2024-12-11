package claims

import "paretosecurity.com/auditor/check"

var All = []Claim{
	{"Access Security", []check.Check{
		check.Register(&check.Autologin{}),
		check.Register(&check.DockerAccess{}),
		check.Register(&check.PasswordToUnlock{}),
		check.Register(&check.SSHKeys{}),
		check.Register(&check.SSHKeysAlgo{}),
	}},
	{"System Updates", []check.Check{
		check.Register(&check.SoftwareUpdates{}),
	}},
	{"Firewall & Sharing", []check.Check{
		check.Register(&check.Firewall{}),
		check.Register(&check.Printer{}),
		check.Register(&check.RemoteLogin{}),
		check.Register(&check.Sharing{}),
	}},

	{"System Integrity", []check.Check{
		check.Register(&check.SecureBoot{}),
		check.Register(&check.EncryptingFS{}),
	}},
}
