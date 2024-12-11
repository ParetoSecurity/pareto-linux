package cmd

import (
	"os"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
	"paretosecurity.com/auditor/shared"
	"paretosecurity.com/auditor/team"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link team with this device",
	Run: func(cc *cobra.Command, args []string) {
		if team.IsLinked() {
			log.Warn("Already linked to a team")
			log.Warn("Unlink first with `pareto unlink`")
			log.Infof("Team ID: %s", shared.Config.TeamID)
			os.Exit(1)
		}
		err := team.LinkAndWaitForTicket()
		if err != nil {
			log.WithError(err).Warn("failed to link")
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
