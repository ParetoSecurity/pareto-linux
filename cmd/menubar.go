package cmd

import (
	"fmt"
	"strings"

	"os/exec"

	"fyne.io/systray"

	"github.com/caarlos0/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"paretosecurity.com/auditor/check"
	"paretosecurity.com/auditor/claims"
	"paretosecurity.com/auditor/shared"
)

func addQuitItem() {
	mQuit := systray.AddMenuItem("Quit", "Quit the Pareto Security")
	mQuit.Enable()
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	systray.AddSeparator()
}

func checkStatusToIcon(status bool) string {
	if status {
		return "✔"
	}
	return "✘"
}

func getIcon() []byte {

	isDark, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme").Output()
	if err == nil && strings.Contains(string(isDark), "prefer-dark") {
		return shared.IconWhite
	}
	isKDE, err := exec.Command("kreadconfig5", "--group", "General", "--key", "ColorScheme").Output()
	if err == nil && strings.Contains(string(isKDE), "Dark") {
		return shared.IconWhite
	}

	return shared.IconBlack
}

func onReady() {
	systray.SetTemplateIcon(shared.IconBlack, shared.IconBlack)
	systray.SetTemplateIcon(getIcon(), getIcon())
	systray.SetTooltip("Pareto Security")
	systray.AddMenuItem("Pareto Security", "").Disable()
	systray.AddSeparator()
	for _, claim := range claims.All {
		mClaim := systray.AddMenuItem(claim.Title, "")
		allStatus := lo.Reduce(claim.Checks, func(acc bool, item check.Check, index int) bool {
			checkStatus, found, _ := shared.GetLastState(item.UUID())
			if !item.IsRunnable() {
				return acc && true
			}
			return acc && checkStatus.State && found
		}, true)

		mClaim.SetTitle(fmt.Sprintf("%s %s", checkStatusToIcon(allStatus), claim.Title))

		for _, chk := range claim.Checks {
			checkStatus, found, _ := shared.GetLastState(chk.UUID())
			state := chk.Passed()
			if found {
				state = checkStatus.State
			}
			mCheck := mClaim.AddSubMenuItem(fmt.Sprintf("%s %s", checkStatusToIcon(state), chk.Name()), "")
			if !chk.IsRunnable() {
				mCheck.Disable()
			}
			go func(chk check.Check, mCheck *systray.MenuItem) {
				for range mCheck.ClickedCh {
					err := exec.Command("open", fmt.Sprintf("https://paretosecurity.com/checks/%s?details=None", chk.UUID())).Run()
					if err != nil {
						log.WithError(err).Error("failed to open check URL")
					}
				}
			}(chk, mCheck)
		}
	}
	systray.AddSeparator()
	addQuitItem()
}

var menubarCmd = &cobra.Command{
	Use:   "menubar",
	Short: "Show the checks in the menubar",
	Run: func(cc *cobra.Command, args []string) {
		onExit := func() {
			log.Info("Exiting...")
		}
		systray.Run(onReady, onExit)
	},
}

func init() {
	rootCmd.AddCommand(menubarCmd)
}