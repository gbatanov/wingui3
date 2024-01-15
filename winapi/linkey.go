//go:build linux
// +build linux

package winapi

import "github.com/jezek/xgb/xproto"

func convertKeyCode(code xproto.Keycode) (string, bool) {
	return "", false
}
