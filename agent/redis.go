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
	"github.com/getlantern/byteexec"
)

func setupChildProcesses() {
	spawnChildServers()
	runOrchestrator()
	runWorker()
}

// Spawns child redis processes, these are terminated automatically when the agent exits
func spawnChildServers() {
	be, err := byteexec.New(redis_server.RedisServer, getRedisPath())
	if err != nil {
		panic(err)
	}

	err = be.Command("--port", libAgent.OrchestratorRedisPort, "--save", "", "--appendonly", "no").Start()
	if err != nil {
		panic(err)
	}

	err = be.Command("--port", libAgent.WorkerRedisPort, "--save", "", "--appendonly", "no").Start()
	if err != nil {
		panic(err)
	}
}

func runOrchestrator() {
	// Set some environment variables
	_ = os.Setenv("ORCHESTRATOR_REDIS_HOST", libAgent.OrchestratorRedisHost)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PORT", libAgent.OrchestratorRedisPort)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PASSWORD", "")

	go orchestrator.Run(false, false)
}

func runWorker() {
	// Set some environment variables
	_ = os.Setenv("CLIENT_HOST", libAgent.WorkerRedisHost)
	_ = os.Setenv("CLIENT_PORT", libAgent.WorkerRedisPort)
	_ = os.Setenv("CLIENT_PASSWORD", "")

	go worker.Run(false)
}

func getRedisPath() string {
	system := runtime.GOOS

	if system == "windows" {
		// Return path to temp directory
		return "redis-server-windows"
	} else if system == "linux" {
		return "redis-server-linux"
	} else if system == "darwin" {
		if runtime.GOARCH == "arm64" {
			return "redis-server-darwin-arm"
		} else {
			return "redis-server-darwin-intel"
		}
	}

	panic(fmt.Sprintf("Unsupported system: %s, arch: %s", system, runtime.GOARCH))
}

// Instructs the redis servers to immediately terminate
func stopRedisClients() {
	orchestratorRedis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libAgent.OrchestratorRedisHost, libAgent.OrchestratorRedisPort),
		Username: "default",
		Password: "",
	})

	workerRedis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libAgent.WorkerRedisHost, libAgent.WorkerRedisPort),
		Username: "default",
		Password: "",
	})

	orchestratorRedis.ShutdownNoSave(context.Background())
	workerRedis.ShutdownNoSave(context.Background())
}
