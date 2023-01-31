package agent

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/APITeamLimited/agent/agent/libAgent"
	"github.com/APITeamLimited/agent/logo"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
)

const agentVersion = "v0.1.20"

func Run() {
	ensureNotAlreadyRunning()

	fmt.Printf("Running APITeam Agent %s\n", agentVersion)

	systrayContent := func() {
		logoIcon := logo.AgentLogo
		systray.SetIcon(logoIcon)
		systray.SetTitle("APITeam Agent")
		systray.SetTooltip("APITeam Agent")

		// Add non clickable menu item with name and icon
		mTitle := systray.AddMenuItem("About APITeam Agent", "About APITeam Agent")

		systray.AddSeparator()
		mAbortAll := systray.AddMenuItem("Abort All", "Abort All")
		mQuit := systray.AddMenuItem("Quit", "Quit APITeam Agent")

		setJobCountFunc := func(count int) {
			if count == 0 {
				mAbortAll.Hide()
			} else {
				mAbortAll.SetTitle(fmt.Sprintf("Abort All (%d)", count))
				mAbortAll.Show()
			}
		}

		setJobCountFunc(0)

		go func() {
			<-mTitle.ClickedCh
			// Open the URL in the default browser
			err := browser.OpenURL("https://apiteam.cloud/agent")
			if err != nil {
				panic(err)
			}
		}()

		setupChildProcesses()
		go runAgentServer(mAbortAll.ClickedCh, setJobCountFunc)

		// Wait for the server to stop before exiting
		<-mQuit.ClickedCh
		systray.Quit()
	}

	systray.Run(systrayContent, func() {
		stopRedisClients()
		os.Exit(0)
	})
}

func ensureNotAlreadyRunning() {
	serverAddress := fmt.Sprintf("http://localhost:%d/version", libAgent.AgentPort)
	// Ping the server to see if it's already running
	_, err := http.Get(serverAddress)

	if err != nil {
		return
	}

	// show popup
	dialog.Message("APITeam Agent is already running. Please close the existing instance before starting a new one.").Title("APITeam Agent").Error()

	os.Exit(0)
}
