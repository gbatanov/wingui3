// SPDX-License-Identifier: Unlicense OR MIT

//go:build windows
// +build windows

package winapi

import (
	syscall "golang.org/x/sys/windows"
)

const (
	ANSI_CHARSET    = 0
	DEFAULT_CHARSET = 1
	RUSSIAN_CHARSET = 204
)

type PAINTSTRUCT struct {
	hdc         syscall.Handle
	fErase      bool
	rcPaint     Rect
	fRestore    bool
	fIncUpdate  bool
	rgbReserved uint32
}

// NRGBA represents a non-alpha-premultiplied 32-bit color.
type NRGBA struct {
	R, G, B, A uint8
}

type CompositionForm struct {
	dwStyle      uint32
	ptCurrentPos Point
	rcArea       Rect
}

type CandidateForm struct {
	dwIndex      uint32
	dwStyle      uint32
	ptCurrentPos Point
	rcArea       Rect
}

type Rect struct {
	Left, Top, Right, Bottom int32
}

type WndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

type Margins struct {
	CxLeftWidth    int32
	CxRightWidth   int32
	CyTopHeight    int32
	CyBottomHeight int32
}

type Msg struct {
	Hwnd     syscall.Handle
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

type Point struct {
	X, Y int32
}

type MinMaxInfo struct {
	PtReserved     Point
	PtMaxSize      Point
	PtMaxPosition  Point
	PtMinTrackSize Point
	PtMaxTrackSize Point
}

type NCCalcSizeParams struct {
	Rgrc  [3]Rect
	LpPos *WindowPos
}

type WindowPos struct {
	HWND            syscall.Handle
	HWNDInsertAfter syscall.Handle
	x               int32
	y               int32
	cx              int32
	cy              int32
	flags           uint32
}

type WindowPlacement struct {
	length           uint32
	flags            uint32
	showCmd          uint32
	ptMinPosition    Point
	ptMaxPosition    Point
	rcNormalPosition Rect
	rcDevice         Rect
}

type MonitorInfo struct {
	cbSize   uint32
	Monitor  Rect
	WorkArea Rect
	Flags    uint32
}

const (
	TRUE       = 1
	IDOK       = 1
	IDCANCEL   = 2
	IDABORT    = 3
	IDRETRY    = 4
	IDIGNORE   = 5
	IDYES      = 6
	IDNO       = 7
	IDCLOSE    = 8
	IDHELP     = 9
	IDTRYAGAIN = 10
	IDCONTINUE = 11

	CPS_CANCEL = 0x0004

	CS_HREDRAW     = 0x0002
	CS_INSERTCHAR  = 0x2000
	CS_NOMOVECARET = 0x4000
	CS_VREDRAW     = 0x0001
	CS_OWNDC       = 0x0020

	CW_USEDEFAULT = -2147483648

	GWL_STYLE = ^(uintptr(16) - 1) // -16

	GCS_COMPSTR       = 0x0008
	GCS_COMPREADSTR   = 0x0001
	GCS_CURSORPOS     = 0x0080
	GCS_DELTASTART    = 0x0100
	GCS_RESULTREADSTR = 0x0200
	GCS_RESULTSTR     = 0x0800

	CFS_POINT        = 0x0002
	CFS_CANDIDATEPOS = 0x0040

	HWND_TOPMOST = uint32(^(int32(0) - 1)) // в такой записи работает функция, не дает ошибку
	HWND_TOP     = 0

	HTCAPTION     = 2
	HTCLIENT      = 1
	HTLEFT        = 10
	HTRIGHT       = 11
	HTTOP         = 12
	HTTOPLEFT     = 13
	HTTOPRIGHT    = 14
	HTBOTTOM      = 15
	HTBOTTOMLEFT  = 16
	HTBOTTOMRIGHT = 17

	IDC_APPSTARTING = 32650 // Standard arrow and small hourglass
	IDC_ARROW       = 32512 // Standard arrow
	IDC_CROSS       = 32515 // Crosshair
	IDC_HAND        = 32649 // Hand
	IDC_HELP        = 32651 // Arrow and question mark
	IDC_IBEAM       = 32513 // I-beam
	IDC_NO          = 32648 // Slashed circle
	IDC_SIZEALL     = 32646 // Four-pointed arrow pointing north, south, east, and west
	IDC_SIZENESW    = 32643 // Double-pointed arrow pointing northeast and southwest
	IDC_SIZENS      = 32645 // Double-pointed arrow pointing north and south
	IDC_SIZENWSE    = 32642 // Double-pointed arrow pointing northwest and southeast
	IDC_SIZEWE      = 32644 // Double-pointed arrow pointing west and east
	IDC_UPARROW     = 32516 // Vertical arrow
	IDC_WAIT        = 32514 // Hour

	INFINITE = 0xFFFFFFFF

	LOGPIXELSX = 88

	MDT_EFFECTIVE_DPI = 0

	MONITOR_DEFAULTTOPRIMARY = 1

	NI_COMPOSITIONSTR = 0x0015

	SIZE_MAXIMIZED = 2
	SIZE_MINIMIZED = 1
	SIZE_RESTORED  = 0

	SCS_SETSTR = GCS_COMPREADSTR | GCS_COMPSTR

	SM_CXSIZEFRAME = 32
	SM_CYSIZEFRAME = 33

	SW_SHOWDEFAULT   = 10
	SW_SHOWMINIMIZED = 2
	SW_SHOWMAXIMIZED = 3
	SW_SHOWNORMAL    = 1
	SW_SHOW          = 5

	SWP_NOSIZE        = 0x0001
	SWP_NOMOVE        = 0x0002
	SWP_NOZORDER      = 0x0004
	SWP_FRAMECHANGED  = 0x0020
	SWP_SHOWWINDOW    = 0x0040
	SWP_HIDEWINDOW    = 0x0080
	SWP_NOOWNERZORDER = 0x0200

	USER_TIMER_MINIMUM = 0x0000000A

	UNICODE_NOCHAR = 65535

	WS_BORDER = 0x00800000
	WS_CHILD  = 0x40000000

	WS_CLIPCHILDREN     = 0x02000000
	WS_CLIPSIBLINGS     = 0x04000000
	WS_MAXIMIZE         = 0x01000000
	WS_ICONIC           = 0x20000000
	WS_VISIBLE          = 0x10000000
	WS_OVERLAPPED       = 0x00000000
	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME |
		WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_CAPTION       = 0x00C00000
	WS_SYSMENU       = 0x00080000
	WS_THICKFRAME    = 0x00040000
	WS_MINIMIZEBOX   = 0x00020000
	WS_MAXIMIZEBOX   = 0x00010000
	WS_SIZEBOX       = 0x00040000
	WS_EX_WINDOWEDGE = 0x00000100
	WS_DLGFRAME      = 0x00400000
	WS_POPUP         = 0x80000000

	WS_EX_APPWINDOW = 0x00040000

	QS_ALLINPUT = 0x04FF

	MWMO_WAITALL        = 0x0001
	MWMO_INPUTAVAILABLE = 0x0004

	WAIT_OBJECT_0 = 0

	PM_REMOVE   = 0x0001
	PM_NOREMOVE = 0x0000

	GHND = 0x0042

	CF_UNICODETEXT = 13
	IMAGE_BITMAP   = 0
	IMAGE_ICON     = 1
	IMAGE_CURSOR   = 2

	LR_CREATEDIBSECTION = 0x00002000
	LR_DEFAULTCOLOR     = 0x00000000
	LR_DEFAULTSIZE      = 0x00000040
	LR_LOADFROMFILE     = 0x00000010
	LR_LOADMAP3DCOLORS  = 0x00001000
	LR_LOADTRANSPARENT  = 0x00000020
	LR_MONOCHROME       = 0x00000001
	LR_SHARED           = 0x00008000
	LR_VGACOLOR         = 0x00000080
)

type VS_FIXEDFILEINFO struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}
