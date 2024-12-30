//go:build windows
// +build windows

package shared

import (
	_ "embed"
)

var (
	//go:embed icon_white.ico
	IconWhite []byte

	//go:embed icon_black.ico
	IconBlack []byte
)
