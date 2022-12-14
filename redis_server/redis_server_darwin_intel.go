//go:build darwin_intel
// +build darwin_intel

package redis_server

import (
	_ "embed"
)

//go:embed redis-server-darwin-intel
var RedisServer []byte