package cmd

import (
	"os"
	"runtime"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("%s@%s %s", shared.Version, shared.Commit, shared.Date)
		log.Infof("Built with %s", runtime.Version())

		device := shared.CurrentReportingDevice()
		log.Infof("Machine UUID: %s", device.MachineUUID)
		log.Infof("Name: %s", device.MachineName)
		log.Infof("OS Version: %s", device.OSVersion)
		log.Infof("Model Name: %s", device.ModelName)
		log.Infof("Model Serial: %s", device.ModelSerial)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
