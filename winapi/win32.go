// SPDX-License-Identifier: Unlicense OR MIT

//go:build windows
// +build windows

package winapi

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	syscall "golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var (
	kernel32          = syscall.NewLazySystemDLL("kernel32.dll")
	_GetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	_GlobalAlloc      = kernel32.NewProc("GlobalAlloc")
	_GlobalFree       = kernel32.NewProc("GlobalFree")
	_GlobalLock       = kernel32.NewProc("GlobalLock")
	_GlobalUnlock     = kernel32.NewProc("GlobalUnlock")

	user32                       = syscall.NewLazySystemDLL("user32.dll")
	_AdjustWindowRectEx          = user32.NewProc("AdjustWindowRectEx")
	_CallMsgFilter               = user32.NewProc("CallMsgFilterW")
	_CloseClipboard              = user32.NewProc("CloseClipboard")
	_CreateWindowEx              = user32.NewProc("CreateWindowExW")
	_DefWindowProc               = user32.NewProc("DefWindowProcW")
	_DestroyWindow               = user32.NewProc("DestroyWindow")
	_DispatchMessage             = user32.NewProc("DispatchMessageW")
	_EmptyClipboard              = user32.NewProc("EmptyClipboard")
	_FillRect                    = user32.NewProc("FillRect")
	_GetWindowRect               = user32.NewProc("GetWindowRect")
	_GetClientRect               = user32.NewProc("GetClientRect")
	_GetClipboardData            = user32.NewProc("GetClipboardData")
	_GetDC                       = user32.NewProc("GetDC")
	_GetDpiForWindow             = user32.NewProc("GetDpiForWindow")
	_GetKeyState                 = user32.NewProc("GetKeyState")
	_GetMessage                  = user32.NewProc("GetMessageW")
	_SendMessage                 = user32.NewProc("SendMessageW")
	_GetMessageTime              = user32.NewProc("GetMessageTime")
	_GetMonitorInfo              = user32.NewProc("GetMonitorInfoW")
	_GetSystemMetrics            = user32.NewProc("GetSystemMetrics")
	_GetWindowLong               = user32.NewProc("GetWindowLongPtrW")
	_GetWindowLong32             = user32.NewProc("GetWindowLongW")
	_GetWindowPlacement          = user32.NewProc("GetWindowPlacement")
	_InvalidateRect              = user32.NewProc("InvalidateRect")
	_KillTimer                   = user32.NewProc("KillTimer")
	_LoadCursor                  = user32.NewProc("LoadCursorW")
	_LoadImage                   = user32.NewProc("LoadImageW")
	_MonitorFromPoint            = user32.NewProc("MonitorFromPoint")
	_MonitorFromWindow           = user32.NewProc("MonitorFromWindow")
	_MoveWindow                  = user32.NewProc("MoveWindow")
	_MsgWaitForMultipleObjectsEx = user32.NewProc("MsgWaitForMultipleObjectsEx")
	_OpenClipboard               = user32.NewProc("OpenClipboard")
	_PeekMessage                 = user32.NewProc("PeekMessageW")
	_PostMessage                 = user32.NewProc("PostMessageW")
	_PostQuitMessage             = user32.NewProc("PostQuitMessage")
	_ReleaseCapture              = user32.NewProc("ReleaseCapture")
	_RegisterClassExW            = user32.NewProc("RegisterClassExW")
	_ReleaseDC                   = user32.NewProc("ReleaseDC")
	_ScreenToClient              = user32.NewProc("ScreenToClient")
	_ShowWindow                  = user32.NewProc("ShowWindow")
	_SetCapture                  = user32.NewProc("SetCapture")
	_SetCursor                   = user32.NewProc("SetCursor")
	_SetClipboardData            = user32.NewProc("SetClipboardData")
	_SetForegroundWindow         = user32.NewProc("SetForegroundWindow")
	_SetFocus                    = user32.NewProc("SetFocus")
	_SetProcessDPIAware          = user32.NewProc("SetProcessDPIAware")
	_SetTimer                    = user32.NewProc("SetTimer")
	_SetWindowLong               = user32.NewProc("SetWindowLongPtrW")
	_SetWindowLong32             = user32.NewProc("SetWindowLongW")
	_SetWindowPlacement          = user32.NewProc("SetWindowPlacement")
	_SetWindowPos                = user32.NewProc("SetWindowPos")
	_SetWindowText               = user32.NewProc("SetWindowTextW")
	_TranslateMessage            = user32.NewProc("TranslateMessage")
	_UnregisterClass             = user32.NewProc("UnregisterClassW")
	_UpdateWindow                = user32.NewProc("UpdateWindow")
	_EnableWindow                = user32.NewProc("EnableWindow")
	_BeginPaint                  = user32.NewProc("BeginPaint")
	_EndPaint                    = user32.NewProc("EndPaint")

	shcore            = syscall.NewLazySystemDLL("shcore")
	_GetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")

	gdi32             = syscall.NewLazySystemDLL("gdi32")
	_GetDeviceCaps    = gdi32.NewProc("GetDeviceCaps")
	_SetBkColor       = gdi32.NewProc("SetBkColor")
	_SetBkMode        = gdi32.NewProc("SetBkMode")
	_BeginPath        = gdi32.NewProc("BeginPath")
	_EndPath          = gdi32.NewProc("EndPath")
	_TextOut          = gdi32.NewProc("TextOutW")
	_SetTextColor     = gdi32.NewProc("SetTextColor")
	_GetStockObject   = gdi32.NewProc("GetStockObject")
	_CreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	_SelectObject     = gdi32.NewProc("SelectObject")
	_SetTextAlign     = gdi32.NewProc("SetTextAlign")
	_CreateFont       = gdi32.NewProc("CreateFontW")

	imm32                    = syscall.NewLazySystemDLL("imm32")
	_ImmGetContext           = imm32.NewProc("ImmGetContext")
	_ImmGetCompositionString = imm32.NewProc("ImmGetCompositionStringW")
	_ImmNotifyIME            = imm32.NewProc("ImmNotifyIME")
	_ImmReleaseContext       = imm32.NewProc("ImmReleaseContext")
	_ImmSetCandidateWindow   = imm32.NewProc("ImmSetCandidateWindow")
	_ImmSetCompositionWindow = imm32.NewProc("ImmSetCompositionWindow")

	dwmapi                        = syscall.NewLazySystemDLL("dwmapi")
	_DwmExtendFrameIntoClientArea = dwmapi.NewProc("DwmExtendFrameIntoClientArea")

	version                = syscall.NewLazyDLL("version.dll")
	getFileVersionInfoSize = version.NewProc("GetFileVersionInfoSizeW")
	getFileVersionInfo     = version.NewProc("GetFileVersionInfoW")
	verQueryValue          = version.NewProc("VerQueryValueW")
)

func AdjustWindowRectEx(r *Rect, dwStyle uint32, bMenu int, dwExStyle uint32) {
	_AdjustWindowRectEx.Call(uintptr(unsafe.Pointer(r)), uintptr(dwStyle), uintptr(bMenu), uintptr(dwExStyle))
}

func CallMsgFilter(m *Msg, nCode uintptr) bool {
	r, _, _ := _CallMsgFilter.Call(uintptr(unsafe.Pointer(m)), nCode)
	return r != 0
}

func CloseClipboard() error {
	r, _, err := _CloseClipboard.Call()
	if r == 0 {
		return fmt.Errorf("CloseClipboard: %v", err)
	}
	return nil
}

func CreateWindowEx(dwExStyle uint32, lpClassName string, lpWindowName string, dwStyle uint32, x, y, w, h int32, hWndParent, hMenu, hInstance syscall.Handle, lpParam uintptr) (syscall.Handle, error) {
	cname := syscall.StringToUTF16Ptr(lpClassName)
	wname := syscall.StringToUTF16Ptr(lpWindowName)
	hwnd, _, err := _CreateWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(unsafe.Pointer(cname)),
		uintptr(unsafe.Pointer(wname)),
		uintptr(dwStyle),
		uintptr(x), uintptr(y),
		uintptr(w), uintptr(h),
		uintptr(hWndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(lpParam))
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowEx failed: %v", err)
	}
	return syscall.Handle(hwnd), nil
}

func DefWindowProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) int {
	r, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
	return int(r)
}

func DestroyWindow(hwnd syscall.Handle) {
	_DestroyWindow.Call(uintptr(hwnd))
}

func DispatchMessage(m *Msg) {
	_DispatchMessage.Call(uintptr(unsafe.Pointer(m)))
}

func DwmExtendFrameIntoClientArea(hwnd syscall.Handle, margins Margins) error {
	r, _, _ := _DwmExtendFrameIntoClientArea.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&margins)))
	if r != 0 {
		return fmt.Errorf("DwmExtendFrameIntoClientArea: %#x", r)
	}
	return nil
}

func EmptyClipboard() error {
	r, _, err := _EmptyClipboard.Call()
	if r == 0 {
		return fmt.Errorf("EmptyClipboard: %v", err)
	}
	return nil
}

func GetWindowRect(hwnd syscall.Handle) Rect {
	var r Rect
	_GetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r
}

func GetClientRect(hwnd syscall.Handle) Rect {
	var r Rect
	_GetClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r
}

func GetClipboardData(format uint32) (syscall.Handle, error) {
	r, _, err := _GetClipboardData.Call(uintptr(format))
	if r == 0 {
		return 0, fmt.Errorf("GetClipboardData: %v", err)
	}
	return syscall.Handle(r), nil
}

func GetDC(hwnd syscall.Handle) (syscall.Handle, error) {
	hdc, _, err := _GetDC.Call(uintptr(hwnd))
	if hdc == 0 {
		return 0, fmt.Errorf("GetDC failed: %v", err)
	}
	return syscall.Handle(hdc), nil
}

func CreateSolidBrush(color int32) (syscall.Handle, error) {
	hbr, _, err := _CreateSolidBrush.Call(uintptr(color))
	if hbr == 0 {
		return 0, fmt.Errorf("CreateSolidBrush failed: %v", err)
	}
	return syscall.Handle(hbr), nil
}
func GetModuleHandle() (syscall.Handle, error) {
	h, _, err := _GetModuleHandleW.Call(uintptr(0))
	if h == 0 {
		return 0, fmt.Errorf("GetModuleHandleW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func getDeviceCaps(hdc syscall.Handle, index int32) int {
	c, _, _ := _GetDeviceCaps.Call(uintptr(hdc), uintptr(index))
	return int(c)
}

func getDpiForMonitor(hmonitor syscall.Handle, dpiType uint32) int {
	var dpiX, dpiY uintptr
	_GetDpiForMonitor.Call(uintptr(hmonitor), uintptr(dpiType), uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
	return int(dpiX)
}

// GetSystemDPI returns the effective DPI of the system.
func GetSystemDPI() int {
	// Check for GetDpiForMonitor, introduced in Windows 8.1.
	if _GetDpiForMonitor.Find() == nil {
		hmon := monitorFromPoint(Point{}, MONITOR_DEFAULTTOPRIMARY)
		return getDpiForMonitor(hmon, MDT_EFFECTIVE_DPI)
	} else {
		// Fall back to the physical device DPI.
		screenDC, err := GetDC(0)
		if err != nil {
			return 96
		}
		defer ReleaseDC(screenDC)
		return getDeviceCaps(screenDC, LOGPIXELSX)
	}
}

func GetKeyState(nVirtKey int32) int16 {
	c, _, _ := _GetKeyState.Call(uintptr(nVirtKey))
	return int16(c)
}

func GetMessage(m *Msg, hwnd syscall.Handle, wMsgFilterMin, wMsgFilterMax uint32) int32 {
	r, _, _ := _GetMessage.Call(uintptr(unsafe.Pointer(m)),
		uintptr(hwnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax))
	return int32(r)
}
func SendMessage(hwnd syscall.Handle, m uint32, wParam, lParam uint32) int32 {
	r, _, _ := _SendMessage.Call(uintptr(hwnd),
		uintptr(m),
		uintptr(wParam),
		uintptr(lParam))
	return int32(r)
}

func GetMessageTime() time.Duration {
	r, _, _ := _GetMessageTime.Call()
	return time.Duration(r) * time.Millisecond
}

func GetSystemMetrics(nIndex int) int {
	r, _, _ := _GetSystemMetrics.Call(uintptr(nIndex))
	return int(r)
}

// GetWindowDPI returns the effective DPI of the window.
func GetWindowDPI(hwnd syscall.Handle) int {
	// Check for GetDpiForWindow, introduced in Windows 10.
	if _GetDpiForWindow.Find() == nil {
		dpi, _, _ := _GetDpiForWindow.Call(uintptr(hwnd))
		return int(dpi)
	} else {
		return GetSystemDPI()
	}
}

func GetWindowPlacement(hwnd syscall.Handle) *WindowPlacement {
	var wp WindowPlacement
	wp.length = uint32(unsafe.Sizeof(wp))
	_GetWindowPlacement.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&wp)))
	return &wp
}

func InvalidateRect(hwnd syscall.Handle, r *Rect, erase int32) {
	_InvalidateRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(r)), uintptr(erase))
}

func GetMonitorInfo(hwnd syscall.Handle) MonitorInfo {
	var mi MonitorInfo
	mi.cbSize = uint32(unsafe.Sizeof(mi))
	v, _, _ := _MonitorFromWindow.Call(uintptr(hwnd), MONITOR_DEFAULTTOPRIMARY)
	_GetMonitorInfo.Call(v, uintptr(unsafe.Pointer(&mi)))
	return mi
}

func GetWindowLong(hwnd syscall.Handle, index uintptr) (val uintptr) {
	if runtime.GOARCH == "386" {
		val, _, _ = _GetWindowLong32.Call(uintptr(hwnd), index)
	} else {
		val, _, _ = _GetWindowLong.Call(uintptr(hwnd), index)
	}
	return
}

func ImmGetContext(hwnd syscall.Handle) syscall.Handle {
	h, _, _ := _ImmGetContext.Call(uintptr(hwnd))
	return syscall.Handle(h)
}

func ImmReleaseContext(hwnd, imc syscall.Handle) {
	_ImmReleaseContext.Call(uintptr(hwnd), uintptr(imc))
}

func ImmNotifyIME(imc syscall.Handle, action, index, value int) {
	_ImmNotifyIME.Call(uintptr(imc), uintptr(action), uintptr(index), uintptr(value))
}

func ImmGetCompositionString(imc syscall.Handle, key int) string {
	size, _, _ := _ImmGetCompositionString.Call(uintptr(imc), uintptr(key), 0, 0)
	if int32(size) <= 0 {
		return ""
	}
	u16 := make([]uint16, size/unsafe.Sizeof(uint16(0)))
	_ImmGetCompositionString.Call(uintptr(imc), uintptr(key), uintptr(unsafe.Pointer(&u16[0])), size)
	return string(utf16.Decode(u16))
}

func ImmGetCompositionValue(imc syscall.Handle, key int) int {
	val, _, _ := _ImmGetCompositionString.Call(uintptr(imc), uintptr(key), 0, 0)
	return int(int32(val))
}

func ImmSetCompositionWindow(imc syscall.Handle, x, y int) {
	f := CompositionForm{
		dwStyle: CFS_POINT,
		ptCurrentPos: Point{
			X: int32(x), Y: int32(y),
		},
	}
	_ImmSetCompositionWindow.Call(uintptr(imc), uintptr(unsafe.Pointer(&f)))
}

func ImmSetCandidateWindow(imc syscall.Handle, x, y int) {
	f := CandidateForm{
		dwStyle: CFS_CANDIDATEPOS,
		ptCurrentPos: Point{
			X: int32(x), Y: int32(y),
		},
	}
	_ImmSetCandidateWindow.Call(uintptr(imc), uintptr(unsafe.Pointer(&f)))
}

func SetWindowLong(hwnd syscall.Handle, idx uintptr, style uintptr) {
	if runtime.GOARCH == "386" {
		_SetWindowLong32.Call(uintptr(hwnd), idx, style)
	} else {
		_SetWindowLong.Call(uintptr(hwnd), idx, style)
	}
}

func SetWindowPlacement(hwnd syscall.Handle, wp *WindowPlacement) {
	_SetWindowPlacement.Call(uintptr(hwnd), uintptr(unsafe.Pointer(wp)))
}

func SetWindowPos(hwnd syscall.Handle, hwndInsertAfter uint32, x, y, dx, dy int32, style uintptr) {
	_SetWindowPos.Call(uintptr(hwnd), uintptr(hwndInsertAfter),
		uintptr(x), uintptr(y),
		uintptr(dx), uintptr(dy),
		style,
	)
}

func SetWindowText(hwnd syscall.Handle, title string) {
	defer func() {
		if val := recover(); val != nil {
			SysLog(1, "SetWindowText")
		}
	}()

	wname := syscall.StringToUTF16Ptr(title)
	_SetWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(wname)))
}

func SetBkColor(hdc syscall.Handle, color uint32) {
	_SetBkColor.Call(uintptr(hdc), uintptr(color))
}

func SetBkMode(hdc syscall.Handle, mode uint32) {
	_SetBkMode.Call(uintptr(hdc), uintptr(mode))
}

func SetTextColor(hdc syscall.Handle, color uint32) {
	_SetTextColor.Call(uintptr(hdc), uintptr(color))
}

func BeginPath(hdc syscall.Handle) {
	_BeginPath.Call(uintptr(hdc))
}

func EndPath(hdc syscall.Handle) {
	_EndPath.Call(uintptr(hdc))
}

func TextOut(hdc syscall.Handle, x int32, y int32, text *string, len int32) {

	_text := syscall.StringToUTF16Ptr(*text)
	_TextOut.Call(uintptr(hdc), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(_text)), uintptr(len))
}

func GetStockObject(i int32) syscall.Handle {
	r1, _, _ := _GetStockObject.Call(uintptr(i))
	return syscall.Handle(r1)
}
func FillRect(hdc syscall.Handle, r *Rect, hbr syscall.Handle) {
	_FillRect.Call(uintptr(hdc), uintptr(unsafe.Pointer(r)), uintptr(unsafe.Pointer(hbr)))
}

func GlobalAlloc(size int) (syscall.Handle, error) {
	r, _, err := _GlobalAlloc.Call(GHND, uintptr(size))
	if r == 0 {
		return 0, fmt.Errorf("GlobalAlloc: %v", err)
	}
	return syscall.Handle(r), nil
}

func GlobalFree(h syscall.Handle) {
	_GlobalFree.Call(uintptr(h))
}

func GlobalLock(h syscall.Handle) (unsafe.Pointer, error) {
	r, _, err := _GlobalLock.Call(uintptr(h))
	if r == 0 {
		return nil, fmt.Errorf("GlobalLock: %v", err)
	}
	return unsafe.Pointer(r), nil
}

func GlobalUnlock(h syscall.Handle) {
	_GlobalUnlock.Call(uintptr(h))
}

func KillTimer(hwnd syscall.Handle, nIDEvent uintptr) error {
	r, _, err := _SetTimer.Call(uintptr(hwnd), uintptr(nIDEvent), 0, 0)
	if r == 0 {
		return fmt.Errorf("KillTimer failed: %v", err)
	}
	return nil
}

func LoadCursor(curID uint16) (syscall.Handle, error) {
	h, _, err := _LoadCursor.Call(0, uintptr(curID))
	if h == 0 {
		return 0, fmt.Errorf("LoadCursorW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

// Загрузка иконки из файла
func LoadIconFromFile(fName string) (syscall.Handle, error) {
	defer func() {
		if val := recover(); val != nil {
			SysLog(1, "LoadIconFromFile")
		}
	}()
	var hInst syscall.Handle = 0

	res := unsafe.Pointer(syscall.StringToUTF16Ptr(fName))
	typ := uint32(IMAGE_ICON)
	cx := 0
	cy := 0
	fuload := uint32(LR_DEFAULTSIZE | LR_LOADFROMFILE)
	h, _, err := _LoadImage.Call(uintptr(hInst), uintptr(res), uintptr(typ), uintptr(cx), uintptr(cy), uintptr(fuload))
	if h == 0 {
		return 0, fmt.Errorf("LoadImageW failed: %v", err)
	}
	return syscall.Handle(h), nil

}

// Загрузка картинки из ресурса
func LoadImage(hInst syscall.Handle, res uint32, typ uint32, cx, cy int, fuload uint32) (syscall.Handle, error) {

	h, _, err := _LoadImage.Call(uintptr(hInst), uintptr(res), uintptr(typ), uintptr(cx), uintptr(cy), uintptr(fuload))
	if h == 0 {
		return 0, fmt.Errorf("LoadImageW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func MoveWindow(hwnd syscall.Handle, x, y, width, height int32, repaint bool) {
	var paint uintptr
	if repaint {
		paint = TRUE
	}
	_MoveWindow.Call(uintptr(hwnd), uintptr(x), uintptr(y), uintptr(width), uintptr(height), paint)
}

func monitorFromPoint(pt Point, flags uint32) syscall.Handle {
	r, _, _ := _MonitorFromPoint.Call(uintptr(pt.X), uintptr(pt.Y), uintptr(flags))
	return syscall.Handle(r)
}

func MsgWaitForMultipleObjectsEx(nCount uint32, pHandles uintptr, millis, mask, flags uint32) (uint32, error) {
	r, _, err := _MsgWaitForMultipleObjectsEx.Call(uintptr(nCount), pHandles, uintptr(millis), uintptr(mask), uintptr(flags))
	res := uint32(r)
	if res == 0xFFFFFFFF {
		return 0, fmt.Errorf("MsgWaitForMultipleObjectsEx failed: %v", err)
	}
	return res, nil
}

func OpenClipboard(hwnd syscall.Handle) error {
	r, _, err := _OpenClipboard.Call(uintptr(hwnd))
	if r == 0 {
		return fmt.Errorf("OpenClipboard: %v", err)
	}
	return nil
}

func PeekMessage(m *Msg, hwnd syscall.Handle, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	r, _, _ := _PeekMessage.Call(uintptr(unsafe.Pointer(m)), uintptr(hwnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax), uintptr(wRemoveMsg))
	return r != 0
}

func PostQuitMessage(exitCode uintptr) {
	_PostQuitMessage.Call(exitCode)
}

func PostMessage(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) error {
	r, _, err := _PostMessage.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	if r == 0 {
		return fmt.Errorf("PostMessage failed: %v", err)
	}
	return nil
}

func ReleaseCapture() bool {
	r, _, _ := _ReleaseCapture.Call()
	return r != 0
}

func RegisterClassEx(cls *WndClassEx) (uint16, error) {
	a, _, err := _RegisterClassExW.Call(uintptr(unsafe.Pointer(cls)))
	if a == 0 {
		return 0, fmt.Errorf("RegisterClassExW failed: %v", err)
	}
	return uint16(a), nil
}

func ReleaseDC(hdc syscall.Handle) {
	_ReleaseDC.Call(uintptr(hdc))
}

func SetForegroundWindow(hwnd syscall.Handle) {
	_SetForegroundWindow.Call(uintptr(hwnd))
}

func SetFocus(hwnd syscall.Handle) {
	_SetFocus.Call(uintptr(hwnd))
}

func SetProcessDPIAware() {
	_SetProcessDPIAware.Call()
}

func SetCapture(hwnd syscall.Handle) syscall.Handle {
	r, _, _ := _SetCapture.Call(uintptr(hwnd))
	return syscall.Handle(r)
}

func SetClipboardData(format uint32, mem syscall.Handle) error {
	r, _, err := _SetClipboardData.Call(uintptr(format), uintptr(mem))
	if r == 0 {
		return fmt.Errorf("SetClipboardData: %v", err)
	}
	return nil
}

func SetCursor(h syscall.Handle) {
	_SetCursor.Call(uintptr(h))
}

func SetTimer(hwnd syscall.Handle, nIDEvent uintptr, uElapse uint32, timerProc uintptr) error {
	r, _, err := _SetTimer.Call(uintptr(hwnd), uintptr(nIDEvent), uintptr(uElapse), timerProc)
	if r == 0 {
		return fmt.Errorf("SetTimer failed: %v", err)
	}
	return nil
}

func ScreenToClient(hwnd syscall.Handle, p *Point) {
	_ScreenToClient.Call(uintptr(hwnd), uintptr(unsafe.Pointer(p)))
}

func ShowWindow(hwnd syscall.Handle, nCmdShow int32) {
	_ShowWindow.Call(uintptr(hwnd), uintptr(nCmdShow))
}

func TranslateMessage(m *Msg) {
	_TranslateMessage.Call(uintptr(unsafe.Pointer(m)))
}

func UnregisterClass(cls uint16, hInst syscall.Handle) {
	_UnregisterClass.Call(uintptr(cls), uintptr(hInst))
}

func UpdateWindow(hwnd syscall.Handle) {
	_UpdateWindow.Call(uintptr(hwnd))
}

func EnableWindow(hwnd syscall.Handle, enable int32) {
	_UpdateWindow.Call(uintptr(hwnd), uintptr(enable))
}

func (p WindowPlacement) Rect() Rect {
	return p.rcNormalPosition
}

func (p WindowPlacement) IsMinimized() bool {
	return p.showCmd == SW_SHOWMINIMIZED
}

func (p WindowPlacement) IsMaximized() bool {
	return p.showCmd == SW_SHOWMAXIMIZED
}

func (p *WindowPlacement) Set(Left, Top, Right, Bottom int) {
	p.rcNormalPosition.Left = int32(Left)
	p.rcNormalPosition.Top = int32(Top)
	p.rcNormalPosition.Right = int32(Right)
	p.rcNormalPosition.Bottom = int32(Bottom)
}

func SelectObject(hdc syscall.Handle, hgdiobj syscall.Handle) syscall.Handle {
	ret, _, _ := _SelectObject.Call(uintptr(hdc), uintptr(hgdiobj))

	return syscall.Handle(ret)
}

func SetTextAlign(hdc syscall.Handle, fMode uint32) uint32 {
	ret, _, _ := _SetTextAlign.Call(uintptr(hdc), uintptr(fMode))
	return uint32(ret)
}

func CreateFont(
	nHeight, nWidth,
	nEscapement,
	nOrientation,
	fnWeight int32,
	fdwItalic,
	fdwUnderline,
	fdwStrikeOut,
	fdwCharSet,
	fdwOutputPrecision,
	fdwClipPrecision,
	fdwQuality,
	fdwPitchAndFamily uint32,
	lpszFace *uint16) syscall.Handle {
	ret, _, _ := _CreateFont.Call(
		uintptr(nHeight),
		uintptr(nWidth),
		uintptr(nEscapement),
		uintptr(nOrientation),
		uintptr(fnWeight),
		uintptr(fdwItalic),
		uintptr(fdwUnderline),
		uintptr(fdwStrikeOut),
		uintptr(fdwCharSet),
		uintptr(fdwOutputPrecision),
		uintptr(fdwClipPrecision),
		uintptr(fdwQuality),
		uintptr(fdwPitchAndFamily),
		uintptr(unsafe.Pointer(lpszFace)))

	return syscall.Handle(ret)
}

func BeginPaint(hwnd syscall.Handle, lpPaint *PAINTSTRUCT) syscall.Handle {
	ret, _, _ := _BeginPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(lpPaint)))
	return syscall.Handle(ret)
}

func EndPaint(hwnd syscall.Handle, lpPaint *PAINTSTRUCT) {
	_EndPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(lpPaint)))
}

func Loword(in uint32) uint16 {
	return uint16(in & 0x0000ffff)
}
func Hiword(in uint32) uint16 {
	return uint16((in >> 8) & 0x0000ffff)
}

// FileVersion concatenates FileVersionMS and FileVersionLS to a uint64 value.
func (fi VS_FIXEDFILEINFO) FileVersion() uint64 {
	return uint64(fi.FileVersionMS)<<32 | uint64(fi.FileVersionLS)
}

// FileDate concatenates FileDateMS and FileDateLS to a uint64 value.
func (fi VS_FIXEDFILEINFO) FileDate() uint64 {
	return uint64(fi.FileDateMS)<<32 | uint64(fi.FileDateLS)
}

func GetFileVersionInfoSize(path string) uint32 {
	if len(path) == 0 {
		return 0
	}

	ret, _, _ := getFileVersionInfoSize.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		0,
	)
	return uint32(ret)
}

func GetFileVersionInfo(path string, data []byte) bool {
	if len(path) == 0 {
		return false
	}
	ret, _, _ := getFileVersionInfo.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		0,
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&data[0])),
	)
	return ret != 0
}

// VerQueryValueRoot calls VerQueryValue
// (https://msdn.microsoft.com/en-us/library/windows/desktop/ms647464(v=vs.85).aspx)
// with `\` (root) to retieve the VS_FIXEDFILEINFO.
func VerQueryValueRoot(block []byte) (VS_FIXEDFILEINFO, bool) {
	var offset uintptr
	var length uint
	blockStart := uintptr(unsafe.Pointer(&block[0]))
	ret, _, _ := verQueryValue.Call(
		blockStart,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(`\`))),
		uintptr(unsafe.Pointer(&offset)),
		uintptr(unsafe.Pointer(&length)),
	)
	if ret == 0 {
		return VS_FIXEDFILEINFO{}, false
	}
	start := int(offset) - int(blockStart)
	end := start + int(length)
	if start < 0 || start >= len(block) || end < start || end > len(block) {
		return VS_FIXEDFILEINFO{}, false
	}
	data := block[start:end]
	info := *((*VS_FIXEDFILEINFO)(unsafe.Pointer(&data[0])))
	return info, true
}

// Загрузка иконки в текстовый файл для включения в тело программы
// Иконки должны размещаться в папке img проекта
func loadImg(name string) ([]byte, error) {
	res, err := os.ReadFile(".\\img\\" + name + ".ico")
	if err == nil {
		res2 := name + "=[]byte{"
		for _, b := range res {
			res2 = res2 + fmt.Sprintf("0x%02x", b) + ","
		}
		res2 = res2 + "}"
		res2 = strings.Replace(res2, ",}", "}", 1)
		os.WriteFile(name+".go", []byte(res2), syscall.O_RDWR)
	}
	return res, err
}

func SysLog(level int, msg string) {
	if runtime.GOOS == "windows" {
		var elog debug.Log
		var name string = "CheckServer"
		var err error

		// Чтобы это работало, надо запускать в режиме Администратора
		err = eventlog.InstallAsEventCreate(name, eventlog.Info|eventlog.Warning|eventlog.Error)
		if err != nil {
			log.Println(msg)
		} else {
			defer eventlog.Remove(name)

			elog, err = eventlog.Open(name)
			if err != nil {
				return
			}
			switch level {
			case 1:
				elog.Error(1, fmt.Sprintf("%s", msg))
			default:
				elog.Info(1, fmt.Sprintf("%s", msg))
			}
		}
	}
}
