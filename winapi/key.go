package winapi

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
