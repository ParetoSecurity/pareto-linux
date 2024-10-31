package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/log"

	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
	"paretosecurity.com/auditor/cmd"
	"paretosecurity.com/auditor/shared"
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
	Use:   "link [url]",
	Short: "Link team to this device",
	Args: func(cc *cobra.Command, args []string) error {
		if len(args) < 1 {
			log.Error("requires a URL argument")
			return nil
		}
		if !shared.IsValidParetoURL(args[0]) {
			log.Error("invalid URL, must start with paretosecurity://")
			return nil
		}
		return nil
	},
	Run: func(cc *cobra.Command, args []string) {
		//team.LinkWithDevice(args[0])
	},
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink team",
	Run: func(cc *cobra.Command, args []string) {
		fmt.Println("Unlinking system components...")
	},
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run as a daemon",
	Run: func(cc *cobra.Command, args []string) {
		fmt.Println("Starting daemon mode...")
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
