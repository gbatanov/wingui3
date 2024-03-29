//go:build linux
// +build linux

package winapi

import (
	"fmt"
	"image"
	"log"
	"log/syslog"
	"strconv"
	"strings"
	"time"

	"github.com/gbatanov/wingui3/img"
	"github.com/jezek/xgb/xproto"
)

// Основной цикл обработки событий
func Loop() {
	for {
		ev, xerr := X.WaitForEvent()
		if ev == nil && xerr == nil {
			// Возникает при закрытии окна по крестику
			log.Println("Window closed. Exiting...")
			return
		}

		if xerr != nil {
			log.Printf("Error: %s\n", xerr.Error())
		}
		//		log.Println("Event", ev)

		switch ev := ev.(type) {
		case xproto.CreateNotifyEvent:
			//			log.Println("CreateNotifyEvent", ev)
		case xproto.UnmapNotifyEvent: // Сворачивание окна
			//			log.Println("UnmapNotifyEvent", ev)
			w := getWindow(ev.Event)
			SetWindowPos(ev.Event, HWND_TOPMOST,
				int32(w.Config.Position.X-5),
				int32(w.Config.Position.Y-27),
				int32(w.Config.Size.X),
				int32(w.Config.Size.Y), 2)
			xproto.MapWindowChecked(X, ev.Event).Check()

			//		case xproto.LeaveNotifyEvent: // Потеря фокуса
			//			log.Println("LeaveNotifyEvent", ev)

			//Клавиши клавиатуры
		case xproto.KeyPressEvent:
			w := getWindow(ev.Event)
			w.createKbEvent("Press", ev.Detail, ev.Time)
		case xproto.KeyReleaseEvent:
			w := getWindow(ev.Event)
			w.createKbEvent("Release", ev.Detail, ev.Time)

			// Кнопки мыши
		case xproto.ButtonPressEvent:
			w := getWindow(ev.Event)
			w.createMouseEvent("Press", ev.Detail, ev.EventX, ev.EventY, ev.Time)
		case xproto.ButtonReleaseEvent:
			w := getWindow(ev.Event)
			w.createMouseEvent("Release", ev.Detail, ev.EventX, ev.EventY, ev.Time)

		case xproto.MotionNotifyEvent:
			//			log.Printf("%d", ev.Time)
			Wind.Config.EventChan <- Event{
				SWin:      getWindow(ev.Event),
				Kind:      Move,
				Source:    Mouse,
				Position:  image.Point{int(ev.EventX), int(ev.EventY)},
				Mbuttons:  Wind.Mbuttons, //uint8
				Time:      time.Duration(time.Duration(ev.Time)),
				Modifiers: getModifiers(),
			}

		case xproto.ReparentNotifyEvent:
			//			log.Println("Reparent notify ", ev)

		case xproto.ConfigureNotifyEvent:

			//			log.Println("Configure notify ", ev)
			w := getWindow(ev.Event)

			// A window's size, position, border, and/or stacking order is reconfigured by calling XConfigureWindow().
			// The window's position in the stacking order is changed by calling XLowerWindow(), XRaiseWindow(), or XRestackWindows().
			// A window is moved by calling XMoveWindow().
			// A window's size is changed by calling XResizeWindow().
			// A window's size and location is changed by calling XMoveResizeWindow().
			// A window is mapped and its position in the stacking order is changed by calling XMapRaised().
			// A window's border width is changed by calling XSetWindowBorderWidth().
			if w == Wind {
				if ev.Width > uint16(w.Config.MaxSize.X) ||
					ev.Height > uint16(w.Config.MaxSize.Y) {
					SetWindowPos(ev.Event, HWND_TOPMOST,
						int32(w.Config.Position.X-1),
						int32(w.Config.Position.Y-7),
						int32(w.Config.MaxSize.X),
						int32(w.Config.MaxSize.Y), 2)

				} else if ev.Width < uint16(w.Config.MinSize.X) ||
					ev.Height < uint16(w.Config.MinSize.Y) {
					SetWindowPos(ev.Event, HWND_TOPMOST,
						int32(w.Config.Position.X-1),
						int32(w.Config.Position.Y-7),
						int32(w.Config.MinSize.X),
						int32(w.Config.MinSize.Y), 2)

				} else {
					w.Config.Position.X = int(ev.X)
					w.Config.Position.Y = int(ev.Y)
				}
			} else {
				w.Config.Position.X = int(ev.X)
				w.Config.Position.Y = int(ev.Y)

			}

		case xproto.MapNotifyEvent: // Отображение окна, включая дочерние
			//			log.Println("Map notify ", ev)

		case xproto.ResizeRequestEvent: // WM_SIZE Работает криво, отключил в маске
			log.Println("Resize Request ", ev)

		case xproto.ClientMessageEvent:
			log.Println("ClientMessage Event ", ev)

		case xproto.ExposeEvent: // аналог WM_PAINT в Windows
			// log.Println("ExposeEvent notify ", ev)
			if ev.Count == 0 {
				w := getWindow(ev.Window)
				w.draw()
			}

		case xproto.DestroyNotifyEvent:
			// На закрытие по крестику не приходит
			// Событие приходит для каждого окна (главное и дочерние)
			// Для обработки будем отправлять событие только для главного окна
			if ev.Window == Wind.Hwnd {
				Wind.Config.EventChan <- Event{
					SWin:   Wind,
					Kind:   Destroy,
					Source: Frame,
				}

				// Небольшая задержка, чтобы основной поток принял сообщение
				// (после завершения цикла канал закрывается)
				time.Sleep(1 * time.Second)
				return // Завершаем цикл
			}
		} // switch
	} //for
} //Loop

func getWindow(wev xproto.Window) *Window {
	w := Wind
	if wev != Wind.Hwnd {
		wind, exists := WinMap.Load(wev)
		if exists {
			w = wind.(*Window)
		}
	}
	return w
}

// в линукс приходят скан-коды, переводим в код символа на клиентской стороне
func (w *Window) createKbEvent(evType string, btn xproto.Keycode, evTime xproto.Timestamp) {
	//	log.Printf("btn: %d 0x%04x\n", btn, btn) // "A" 38 0x26
	mod := getModifiers()
	//	log.Printf("mod before: 0x%04x\n", mod)
	var keyCode xproto.Keysym = 0

	keycodeIndx := (int(btn) - int(Wind.FirstCode)) * int(Wind.KeysymsPerKeycode)
	keyCode = Wind.Keymap[keycodeIndx]
	if keyCode < 255 { // "нормальные символы"
		if mod&ModShift != 0 {
			keycodeIndx++
			keyCode = Wind.Keymap[keycodeIndx]
		}
	} else {
		SetKeyState(keyCode, evType == "Press")
		return
	}

	evnt := Event{
		SWin:      w,
		Source:    Keyboard,
		Position:  image.Point{0, 0},
		Mbuttons:  w.Mbuttons, //uint8
		Time:      time.Duration(time.Duration(evTime)),
		Modifiers: mod,
		Name:      "",
		Keycode:   xproto.Keycode(keyCode),
	}
	if evType == "Press" {
		evnt.Kind = Press
	} else if evType == "Release" {
		evnt.Kind = Release
	}
	if n, ok := convertKeyCode(uintptr(keyCode)); ok {
		evnt.Name = n
	}

	//	log.Println("evnt ", evnt)
	Wind.Config.EventChan <- evnt
}

func (w *Window) WinTranslateCoordinates(x int, y int) (int, int, error) {

	if w != Wind {
		tc := xproto.TranslateCoordinates(X, w.Hwnd, w.Parent, int16(x), int16(y))
		tcR, err := tc.Reply()
		if err != nil {
			log.Println(err.Error())
			return y, y, err
		} else {
			return int(tcR.DstX), int(tcR.DstY), nil
		}

		//		log.Printf("TranslateCoordinates eventX: %d, eventY: %d,  X: %d Y: %d\n", x, y, tcR.DstX, tcR.DstY)
		//		log.Printf("TranslateCoordinates Child:%v wn.Parent: %v\n", tcR.Child, w.Parent)
	} else {
		return x, y, nil
	}
}
func (w *Window) createMouseEvent(evType string, btn xproto.Button, eventX int16, eventY int16, evTime xproto.Timestamp) {
	prevButtons := w.Mbuttons

	evnt := Event{
		SWin: w,
		//		Kind:      Press,
		Source:    Mouse,
		Position:  image.Point{int(eventX), int(eventY)},
		Mbuttons:  w.Mbuttons, //uint8
		Time:      time.Duration(time.Duration(evTime)),
		Modifiers: getModifiers(),
	}
	if evType == "Press" {
		evnt.Kind = Press
		switch btn {
		case 1:
			w.Mbuttons = w.Mbuttons | ButtonPrimary

		case 2:
			w.Mbuttons = w.Mbuttons | ButtonTertiary

		case 3:
			w.Mbuttons = w.Mbuttons | ButtonSecondary
		}
	} else if evType == "Release" {
		evnt.Kind = Release
		switch btn {
		case 1:
			w.Mbuttons = w.Mbuttons ^ ButtonPrimary

		case 2:
			w.Mbuttons = w.Mbuttons ^ ButtonTertiary

		case 3:
			w.Mbuttons = w.Mbuttons ^ ButtonSecondary
		}
	}
	evnt.Mbuttons = w.Mbuttons ^ prevButtons // меняющееся состояние

	Wind.Config.EventChan <- evnt
}

func GetFileVersion() string {
	return ""
}

// Обрабатываем ALT CTRL SHIFT
// Левый и правый считаем за одно и то же
func SetKeyState(key xproto.Keysym, state bool) Modifiers {
	if key == 0xffe1 || key == 0xffe2 {
		if state {
			Wind.ModKeyState |= (ModShift)
		} else {
			Wind.ModKeyState ^= (ModShift)
		}
	} else if key == 0xffe3 || key == 0xffe4 {
		if state {
			Wind.ModKeyState |= (ModCtrl)
		} else {
			Wind.ModKeyState ^= (ModCtrl)
		}
	} else if key == 0xffe9 || key == 0xffea {
		if state {
			Wind.ModKeyState |= (ModAlt)
		} else {
			Wind.ModKeyState ^= (ModAlt)
		}
	}
	return Wind.ModKeyState
}

// Надо реализовать самому
// LSHIFT 0x0032 0xffe1 0x10
// RSHIFT 0x003e 0xffe2 0x10
// CAPSLOCK 0x0042  0xffe5 -
// LCtrl 0x0025 0xffe3 0x11
// RCtl 0x0069 0xffe4 0x11
// LAlt 0x0040 0xffe9 0x12
// RAlt 0x006c 0xffeA 0x12
// LWin --     --      0x5b
// RWin --     --      0x5C
// Menu 0x0087 0xff67

func getModifiers() Modifiers {
	return Wind.ModKeyState
}

// Заглушка для совместимости с Windows
func SendMessage(hwnd xproto.Window, m uint32, wParam, lParam uint32) {

}

// Меняем положение окна
func SetWindowPos(hwnd xproto.Window,
	HWND_TOPMOST,
	x, y, w, h, move int32,
) {

	var mask = uint16(0)
	var values []uint32 = make([]uint32, 0)

	wind, exists := WinMap.Load(hwnd)
	wn := &Window{}
	if exists {
		wn = wind.(*Window)

		if wn.Config.Position.X != int(x) || move != 0 {
			wn.Config.Position.X = int(x)
			mask |= xproto.ConfigWindowX
			values = append(values, uint32(x))
		}
		if wn.Config.Position.Y != int(y) || move != 0 {
			wn.Config.Position.Y = int(y)
			mask |= xproto.ConfigWindowY
			values = append(values, uint32(y))
		}
		if wn.Config.Size.X != int(w) || move != 0 {
			wn.Config.Size.X = int(w)
			mask |= xproto.ConfigWindowWidth
			values = append(values, uint32(w))
		}
		if wn.Config.Size.Y != int(h) || move != 0 {
			wn.Config.Size.Y = int(h)
			mask |= xproto.ConfigWindowHeight
			values = append(values, uint32(h))
		}
		if mask != 0 {
			xproto.ConfigureWindow(X, hwnd, mask, values)
			wn.draw()
		}
	}
}

func SetIcon(smenu int) {
	// Установка иконки окна
	var property xproto.Atom
	propertyC := xproto.InternAtom(X, true, uint16(len("_NET_WM_ICON")), "_NET_WM_ICON")
	propertyA, err := propertyC.Reply()
	if err == nil {
		property = propertyA.Atom
	} else {
		log.Println(err.Error())
		return
	}

	var mode byte = uint8(xproto.PropModeReplace)
	var pformat byte = 32
	var ptype xproto.Atom = xproto.AtomCardinal

	ndata, dataP, err := img.LoadIcon(smenu < 2)

	if err == nil {
		err = xproto.ChangePropertyChecked(X, mode, Wind.Hwnd, property, ptype, pformat, uint32(ndata), dataP).Check()
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}

}

func CloseWindow() {
	xproto.DestroyWindow(X, Wind.Hwnd)
}

// Отрисовка окна
func (w *Window) draw() {
	if !w.IsMain { // дочернее окно
		draw := xproto.Drawable(w.Hwnd)
		font, err := xproto.NewFontId(X)
		if err != nil {
			log.Println("error creating font id:", err)
			return
		} else {
			// X Logical Font Description Conventions
			//-FOUNDRY-FAMILY_NAME-WEIGHT_NAME-SLANT-SETWIDTH_NAME-ADD_STYLE_NAME-PIXEL_SIZE-POINT_SIZE-RESOLUTION_X
			// -RESOLUTION_Y-SPACING-AVERAGE_WIDTH-CHARSET_REGISTRY-CHARSET_ENCODING
			// можно подбирать через xfontsel
			fontname := "-*-*-bold-r-normal--" + strconv.Itoa(int(w.Config.FontSize)) + "-*-*-*-*-*-*-r"
			err = xproto.OpenFontChecked(X, font, uint16(len(fontname)), fontname).Check()

			if err != nil {
				fmt.Println("failed opening the font:", err)
				return
			} else {

				// And create a context from it. We simply pass the font's ID to the GcFont property.
				textCtx, err := xproto.NewGcontextId(X) //uint32
				if err != nil {
					fmt.Println("error creating text context:", err)
					return
				}

				mask := uint32(xproto.GcForeground | xproto.GcBackground | xproto.GcFont)
				values := []uint32{w.Config.TextColor, w.Config.BgColor, uint32(font)}
				xproto.CreateGC(X, textCtx, draw, mask, values)
				text := convertStringToChar2b(w.Config.Title)
				top := int16(25)
				if strings.ToUpper(w.Config.Class) == "BUTTON" {
					top = 18
				}
				xproto.ImageText16(X, byte(len(text)), draw, textCtx, 5, top, text) // по вертикали считается от верха до базовой линии
				// Close the font handle:
				xproto.CloseFont(X, font)
			}
		}
		/*
			   // Если требуется рамка (не border!) в определенной позиции
						   thick, err := xproto.NewGcontextId(X)
						   if err != nil {
							   fmt.Println("error creating thick context:", err)
							   return
						   } else {

							   mask := uint32(xproto.GcLineWidth)
							   values := []uint32{2}
							   xproto.CreateGC(X, thick, draw, mask, values)
							   rectangles := []xproto.Rectangle{
								   {X: 0, Y: 0, Width: 190, Height: 29},
								   {X: 180, Y: 20, Width: 10, Height: 10},
							   }
							   xproto.PolyRectangle(X, draw, thick, rectangles)
						   }
		*/
	}
}

func SysLog(level int, msg string) {
	sysLog, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_ERR, "wingui3")
	if err != nil {
		log.Println(msg)
		return
	}
	if level == 1 {
		sysLog.Err(msg)
	} else {
		sysLog.Info(msg)
	}
}
