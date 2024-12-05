package cmd

import (
	"os"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
	"paretosecurity.com/auditor/shared"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink team",
	Run: func(cc *cobra.Command, args []string) {
		log.Info("Unlinking device ...")
		shared.Config.TeamID = ""
		shared.Config.AuthToken = ""
		if err := shared.SaveConfig(); err != nil {
			log.WithError(err).Warn("failed to save config")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}
