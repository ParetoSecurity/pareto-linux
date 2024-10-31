package cmd

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pterm/pterm"
	"paretosecurity.com/auditor/check"
	"paretosecurity.com/auditor/claims"
	"paretosecurity.com/auditor/shared"
)

func Check() {
	multi := pterm.DefaultMultiPrinter
	var wg sync.WaitGroup
	pterm.Println("Starting checks...")
	if _, err := multi.Start(); err != nil {
		fmt.Println(err)
	}
	for _, claim := range claims.All {
		for _, chk := range claim.Checks {
			wg.Add(1)
			go func(claim claims.Claim, chk check.Check) {
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
				if err := chk.Run(); err != nil {
					spinner.Fail(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Red(err.Error()))
				}

				if chk.Passed() {
					spinner.Success(pterm.White(claim.Title), pterm.White(": "), pterm.Green(chk.Status()))
				} else {
					spinner.Fail(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Red(chk.Status()))
				}
				wg.Done()
			}(claim, chk)
		}
	}
	wg.Wait()
	if _, err := multi.Stop(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
	pterm.Println("Checks completed.")
	if err := shared.SaveConfig(); err != nil {
		fmt.Println(err)
	}
}

func CheckJSON() {
	status := make(map[string]string)
	for _, claim := range claims.All {
		for _, chk := range claim.Checks {
			if err := chk.Run(); err != nil {
				status[chk.UUID()] = err.Error()
				continue
			}
			status[chk.UUID()] = chk.Status()
		}
	}
	if err := shared.SaveConfig(); err != nil {
		fmt.Println(err)
	}
	out, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
