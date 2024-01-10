//go:build windows
// +build windows

package winapi

import (
	"fmt"
	"image"
	"log"
	"os"
	"unicode"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

func Loop() {
	msg := new(Msg)
	for GetMessage(msg, 0, 0, 0) > 0 {
		TranslateMessage(msg)
		DispatchMessage(msg)
	}
}

// Основной обработчик событий главного окна
func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) int {
	win, exists := WinMap.Load(hwnd)
	if !exists {
		// Эти сообщения появляются еще до создания окна, поэтому его хэндла нет в WinMap!!!
		/*
			Я не собираюсь использовать этии сообщения, поэтому игнорирую
			if msg == WM_NCCREATE { // сначала приходит это
				// return 1 // Если вернуть 0 - окно не создается, если 1 - нет title
				// если не обрабатываем, то передаем обработку функции DefWindowProc
			} else if msg == WM_CREATE { // потом это
				return 0 // Если вернуть -1 - окно не создается
			}
		*/
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	// Далее работа с сообщениями от уже созданного окна
	w := win.(*Window)

	switch msg {
	case WM_DESTROY:
		if w.Hdc != 0 {
			ReleaseDC(w.Hdc)
			w.Hdc = 0
		}
		w.Hwnd = 0
		PostQuitMessage(0)

	case WM_UNICHAR:
		if wParam == UNICODE_NOCHAR {
			return TRUE
		}
		fallthrough
	case WM_CHAR:
		if r := rune(wParam); unicode.IsPrint(r) {
			//			w.w.EditorInsert(string(r))
		}
		// The message is processed.
		return TRUE
	case WM_DPICHANGED:
		// Let Windows know we're prepared for runtime DPI changes.
		return TRUE
	case WM_ERASEBKGND:
		// Avoid flickering between GPU content and background color.
		return TRUE
	case WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, WM_SYSKEYUP:

		if n, ok := convertKeyCode(wParam); ok {
			e := Event{
				Name:      n,
				Modifiers: getModifiers(),
				State:     Press,
			}
			if msg == WM_KEYUP || msg == WM_SYSKEYUP {
				e.State = Release
			}

			w.Config.EventChan <- (e)

			if (wParam == VK_F10) && (msg == WM_SYSKEYDOWN || msg == WM_SYSKEYUP) {
				// Reserve F10 for ourselves, and don't let it open the system menu. Other Windows programs
				// such as cmd.exe and graphical debuggers also reserve F10.
				return 0
			}
		}
	case WM_LBUTTONDOWN:
		w.pointerButton(ButtonPrimary, true, lParam, getModifiers())
	case WM_LBUTTONUP:
		w.pointerButton(ButtonPrimary, false, lParam, getModifiers())
	case WM_RBUTTONDOWN:
		w.pointerButton(ButtonSecondary, true, lParam, getModifiers())
	case WM_RBUTTONUP:
		w.pointerButton(ButtonSecondary, false, lParam, getModifiers())
	case WM_MBUTTONDOWN:
		w.pointerButton(ButtonTertiary, true, lParam, getModifiers())
	case WM_MBUTTONUP:
		w.pointerButton(ButtonTertiary, false, lParam, getModifiers())
	case WM_CANCELMODE:
		log.Println("Cancel")
		// Если обрабатываем, вернуть 0
		// При отправке сообщения WM_CANCELMODE функция DefWindowProc отменяет внутреннюю обработку
		// стандартных входных данных полосы прокрутки,
		// отменяет обработку внутреннего меню и освобождает захват мыши.
		/*
			w.Config.EventChan <- Event{
				SWin:   w,
				Kind:   Cancel,
				Source: Frame,
			}
			return 0
		*/
	case WM_SETFOCUS:
		// Это щелчок в окне
		w.Focused = true
		x, y := coordsFromlParam(lParam)
		w.Config.EventChan <- Event{
			SWin:      w,
			Kind:      Enter,
			Source:    Mouse,
			Position:  image.Point{X: x, Y: y},
			Mbuttons:  w.Mbuttons,
			Time:      GetMessageTime(),
			Modifiers: getModifiers(),
		}
	case WM_KILLFOCUS:
		// Щелчок вне нашего главного окна
		// Щелчок по кнопке тоже дает это событие
		w.Focused = false
		w.Config.EventChan <- Event{
			SWin:      w,
			Kind:      Leave,
			Source:    Mouse,
			Position:  image.Point{X: -1, Y: -1},
			Mbuttons:  w.Mbuttons,
			Time:      GetMessageTime(),
			Modifiers: getModifiers(),
		}

	case WM_MOUSEMOVE:
		// Это событие будет, даже если наше окно не в фокусе
		// и может быть даже частично перекрыто другим окном
		x, y := coordsFromlParam(lParam)
		p := image.Point{X: x, Y: y}

		w.Config.EventChan <- Event{
			SWin:      w,
			Kind:      Move,
			Source:    Mouse,
			Position:  p,
			Mbuttons:  w.Mbuttons,
			Time:      GetMessageTime(),
			Modifiers: getModifiers(),
		}

	case WM_MOUSEWHEEL:
		// Поворот колесика +- WHEEL_DELTA (120) HIWORD wParam
		//		w.scrollEvent(wParam, lParam, false, getModifiers())
	case WM_MOUSEHWHEEL:
		// Поворот горизонтального колесика +- WHEEL_DELTA (120) HIWORD wParam
		//		w.scrollEvent(wParam, lParam, true, getModifiers())
	case WM_NCACTIVATE:
		// Отправляется в окно, когда его неклиентская область должна быть изменена,
		// чтобы указать активное или неактивное состояние.
		if w.Stage >= StageInactive {
			if wParam == TRUE {
				w.Stage = StageRunning
			} else {
				w.Stage = StageInactive
			}
		}

	case WM_NCHITTEST:
		// Отправляется в окно, чтобы определить,
		// какая часть окна соответствует определенной экранной координате.
		// Если окно с заголовком, его обрабатывает дефолтная процедура
		if w.Config.SysMenu > 0 {
			break
		}

		x, y := coordsFromlParam(lParam)
		log.Printf("x: %d y: %d", x, y)
		np := Point{X: int32(x), Y: int32(y)}
		ScreenToClient(w.Hwnd, &np)
		log.Printf("np.x: %d np.y: %d", np.X, np.Y)
		area := w.hitTest(int(np.X), int(np.Y))
		log.Printf("area: %d", area)
		return area

	case WM_NCCALCSIZE:
		// Отправляется, когда необходимо вычислить размер и положение клиентской области окна.
		//  Обрабатывая это сообщение,  приложение может управлять содержимым
		// клиентской области окна при изменении размера или положения окна.
		// Если окно с заголовком, его обрабатывает дефолтная процедура
		if w.Config.SysMenu > 0 {
			break
		}

		// No client areas; we draw decorations ourselves.
		if wParam != 1 {
			return 0
		}
		// lParam contains an NCCALCSIZE_PARAMS for us to adjust.
		place := GetWindowPlacement(w.Hwnd)
		if !place.IsMaximized() {
			// Nothing do adjust.
			return 0
		}
		// Adjust window position to avoid the extra padding in maximized
		// state. See https://devblogs.microsoft.com/oldnewthing/20150304-00/?p=44543.
		// Note that trying to do the adjustment in WM_GETMINMAXINFO is ignored by
		szp := (*NCCalcSizeParams)(unsafe.Pointer(uintptr(lParam)))
		mi := GetMonitorInfo(w.Hwnd)
		szp.Rgrc[0] = mi.WorkArea
		return 0

	case WM_PAINT:
		w.draw(true)

	case WM_MOVE:
		w.update()
		w.draw(true)
		return 0

	case WM_SIZE:

		switch wParam {
		case SIZE_MINIMIZED:
			w.Config.Mode = Minimized
			w.Stage = StagePaused
		case SIZE_MAXIMIZED:
			w.Config.Mode = Maximized
			w.Stage = StageRunning
		case SIZE_RESTORED:
			if w.Config.Mode != Fullscreen {
				w.Config.Mode = Windowed
			}
			w.Stage = StageRunning
		}
		InvalidateRect(hwnd, nil, 1)
		UpdateWindow(hwnd)
		return 0
	case WM_GETMINMAXINFO:
		// Отправляется в окно, когда размер или положение окна вот-вот изменится.
		// Приложение может использовать это сообщение, чтобы переопределить
		// развернутый размер и положение окна по умолчанию,
		// а также его минимальный или максимальный размер отслеживания по умолчанию.
		mm := (*MinMaxInfo)(unsafe.Pointer(uintptr(lParam)))
		var bw, bh int32 = 0, 0
		//		if w.Config.Decorated {
		r := GetWindowRect(w.Hwnd)
		cr := GetClientRect(w.Hwnd)
		bw = r.Right - r.Left - (cr.Right - cr.Left)
		bh = r.Bottom - r.Top - (cr.Bottom - cr.Top)
		//		}
		if p := w.Config.MinSize; p.X > 0 || p.Y > 0 {
			mm.PtMinTrackSize = Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		if p := w.Config.MaxSize; p.X > 0 || p.Y > 0 {
			mm.PtMaxTrackSize = Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		return 0
	case WM_SETCURSOR:

		w.CursorIn = (lParam & 0xffff) == HTCLIENT
		if w.CursorIn {
			SetCursor(w.Cursor)
			return TRUE
		}

	case WM_CTLCOLORSTATIC:
		// Установка параметров текста для статических элементов окна (STATIC или ReadOnly EDIT)
		wc := w.Childrens[1]
		log.Println(wc.Hdc, syscall.Handle(wParam))

		SetTextColor(syscall.Handle(wParam), wc.Config.TextColor) // цвет самого теста
		SetBkColor(syscall.Handle(wParam), wc.Config.BgColor)     // цвет подложки текста

		// Если приложение обрабатывает это сообщение, возвращаемое значение представляет собой дескриптор кисти,
		// которую система использует для рисования фона статического элемента управления.
		hbrBkgnd, err := CreateSolidBrush(int32(wc.Config.BgColor)) // цвет заливки окна
		if err == nil {
			return int(hbrBkgnd)
		}
	case WM_COMMAND:
		// Коды команд меню и активных элементов окна (типа кнопки) через присвоенный им код
		log.Printf("WM_COMMAND 0x%08x 0x%08x \n", wParam, lParam)
		// Если мы прописали ID кнопки в качестве hMenu, при создании окна,
		// то в wParam в LOWORD придет этот код
		if Loword(uint32(wParam)) == ID_BUTTON_1 || Loword(uint32(wParam)) == ID_BUTTON_2 {
			// в lParam приходит Handle окна кнопки
			win2, exists := WinMap.Load(syscall.Handle(lParam))
			if exists {
				w2 := win2.(*Window)
				go w.HandleButton(w2, wParam)
				return 0 // если мы обрабатываем, должны вернуть 0
			}
		}
	case WM_NOTIFY:
		log.Printf("WM_NOTIFY 0x%08x 0x%08x \n", wParam, lParam)
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

// ----------------------------------------
func (w *Window) HandleButton(w2 *Window, wParam uintptr) {

	switch Loword(uint32(wParam)) {
	case ID_BUTTON_1:
		log.Println(w2.Config.Title)
		// И какие-то действия

	case ID_BUTTON_2:
		log.Println(w2.Config.Title)
		CloseWindow()
	}

}

// hitTest возвращает область, в которую попал указатель мыши,
// HTCLIENT возвращается при перемещении мыши ынутри клиентской области
// Другие показывают направление относительно клиентской области
// нужно для обработки сообщения  WM_NCHITTEST.
func (w *Window) hitTest(x, y int) int {
	if w.Config.Mode == Fullscreen {
		return HTCLIENT
	}
	if w.Config.Mode != Windowed {
		// Only windowed mode should allow resizing.
		return HTCLIENT
	}
	// Check for resize handle before system actions; otherwise it can be impossible to
	// resize a custom-decorations window when the system move area is flush with the
	// edge of the window.
	top := y <= w.Config.BorderSize.Y
	bottom := y >= w.Config.Size.Y-w.Config.BorderSize.Y
	left := x <= w.Config.BorderSize.X
	right := x >= w.Config.Size.X-w.Config.BorderSize.X
	switch {
	case top && left:
		log.Println("HTTOPLEFT")
		return HTTOPLEFT
	case top && right:
		log.Println("HTTOPRIGHT")
		return HTTOPRIGHT
	case bottom && left:
		log.Println("HTBOTTOMLEFT")
		return HTBOTTOMLEFT
	case bottom && right:
		log.Println("HTBOTTOMRIGHT")
		return HTBOTTOMRIGHT
	case top:
		log.Println("HTTOP")
		return HTTOP
	case bottom:
		log.Println("HTBOTTOM")
		return HTBOTTOM
	case left:
		log.Println("HTLEFT")
		return HTLEFT
	case right:
		log.Println("HTRIGHT")
		return HTRIGHT
	}
	/*
		p := f32.Pt(float32(x), float32(y))

		if a, ok := w.w.ActionAt(p); ok && a == system.ActionMove {
			return  HTCAPTION
		}

	*/

	return HTCLIENT
}

// Перерисовка окна
func (w *Window) draw(sync bool) {
	if w.Config.Size.X == 0 || w.Config.Size.Y == 0 {
		return
	}

	r1 := GetClientRect(w.Hwnd)
	hbrBkgnd, _ := CreateSolidBrush(int32(w.Config.BgColor))
	FillRect(w.Hdc, &r1, hbrBkgnd)

	// Отрисовка текста и фона в статических дочерних окнах
	for _, w2 := range w.Childrens {
		switch w2.Config.Class {
		case "Static":
			w.drawStaticText(w2)
		case "Button":
			InvalidateRect(w2.Hwnd, nil, 1)
			UpdateWindow(w2.Hwnd)
			//		w.drawButton(w2)
		}
	}
}

func (w *Window) drawStaticText(w2 *Window) {
	r1 := GetClientRect(w2.Hwnd)
	hbrBkgnd, _ := CreateSolidBrush(int32(w2.Config.BgColor))
	FillRect(w2.Hdc, &r1, hbrBkgnd)

	var ps PAINTSTRUCT = PAINTSTRUCT{}
	BeginPaint(w2.Hwnd, &ps)
	fontSize := w2.Config.FontSize
	hFont := CreateFont(fontSize, int32(float32(fontSize)*0.4), 0, 0, 0,
		0, 0, 0,
		DEFAULT_CHARSET,
		0, 0, 0, 0,
		syscall.StringToUTF16Ptr("Tahoma"))

	oldFont := SelectObject(w2.Hdc, hFont)
	SetTextColor(w2.Hdc, w2.Config.TextColor) // цвет самого текста
	SetBkColor(w2.Hdc, w2.Config.BgColor)     // цвет подложки текста
	txt := w2.Config.Title
	left := int32(4) // TODO: для кнопок отцентровать
	top := int32(2)
	TextOut(w2.Hdc, left, top, &txt, int32(len(txt)))
	SelectObject(w2.Hdc, oldFont)

	EndPaint(w2.Hwnd, &ps)
}

func (w *Window) drawButton(w2 *Window) {
	r1 := GetClientRect(w2.Hwnd)
	hbrBkgnd, _ := CreateSolidBrush(int32(w2.Config.BgColor))
	FillRect(w2.Hdc, &r1, hbrBkgnd)

	var ps PAINTSTRUCT = PAINTSTRUCT{}
	BeginPaint(w2.Hwnd, &ps)
	fontSize := w2.Config.FontSize
	hFont := CreateFont(fontSize, int32(float32(fontSize)*0.4), 0, 0, 0,
		0, 0, 0,
		DEFAULT_CHARSET,
		0, 0, 0, 0,
		syscall.StringToUTF16Ptr("Tahoma"))

	oldFont := SelectObject(w2.Hdc, hFont)
	SetTextColor(w2.Hdc, w2.Config.TextColor) // цвет самого текста
	SetBkColor(w2.Hdc, w2.Config.BgColor)     // цвет подложки текста
	txt := w2.Config.Title
	left := int32(4) // TODO: для кнопок отцентровать
	top := int32(2)
	TextOut(w2.Hdc, left, top, &txt, int32(len(txt)))
	SelectObject(w2.Hdc, oldFont)

	EndPaint(w2.Hwnd, &ps)
}

// update() handles changes done by the user, and updates the configuration.
// It reads the window style and size/position and updates w.config.
// If anything has changed it emits a ConfigEvent to notify the application.
func (w *Window) update() {

	cr := GetClientRect(w.Hwnd)
	w.Config.Size = image.Point{
		X: int(cr.Right - cr.Left),
		Y: int(cr.Bottom - cr.Top),
	}

	w.Config.BorderSize = image.Pt(
		GetSystemMetrics(SM_CXSIZEFRAME),
		GetSystemMetrics(SM_CYSIZEFRAME),
	)
}

func (w *Window) SetCursor(cursor Cursor) {
	c, err := loadCursor(cursor)
	if err != nil {
		c = resources.cursor
	}
	w.Cursor = c
	SetCursor(w.Cursor) // Win32 API function
}

func loadCursor(cursor Cursor) (syscall.Handle, error) {
	switch cursor {
	case CursorDefault:
		return resources.cursor, nil
	case CursorNone:
		return 0, nil
	default:
		return LoadCursor(windowsCursor[cursor])
	}
}

func (w *Window) pointerButton(btn MButtons, press bool, lParam uintptr, kmods Modifiers) {
	if !w.Focused {
		SetFocus(w.Hwnd)
	}
	log.Println("pointerButton", btn, press)
	var kind Kind
	if press {
		kind = Press
		if w.Mbuttons == 0 {
			SetCapture(w.Hwnd) // Захват событий мыши окном
		}
		w.Mbuttons |= btn
	} else {
		kind = Release
		w.Mbuttons &^= btn
		if w.Mbuttons == 0 {
			ReleaseCapture() // Освобождение событий мыши окном
		}
	}

	x, y := coordsFromlParam(lParam)
	p := image.Point{X: (x), Y: (y)}
	w.Config.EventChan <- Event{
		SWin:      w,
		Kind:      kind,
		Source:    Mouse,
		Position:  p,
		Mbuttons:  w.Mbuttons,
		Time:      GetMessageTime(),
		Modifiers: kmods,
	}

}

func coordsFromlParam(lParam uintptr) (int, int) {
	x := int(int16(lParam & 0xffff))
	y := int(int16((lParam >> 16) & 0xffff))
	return x, y
}

// Текст в статическом окне
func (w *Window) SetText(text string) {
	if w.Parent != nil {
		w.Config.Title = text
		r := GetClientRect(w.Hwnd)
		InvalidateRect(w.Hwnd, &r, 0)
	}
}

func (w *Window) Invalidate() {
	w.update()
	InvalidateRect(w.Hwnd, nil, 1)
	UpdateWindow(w.Hwnd)
	//	w.draw(true)
}

func GetFileVersion() string {
	size := GetFileVersionInfoSize(os.Args[0])
	if size > 0 {
		info := make([]byte, size)
		ok := GetFileVersionInfo(os.Args[0], info)
		if ok {
			fixed, ok := VerQueryValueRoot(info)
			if ok {
				version := fixed.FileVersion()
				VERSION := fmt.Sprintf("v%d.%d.%d",
					version&0xFFFF000000000000>>48,
					version&0x0000FFFF00000000>>32,
					version&0x00000000FFFF0000>>16,
				)
				log.Println("Ver: ", VERSION)
				return VERSION
			}
		}
	}
	return ""
}
