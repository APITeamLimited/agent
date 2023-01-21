package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

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
	redisPath := getRedisPath()

	if runtime.GOOS == "windows" {
		// Execute redis-server.exe using os/exec
		orchestratorRedisCommand := exec.Command(redisPath, "--port", libAgent.OrchestratorRedisPort, "--save", "", "--appendonly", "no")

		orchestratorRedisCommand.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		orchestratorRedisCommand.Stdout = os.Stdout
		err := orchestratorRedisCommand.Start()
		if err != nil {
			panic(err)
		}

		workerRedisCommand := exec.Command(redisPath, "--port", libAgent.WorkerRedisPort, "--save", "", "--appendonly", "no")

		workerRedisCommand.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		workerRedisCommand.Stdout = os.Stdout
		err = workerRedisCommand.Start()
		if err != nil {
			panic(err)
		}

		return
	}

	be, err := byteexec.New(redis_server.RedisServer, redisPath)
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
		userName := os.Getenv("USERNAME")

		// Return path to temp directory
		return fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\Temp\\redis-server-windows", userName)
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
