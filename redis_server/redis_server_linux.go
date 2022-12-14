//go:build linux
// +build linux

package redis_server

import (
	_ "embed"
)

//go:embed redis-server-linux
var RedisServer []byte
