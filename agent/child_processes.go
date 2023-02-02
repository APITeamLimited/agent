package agent

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/APITeamLimited/agent/agent/libAgent"
	"github.com/APITeamLimited/agent/redis_server"
	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"
	"github.com/APITeamLimited/redis/v9"
)

func setupChildProcesses(windowsTerminationChan chan bool) {
	redis_server.SpawnChildServers(windowsTerminationChan)
	runOrchestrator()
	runWorkerServer()
}

func runOrchestrator() {
	// Set some environment variables
	_ = os.Setenv("ORCHESTRATOR_REDIS_HOST", libAgent.OrchestratorRedisHost)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PORT", libAgent.OrchestratorRedisPort)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PASSWORD", "")
	_ = os.Setenv("SERVICE_URL_OVERRIDE_0", fmt.Sprintf("http://localhost:%s", libAgent.WorkerServerPort))

	go orchestrator.Run(false)
}

func runWorkerServer() {
	_ = os.Setenv("WORKER_SERVER_PORT", libAgent.WorkerServerPort)
	_ = os.Setenv("WORKER_0_DISPLAY_NAME", libAgent.AgentWorkerName)
	_ = os.Setenv("WORKER_STANDALONE", "false")

	go worker.RunWorkerServer()
}

// Instructs the redis servers to immediately terminate
func stopRedisClients(windowsTerminationChan chan bool) {
	if runtime.GOOS == "windows" {
		windowsTerminationChan <- true
		return
	}

	orchestratorRedis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libAgent.OrchestratorRedisHost, libAgent.OrchestratorRedisPort),
		Username: "default",
		Password: "",
	})

	orchestratorRedis.ShutdownNoSave(context.Background())
}
