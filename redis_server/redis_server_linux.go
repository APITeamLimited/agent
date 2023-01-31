//go:build linux
// +build linux

package redis_server

import (
	_ "embed"
)

//go:embed redis-server-linux
var RedisServer []byte

// Spawns child redis processes, these are terminated automatically when the agent exits
func SpawnChildServers(windowsTerminationChan chan bool) {
	spawnOrchestratorRedisUnix(windowsTerminationChan)
}
