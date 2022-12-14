//go:build linux || darwin_intel || darwin_arm
// +build linux darwin_intel darwin_arm

package logo

import (
	_ "embed"
)

//go:embed apiteam-logo.png
var AgentLogo []byte
