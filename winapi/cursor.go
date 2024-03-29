//go:build windows
// +build windows

package winapi

// windowsCursor contains mapping from  Cursor to an IDC.
var windowsCursor = [...]uint16{
	CursorDefault:                  IDC_ARROW,
	CursorNone:                     0,
	CursorText:                     IDC_IBEAM,
	CursorVerticalText:             IDC_IBEAM,
	CursorPointer:                  IDC_HAND,
	CursorCrosshair:                IDC_CROSS,
	CursorAllScroll:                IDC_SIZEALL,
	CursorColResize:                IDC_SIZEWE,
	CursorRowResize:                IDC_SIZENS,
	CursorGrab:                     IDC_SIZEALL,
	CursorGrabbing:                 IDC_SIZEALL,
	CursorNotAllowed:               IDC_NO,
	CursorWait:                     IDC_WAIT,
	CursorProgress:                 IDC_APPSTARTING,
	CursorNorthWestResize:          IDC_SIZENWSE,
	CursorNorthEastResize:          IDC_SIZENESW,
	CursorSouthWestResize:          IDC_SIZENESW,
	CursorSouthEastResize:          IDC_SIZENWSE,
	CursorNorthSouthResize:         IDC_SIZENS,
	CursorEastWestResize:           IDC_SIZEWE,
	CursorWestResize:               IDC_SIZEWE,
	CursorEastResize:               IDC_SIZEWE,
	CursorNorthResize:              IDC_SIZENS,
	CursorSouthResize:              IDC_SIZENS,
	CursorNorthEastSouthWestResize: IDC_SIZENESW,
	CursorNorthWestSouthEastResize: IDC_SIZENWSE,
}
