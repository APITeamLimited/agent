//go:build darwin_intel
// +build darwin_intel

package redis_server

import (
	_ "embed"
)

//go:embed redis-server-darwin-intel
var RedisServer []byte

// Spawns child redis processes, these are terminated automatically when the agent exits
func SpawnChildServers(windowsTerminationChan chan bool) {
	spawnOrchestratorRedisUnix(windowsTerminationChan)
}
