package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ParetoSecurity/pareto-linux/check"
	"github.com/ParetoSecurity/pareto-linux/claims"
	"github.com/ParetoSecurity/pareto-linux/shared"
	"github.com/ParetoSecurity/pareto-linux/team"
	"github.com/caarlos0/log"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [--json] [--schema] [--install] [--uninstall]",
	Short: "Run checks",
	Run: func(cc *cobra.Command, args []string) {
		jsonOutput, _ := cc.Flags().GetBool("json")
		schemaOutput, _ := cc.Flags().GetBool("schema")
		installFlag, _ := cc.Flags().GetBool("install")
		uninstallFlag, _ := cc.Flags().GetBool("uninstall")

		if installFlag {
			installUserTimer()
			return
		}
		if uninstallFlag {
			uninstallUserTimer()
			return
		}
		if schemaOutput {
			PrintSchemaJSON()
			return
		}

		if jsonOutput {
			CheckJSON()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		done := make(chan struct{})
		go func() {
			Check(ctx)
			close(done)
		}()

		select {
		case <-done:
			if shared.IsLinked() {
				err := team.ReportToTeam()
				if err != nil {
					log.WithError(err).Warn("failed to report to team")
				}
			} else {
				showLinkingMessage()
			}
			if !isUserTimerInstalled() {
				log.Info("To ensure your system is checked every hour, please run `paretosecurity check --install` to set it up.")
			}
		case <-ctx.Done():
			log.Warn("Check run timed out")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().Bool("json", false, "output JSON")
	checkCmd.Flags().Bool("schema", false, "output schema for all checks")
	checkCmd.Flags().Bool("install", false, "setup hourly checks")
	checkCmd.Flags().Bool("uninstall", false, "remove hourly checks")
}

func showLinkingMessage() {
	log.Info("To link your account with the team, please run `paretosecurity link`.")
	log.Info("For more information, please visit https://paretosecurity.com/dashboard")
}

func Check(ctx context.Context) {
	multi := pterm.DefaultMultiPrinter
	var wg sync.WaitGroup
	log.Info("Starting checks...")
	if _, err := multi.Start(); err != nil {
		log.WithError(err).Warn("failed to stop multi printer")
	}
	for _, claim := range claims.All {
		for _, chk := range claim.Checks {
			wg.Add(1)
			go func(claim claims.Claim, chk check.Check) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
					spinner, _ := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start(fmt.Sprintf("%s: %s", claim.Title, chk.Name()))
					spinner.FailPrinter = &pterm.PrefixPrinter{
						MessageStyle: &pterm.Style{pterm.FgLightRed},
						Prefix: pterm.Prefix{
							Style: &pterm.Style{pterm.BgRed, pterm.FgLightRed},
							Text:  "✗",
						},
					}
					spinner.SuccessPrinter = &pterm.PrefixPrinter{
						MessageStyle: &pterm.Style{pterm.FgLightGreen},
						Prefix: pterm.Prefix{
							Style: &pterm.Style{pterm.BgGreen, pterm.FgLightGreen},
							Text:  "✓",
						},
					}

					// Skip checks that are not runnable
					if !chk.IsRunnable() {
						spinner.Warning(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Yellow("skipped"))
						return
					}

					if err := chk.Run(); err != nil {
						spinner.Fail(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Red(err.Error()))
					}

					if chk.Passed() {
						spinner.Success(pterm.White(claim.Title), pterm.White(": "), pterm.Green(chk.Status()))
					} else {
						spinner.Fail(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Red(chk.Status()))
					}
					shared.UpdateLastState(shared.LastState{
						UUID:    chk.UUID(),
						State:   chk.Passed(),
						Details: chk.Status(),
					})
				}
			}(claim, chk)
		}
	}
	wg.Wait()
	if err := shared.CommitLastState(); err != nil {
		log.WithError(err).Warn("failed to commit last state")
	}
	if _, err := multi.Stop(); err != nil {
		log.WithError(err).Warn("failed to stop multi printer")
	}

	log.Info("Checks completed.")
	if err := shared.SaveConfig(); err != nil {
		log.WithError(err).Warn("cannot save config")
	}
}

func CheckJSON() {
	status := make(map[string]string)
	for _, claim := range claims.All {
		for _, chk := range claim.Checks {

			if !chk.IsRunnable() {
				status[chk.UUID()] = "skipped"
				continue
			}

			if err := chk.Run(); err != nil {
				status[chk.UUID()] = err.Error()
				continue
			}
			if chk.Passed() {
				status[chk.UUID()] = "passed"
			} else {
				status[chk.UUID()] = "failed"
			}
		}
	}
	if err := shared.SaveConfig(); err != nil {
		log.WithError(err).Warn("cannot save config")
	}
	out, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.WithError(err).Warn("cannot marshal status")
	}
	fmt.Println(string(out))
}

func PrintSchemaJSON() {
	schema := make(map[string]map[string][]string)
	for _, claim := range claims.All {
		checks := make(map[string][]string)
		for _, chk := range claim.Checks {
			checks[chk.UUID()] = []string{chk.PassedMessage(), chk.FailedMessage()}
		}
		schema[claim.Title] = checks
	}
	out, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.WithError(err).Warn("cannot marshal schema")
	}
	fmt.Println(string(out))
}
