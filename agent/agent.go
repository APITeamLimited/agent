package agent

import (
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/APITeamLimited/agent/agent/libAgent"
	"github.com/APITeamLimited/agent/logo"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
)

const agentVersion = "v0.1.20"

func Run() {
	ensureNotAlreadyRunning()
	determineIfPortsFree()

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
			// CLicked callback
			for {
				<-mTitle.ClickedCh
				// Open the URL in the default browser
				browser.OpenURL("https://apiteam.cloud/agent")
			}
		}()

		// This is used for telling windows child redis clients to terminate
		windowsTerminationChan := make(chan bool)

		setupChildProcesses(windowsTerminationChan)
		go runAgentServer(mAbortAll.ClickedCh, setJobCountFunc)

		// Wait for the server to stop before exiting
		<-mQuit.ClickedCh
		stopRedisClients(windowsTerminationChan)
		systray.Quit()
	}

	systray.Run(systrayContent, func() {
		os.Exit(0)
	})
}

func ensureNotAlreadyRunning() {
	serverAddress := fmt.Sprintf("http://localhost:%s/version", libAgent.AgentPort)
	// Ping the server to see if it's already running
	_, err := http.Get(serverAddress)

	if err != nil {
		return
	}

	// show popup
	dialog.Message("APITeam Agent is already running. Please close the existing instance before starting a new one.").Title("APITeam Agent").Error()

	os.Exit(0)
}

func determineIfPortsFree() {
	// Check if libAgent.AgentPort is free
	// Check if libAgent.WorkerServerPort is free

	// If not free, show popup

	for _, port := range []string{libAgent.AgentPort, libAgent.WorkerServerPort} {
		if isOpened(port) {
			dialog.Message(fmt.Sprintf("Port %s is already in use and is required by APITeam Agent. Please close the application using it and try again.", port)).Title("APITeam Agent").Error()
			os.Exit(0)
		}
	}
}

func isOpened(port string) bool {
	timeout := 30 * time.Millisecond
	target := fmt.Sprintf("localhost:%s", port)

	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()
		return true
	}

	return false
}
