//go:build darwin_arm
// +build darwin_arm

package redis_server

import (
	_ "embed"
)

//go:embed redis-server-darwin-arm
var RedisServer []byte
