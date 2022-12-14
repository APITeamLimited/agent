//go:build windows
// +build windows

package logo

import (
	_ "embed"
)

//go:embed apiteam-logo.ico
var AgentLogo []byte
