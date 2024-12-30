//go:build linux
// +build linux

package shared

type ReportingDevice struct {
	MachineUUID string `json:"machineUUID"` // e.g. 123e4567-e89b-12d3-a456-426614174000
	MachineName string `json:"machineName"` // e.g. MacBook-Pro.local
	Auth        string `json:"auth"`
	OSVersion   string `json:"linuxOSVersion"` // e.g. Ubuntu 20.04
	ModelName   string `json:"modelName"`      // e.g. MacBook Pro
	ModelSerial string `json:"modelSerial"`    // e.g. C02C1234
}
