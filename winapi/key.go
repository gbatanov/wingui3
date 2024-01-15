package winapi

func convertKeyCode(code uintptr) (string, bool) {
	if '0' <= code && code <= '9' ||
		'A' <= code && code <= 'Z' ||
		'a' <= code && code <= 'z' {
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

// Modifiers
type Modifiers uint32

const (
	// ModCtrl is the ctrl modifier
	ModCtrl Modifiers = 1 << iota
	// ModCommand is the command modifier key
	// found on Apple keyboards.
	ModCommand
	// ModShift is the shift modifier
	ModShift
	// ModAlt is the alt modifier key, or the option
	// key on Apple keyboards.
	ModAlt
	// ModSuper is the "logo" modifier key, often
	// represented by a Windows logo.
	ModSuper
)

func getModifiers() Modifiers {
	var kmods Modifiers
	if GetKeyState(VK_LWIN)&0x1000 != 0 || GetKeyState(VK_RWIN)&0x1000 != 0 {
		kmods |= ModSuper
	}
	if GetKeyState(VK_MENU)&0x1000 != 0 {
		kmods |= ModAlt
	}
	if GetKeyState(VK_CONTROL)&0x1000 != 0 {
		kmods |= ModCtrl
	}
	if GetKeyState(VK_SHIFT)&0x1000 != 0 {
		kmods |= ModShift
	}
	return kmods
}

const (
	// Names for special keys.
	NameLeftArrow      = "←"
	NameRightArrow     = "→"
	NameUpArrow        = "↑"
	NameDownArrow      = "↓"
	NameReturn         = "⏎"
	NameEnter          = "⌤"
	NameEscape         = "⎋"
	NameHome           = "⇱"
	NameEnd            = "⇲"
	NameDeleteBackward = "⌫"
	NameDeleteForward  = "⌦"
	NamePageUp         = "⇞"
	NamePageDown       = "⇟"
	NameTab            = "Tab"
	NameSpace          = "Space"
	NameCtrl           = "Ctrl"
	NameShift          = "Shift"
	NameAlt            = "Alt"
	NameSuper          = "Super"
	NameCommand        = "⌘"
	NameF1             = "F1"
	NameF2             = "F2"
	NameF3             = "F3"
	NameF4             = "F4"
	NameF5             = "F5"
	NameF6             = "F6"
	NameF7             = "F7"
	NameF8             = "F8"
	NameF9             = "F9"
	NameF10            = "F10"
	NameF11            = "F11"
	NameF12            = "F12"
	NameBack           = "Back"
)
