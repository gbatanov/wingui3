package winapi

import (
	"image"
	"sync"
)

type Stage uint8

const (
	// StagePaused is the stage for windows that have no on-screen representation.
	// Paused windows don't receive FrameEvent.
	StagePaused Stage = iota
	// StageInactive is the stage for windows that are visible, but not active.
	// Inactive windows receive FrameEvent.
	StageInactive
	// StageRunning is for active and visible
	// Running windows receive FrameEvent.
	StageRunning
)

// winMap maps win32 HWNDs to *Window
var WinMap sync.Map

type WindowMode uint8

const (
	// Windowed is the normal window mode with OS specific window decorations.
	Windowed WindowMode = iota
	// Fullscreen is the full screen window mode.
	Fullscreen
	// Minimized is for systems where the window can be minimized to an icon.
	Minimized
	// Maximized is for systems where the window can be made to fill the available monitor area.
	Maximized
)

type Config struct {
	ID         uintptr // используется в дочерних активных элементах, как hMenu
	Position   image.Point
	Size       image.Point
	MinSize    image.Point
	MaxSize    image.Point
	Mode       WindowMode
	SysMenu    int // 0 - нет шапки, 1- только заголовок, 2 - иконка и кнопка закрытия
	Title      string
	EventChan  chan Event
	BorderSize image.Point
	TextColor  uint32
	FontSize   int32
	BgColor    uint32
	Class      string
}

const (
	VK_CONTROL = 0x11
	VK_LWIN    = 0x5B
	VK_MENU    = 0x12
	VK_RWIN    = 0x5C
	VK_SHIFT   = 0x10

	VK_BACK   = 0x08
	VK_DELETE = 0x2e
	VK_DOWN   = 0x28
	VK_END    = 0x23
	VK_ESCAPE = 0x1b
	VK_HOME   = 0x24
	VK_LEFT   = 0x25
	VK_NEXT   = 0x22
	VK_PRIOR  = 0x21
	VK_RIGHT  = 0x27
	VK_RETURN = 0x0d
	VK_SPACE  = 0x20
	VK_TAB    = 0x09
	VK_UP     = 0x26

	VK_F1  = 0x70
	VK_F2  = 0x71
	VK_F3  = 0x72
	VK_F4  = 0x73
	VK_F5  = 0x74
	VK_F6  = 0x75
	VK_F7  = 0x76
	VK_F8  = 0x77
	VK_F9  = 0x78
	VK_F10 = 0x79
	VK_F11 = 0x7A
	VK_F12 = 0x7B

	VK_OEM_1      = 0xba
	VK_OEM_PLUS   = 0xbb
	VK_OEM_COMMA  = 0xbc
	VK_OEM_MINUS  = 0xbd
	VK_OEM_PERIOD = 0xbe
	VK_OEM_2      = 0xbf
	VK_OEM_3      = 0xc0
	VK_OEM_4      = 0xdb
	VK_OEM_5      = 0xdc
	VK_OEM_6      = 0xdd
	VK_OEM_7      = 0xde
	VK_OEM_102    = 0xe2

	WM_CREATE               = 0x0001
	WM_DESTROY              = 0x0002
	WM_MOVE                 = 0x0003
	WM_SIZE                 = 0x0005
	WM_SETFOCUS             = 0x0007
	WM_KILLFOCUS            = 0x0008
	WM_PAINT                = 0x000F
	WM_CLOSE                = 0x0010
	WM_QUIT                 = 0x0012
	WM_ERASEBKGND           = 0x0014
	WM_SHOWWINDOW           = 0x0018
	WM_ACTIVATEAPP          = 0x001C
	WM_CANCELMODE           = 0x001F
	WM_SETCURSOR            = 0x0020
	WM_CHILDACTIVATE        = 0x0022
	WM_GETMINMAXINFO        = 0x0024
	WM_WINDOWPOSCHANGED     = 0x0047
	WM_NOTIFY               = 0x004E
	WM_NCCREATE             = 0x0081
	WM_NCCALCSIZE           = 0x0083
	WM_NCHITTEST            = 0x0084
	WM_NCACTIVATE           = 0x0086
	WM_KEYDOWN              = 0x0100
	WM_KEYUP                = 0x0101
	WM_CHAR                 = 0x0102
	WM_SYSKEYDOWN           = 0x0104
	WM_SYSKEYUP             = 0x0105
	WM_UNICHAR              = 0x0109
	WM_IME_STARTCOMPOSITION = 0x010D
	WM_IME_ENDCOMPOSITION   = 0x010E
	WM_IME_COMPOSITION      = 0x010F
	WM_COMMAND              = 0x0111
	WM_TIMER                = 0x0113
	WM_CTLCOLORSTATIC       = 0x0138
	WM_MOUSEMOVE            = 0x0200
	WM_LBUTTONDOWN          = 0x0201
	WM_LBUTTONUP            = 0x0202
	WM_RBUTTONDOWN          = 0x0204
	WM_RBUTTONUP            = 0x0205
	WM_MBUTTONDOWN          = 0x0207
	WM_MBUTTONUP            = 0x0208
	WM_MOUSEWHEEL           = 0x020A
	WM_MOUSEHWHEEL          = 0x020E
	WM_DPICHANGED           = 0x02E0
	WM_USER                 = 0x0400
)
