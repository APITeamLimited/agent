//go:build windows
// +build windows

package redis_server

import (
	_ "embed"
)

//go:embed redis-server-windows.exe
var RedisServer []byte
