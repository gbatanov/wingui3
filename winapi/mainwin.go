//go:build windows
// +build windows

package winapi

import (
	"sync"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

type Window struct {
	Hwnd        syscall.Handle
	Hdc         syscall.Handle
	HInst       syscall.Handle
	Focused     bool
	Stage       Stage
	Config      Config
	Cursor      syscall.Handle
	PointerBtns MButtons //Кнопки мыши
	Parent      *Window
	Childrens   map[int]*Window
	// cursorIn tracks whether the cursor was inside the window according
	// to the most recent WM_SETCURSOR.
	CursorIn bool
	Mbuttons MButtons
}

// iconID это ID в winres.json (#1)
const iconID = 1

var resources struct {
	once sync.Once
	// handle is the module handle from GetModuleHandle.
	handle syscall.Handle
	// cursor is the arrow cursor resource.
	cursor syscall.Handle
}

// initResources initializes the resources global.
func initResources(config Config) error {
	SetProcessDPIAware()
	hInst, err := GetModuleHandle()
	if err != nil {
		return err
	}

	c, err := LoadCursor(IDC_ARROW)
	if err != nil {
		return err
	}

	//	icon, err := LoadIconFromFile(".\\img\\logo.ico") // вариант иконки из файла
	// но лучше брать из предварительно подготовленного ресурса (файл .syso)
	icon, _ := LoadImage(hInst, iconID, IMAGE_ICON, 0, 0, LR_DEFAULTSIZE|LR_SHARED)

	wcls := WndClassEx{
		CbSize:    uint32(unsafe.Sizeof(WndClassEx{})),
		HInstance: hInst,
	}

	wcls.Style = CS_HREDRAW | CS_VREDRAW | CS_OWNDC
	wcls.HIcon = icon
	wcls.LpszClassName = syscall.StringToUTF16Ptr(config.Class)

	wcls.LpfnWndProc = syscall.NewCallback(windowProc)
	_, err = RegisterClassEx(&wcls)
	if err != nil {
		return err
	}
	resources.handle = hInst
	resources.cursor = c

	return nil
}

// Создание основного окна программы
func CreateNativeMainWindow(config Config) (*Window, error) {

	var resErr error
	resources.once.Do(func() {
		resErr = initResources(config)
	})
	if resErr != nil {
		return nil, resErr
	}
	// WS_CAPTION включает в себя WS_BORDER
	var dwExStyle uint32 = 0
	var dwStyle uint32 = 0
	if config.SysMenu == 2 {
		dwStyle = WS_SYSMENU | WS_CAPTION | WS_SIZEBOX
	} else if config.SysMenu == 1 {
		dwStyle = WS_CAPTION | WS_SIZEBOX
	} else {
		dwStyle = WS_POPUP
	}

	if config.Position.X < 0 {
		mi := GetMonitorInfo(0)
		config.Position.X = int(mi.WorkArea.Right) + config.Position.X - config.Size.X //+ int(mi.cbSize)
	}
	if config.Position.Y < 0 {
		mi := GetMonitorInfo(0)
		config.Position.Y = int(mi.WorkArea.Bottom) + config.Position.Y - config.Size.Y //+ int(mi.cbSize)
	}

	hwnd, err := CreateWindowEx(
		dwExStyle,
		config.Class,                                       //	resourceMain.class,                                 //lpClassame
		config.Title,                                       // lpWindowName
		dwStyle,                                            //dwStyle
		int32(config.Position.X), int32(config.Position.Y), //x, y
		int32(config.Size.X), int32(config.Size.Y), //w, h
		0,                //hWndParent
		0,                // hMenu
		resources.handle, //hInstance
		0)                // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Hwnd:      hwnd,
		HInst:     resources.handle,
		Config:    config,
		Parent:    nil,
		Childrens: make(map[int]*Window, 0),
	}
	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}

	WinMap.Store(w.Hwnd, w)

	SetForegroundWindow(w.Hwnd)
	SetFocus(w.Hwnd)
	w.SetCursor(CursorDefault)
	ShowWindow(w.Hwnd, SW_SHOWNORMAL)
	return w, nil
}

// Заглушка
func SetIcon() {

}
