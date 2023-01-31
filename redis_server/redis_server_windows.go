//go:build windows
// +build windows

package redis_server

import (
	_ "embed"
	"os"
	"os/exec"
	"syscall"

	"github.com/APITeamLimited/agent/agent/libAgent"
)

//go:embed redis-server-windows.exe
var RedisServer []byte

// Spawns child redis processes, these are terminated automatically when the agent exits
func SpawnChildServers(windowsTerminationChan chan bool) {
	redisPath := getRedisPath()

	go func() {
		// Execute redis-server.exe using os/exec
		orchestratorRedisCommand := exec.Command(redisPath, "--port", libAgent.OrchestratorRedisPort, "--save", "", "--appendonly", "no")

		orchestratorRedisCommand.SysProcAttr = &syscall.SysProcAttr{
			// Ignore this error
			// @go:linkname syscall_windows_hideWindow syscall.sub_windows_hideWindow
			HideWindow: true,
		}

		orchestratorRedisCommand.Stdout = os.Stdout
		err := orchestratorRedisCommand.Start()
		if err != nil {
			panic(err)
		}

		workerRedisCommand := exec.Command(redisPath, "--port", libAgent.WorkerRedisPort, "--save", "", "--appendonly", "no", "--protected-mode", "no")

		workerRedisCommand.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		workerRedisCommand.Stdout = os.Stdout
		err = workerRedisCommand.Start()
		if err != nil {
			panic(err)
		}

		// Wait for the agent to exit
		<-windowsTerminationChan

		// Kill the redis servers
		_ = orchestratorRedisCommand.Process.Kill()
		_ = workerRedisCommand.Process.Kill()
	}()

}
