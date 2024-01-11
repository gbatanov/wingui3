//go:build windows
// +build windows

package img

import (
	_ "embed"
)

//go:embed check.ico
var OkIco []byte

//go:embed stop.ico
var ErrIco []byte
