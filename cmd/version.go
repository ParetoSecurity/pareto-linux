package cmd

import (
	"os"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
	"paretosecurity.com/auditor/shared"
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
