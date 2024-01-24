//go:build windows
// +build windows

package winapi

import (
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
	syscall "golang.org/x/sys/windows"
)

type Window struct {
	Hwnd    syscall.Handle
	Hdc     syscall.Handle
	HInst   syscall.Handle
	Focused bool
	Stage   Stage
	Config  Config
	Cursor  syscall.Handle
	Parent  *Window
	// cursorIn tracks whether the cursor was inside the window according
	// to the most recent WM_SETCURSOR.
	CursorIn  bool
	Mbuttons  MButtons //Кнопки мыши
	IsMain    bool
	Childrens *map[int]*Window
}

// iconID это ID в winres.json (#1)
const iconID = 1

var Wind *Window
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
	child := make(map[int]*Window, 0)
	win := &Window{
		Hwnd:      hwnd,
		HInst:     resources.handle,
		Config:    config,
		Parent:    nil,
		Childrens: &child,
		IsMain:    true,
	}
	win.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}

	WinMap.Store(win.Hwnd, win)
	WinMap.Store(0, win) // Основное окно дублируем с нулевым ключчом, чтобы иметь доступ всегда

	SetForegroundWindow(win.Hwnd)
	SetFocus(win.Hwnd)
	win.SetCursor(CursorDefault)
	ShowWindow(win.Hwnd, SW_SHOWNORMAL)
	Wind = win
	return win, nil
}

// Заглушка для совместимости с Линукс
func SetIcon(smenu int) {

}

// Программное закрытие окна (совместимость с Линукс)
func CloseWindow() {
	w, exists := WinMap.Load(0)
	if exists {
		SendMessage(w.(*Window).Hwnd, WM_CLOSE, 0, 0)
	}
}

// Label
func CreateLabel(parent *Window, config Config) (*Window, error) {
	win, err := CreateChildWindow(parent, config)
	return win, err
}

// Button
func CreateButton(parent *Window, config Config) (*Window, error) {
	win, err := CreateChildWindow(parent, config)
	return win, err
}

// Создаем статическое окно
func CreateChildWindow(parent *Window, config Config) (*Window, error) {

	var dwStyle uint32 = WS_CHILD | WS_VISIBLE
	if config.BorderSize > 0 {
		dwStyle |= WS_BORDER
	}

	// Для дочернего окна hMenu указывает идентификатор дочернего окна, целочисленное значение,
	// используемое элементом управления диалогового окна для уведомления родительского элемента управления
	// о событиях. Приложение определяет идентификатор дочернего окна;
	// оно должно быть уникальным для всех дочерних окон с одинаковым родительским окном.
	hMenu := windows.Handle(config.ID)

	hwnd, err := CreateWindowEx(
		0,
		config.Class,                                       // standard static class,
		config.Title,                                       // lpWindowName
		dwStyle,                                            //dwStyle
		int32(config.Position.X), int32(config.Position.Y), //x, y
		int32(config.Size.X), int32(config.Size.Y), //w, h
		parent.Hwnd,  //hWndParent
		hMenu,        // hMenu
		parent.HInst, //hInstance
		0)            // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Hwnd:      hwnd,
		HInst:     parent.HInst,
		Config:    config,
		Parent:    parent,
		Childrens: nil,
		IsMain:    false,
	}
	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}
	w.SetCursor(CursorDefault)
	WinMap.Store(w.Hwnd, w)

	return w, nil
}
