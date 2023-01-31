package redis_server

import (
	"fmt"
	"os"
	"runtime"

	"github.com/APITeamLimited/agent/agent/libAgent"
	"github.com/getlantern/byteexec"
)

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

// Non-windows archs use byteexec to embed the redis-server binary
func spawnOrchestratorRedisUnix(windowsTerminationChan chan bool) {
	redisPath := getRedisPath()

	be, err := byteexec.New(RedisServer, redisPath)
	if err != nil {
		panic(err)
	}

	err = be.Command("--port", libAgent.OrchestratorRedisPort, "--save", "", "--appendonly", "no").Start()
	if err != nil {
		panic(err)
	}
}
