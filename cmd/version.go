package cmd

import (
	"os"

	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("%s@%s %s", shared.Version, shared.Commit, shared.Date)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
