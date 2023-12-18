//go:build linux
// +build linux

package winapi

const ID_BUTTON_1 = 101
const ID_BUTTON_2 = 102

const HWND_TOPMOST = -1
const SWP_NOMOVE = 2

type Window struct {
	Hwnd      uintptr
	Childrens map[int]*Window
	Config    Config
}

func GetKeyState(key int32) int16 {
	return 0
}

func CreateNativeMainWindow(config Config) (*Window, error) {
	return nil, nil
}

func CreateLabel(win *Window, config Config) (*Window, error) {
	return nil, nil
}

func SendMessage(hwnd uintptr, m uint32, wParam, lParam uint32) int32 {

	return 0
}

func SetWindowPos(hwnd uintptr,
	HWND_TOPMOST,
	x, y, w, h, move int32,
) {

}

func Loop() {

}
func GetFileVersion() string {
	return ""
}
