//go:build windows
// +build windows

package winapi

func convertKeyCode(code uintptr) (string, bool) {
	if '0' <= code && code <= '9' || 'A' <= code && code <= 'Z' {
		return string(rune(code)), true
	}
	var r string

	switch code {
	case VK_ESCAPE:
		r = NameEscape
	case VK_LEFT:
		r = NameLeftArrow
	case VK_RIGHT:
		r = NameRightArrow
	case VK_RETURN:
		r = NameReturn
	case VK_UP:
		r = NameUpArrow
	case VK_DOWN:
		r = NameDownArrow
	case VK_HOME:
		r = NameHome
	case VK_END:
		r = NameEnd
	case VK_BACK:
		r = NameDeleteBackward
	case VK_DELETE:
		r = NameDeleteForward
	case VK_PRIOR:
		r = NamePageUp
	case VK_NEXT:
		r = NamePageDown
	case VK_F1:
		r = NameF1
	case VK_F2:
		r = NameF2
	case VK_F3:
		r = NameF3
	case VK_F4:
		r = NameF4
	case VK_F5:
		r = NameF5
	case VK_F6:
		r = NameF6
	case VK_F7:
		r = NameF7
	case VK_F8:
		r = NameF8
	case VK_F9:
		r = NameF9
	case VK_F10:
		r = NameF10
	case VK_F11:
		r = NameF11
	case VK_F12:
		r = NameF12
	case VK_TAB:
		r = NameTab
	case VK_SPACE:
		r = NameSpace
	case VK_OEM_1:
		r = ";"
	case VK_OEM_PLUS:
		r = "+"
	case VK_OEM_COMMA:
		r = ","
	case VK_OEM_MINUS:
		r = "-"
	case VK_OEM_PERIOD:
		r = "."
	case VK_OEM_2:
		r = "/"
	case VK_OEM_3:
		r = "`"
	case VK_OEM_4:
		r = "["
	case VK_OEM_5, VK_OEM_102:
		r = "\\"
	case VK_OEM_6:
		r = "]"
	case VK_OEM_7:
		r = "'"
	case VK_CONTROL:
		r = NameCtrl
	case VK_SHIFT:
		r = NameShift
	case VK_MENU:
		r = NameAlt
	case VK_LWIN, VK_RWIN:
		r = NameSuper
	default:
		return "", false
	}
	return r, true
}
