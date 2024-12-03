package main

import (
	"os"

	"github.com/caarlos0/log"

	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
	"paretosecurity.com/auditor/cmd"
	"paretosecurity.com/auditor/shared"
	"paretosecurity.com/auditor/team"
)

var rootCmd = &cobra.Command{
	Use:   "pareto [command]",
	Short: "Pareto Security CLI",
	Long: `Pareto Security CLI is a tool for running and reporting audits 
to paretosecurity.com.`,
}

var checkCmd = &cobra.Command{
	Use:   "check [--json]",
	Short: "Check system status",
	Run: func(cc *cobra.Command, args []string) {
		jsonOutput, _ := cc.Flags().GetBool("json")
		if jsonOutput {
			cmd.CheckJSON()
			return
		}
		cmd.Check()
	},
}

func init() { // enable colored output on github actions et al

	// automatically set GOMAXPROCS to match available CPUs.
	// GOMAXPROCS will be used as the default value for the --parallelism flag.
	if _, err := maxprocs.Set(); err != nil {
		log.WithError(err).Warn("failed to set GOMAXPROCS")
	}
	checkCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link team to this device",
	Run: func(cc *cobra.Command, args []string) {
		err := team.LinkAndWaitForTicket()
		if err != nil {
			log.WithError(err).Warn("failed to link")
			os.Exit(1)
		}
	},
}

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

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run as a daemon",
	Run: func(cc *cobra.Command, args []string) {
		log.Info("Starting daemon mode...")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(unlinkCmd)
}

func main() {

	if err := shared.LoadConfig(); err != nil {
		log.WithError(err).Warn("failed to load config")
	}

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Warn("failed to execute command")
		os.Exit(1)
	}
}
