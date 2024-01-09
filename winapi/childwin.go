//go:build windows
// +build windows

package winapi

import "golang.org/x/sys/windows"

// Пользовательские элементы начинаются со 100 (0x64)
const ID_BUTTON_1 = 100
const ID_BUTTON_2 = 101

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
	if config.BorderSize.X > 0 {
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
