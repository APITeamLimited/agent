package agent

import (
	"fmt"
	"os"
	"runtime"

	"github.com/APITeamLimited/agent/agent/libAgent"
	"github.com/APITeamLimited/agent/redis_server"
	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"
	"github.com/getlantern/byteexec"
)

func setupChildProcesses() {
	spawnChildServers()
	runOrchestrator()
	runWorker()
}

// Spawns child redis processes, these are terminated automatically when the agent exits
func spawnChildServers() {
	be, err := byteexec.New(redis_server.RedisServer, getRedisFileName())
	if err != nil {
		panic(err)
	}

	err = be.Command(fmt.Sprintf("--port %s", libAgent.OrchestratorRedisPort), "--save", "", "--appendonly", "no").Start()
	if err != nil {
		panic(err)
	}

	err = be.Command(fmt.Sprintf("--port %s", libAgent.WorkerRedisPort), "--save", "", "--appendonly", "no").Start()
	if err != nil {
		panic(err)
	}
}

func runOrchestrator() {
	// Set some environment variables
	_ = os.Setenv("ORCHESTRATOR_REDIS_HOST", libAgent.OrchestratorRedisHost)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PORT", libAgent.OrchestratorRedisPort)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PASSWORD", "")

	go orchestrator.Run(false)
}

func runWorker() {
	// Set some environment variables
	_ = os.Setenv("CLIENT_HOST", libAgent.WorkerRedisHost)
	_ = os.Setenv("CLIENT_PORT", libAgent.WorkerRedisPort)
	_ = os.Setenv("CLIENT_PASSWORD", "")

	go worker.Run(false)
}

func getRedisFileName() string {
	system := runtime.GOOS

	if system == "windows" {
		return "redis-server-agent.exe"
	}

	return "redis-server-agent"
}
