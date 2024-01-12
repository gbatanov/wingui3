//go:build linux
// +build linux

package winapi

import (
	"fmt"
	"image"
	"log"
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
		///	log.Println("Event", ev)

		switch ev := ev.(type) {
		case xproto.CreateNotifyEvent:
			//			log.Println("CreateNotifyEvent", ev)

		case xproto.KeyPressEvent:
			//			log.Printf("Key pressed: %d\n", ev.Detail)
			if ev.Detail == VK_Q { //0x18
				return // exit on q
			}
		case xproto.KeyReleaseEvent:
			//			log.Printf("Key released: %d\n", ev.Detail)

		case xproto.ButtonPressEvent:
			w := getWindow(ev.Event)
			w.createMouseEvent("Press", ev.Detail, ev.EventX, ev.EventY, ev.Time)

		case xproto.ButtonReleaseEvent:
			w := getWindow(ev.Event)
			w.createMouseEvent("Release", ev.Detail, ev.EventX, ev.EventY, ev.Time)

		case xproto.MotionNotifyEvent:

			Wind.Config.EventChan <- Event{
				SWin:      Wind,
				Kind:      Move,
				Source:    Mouse,
				Position:  image.Point{int(ev.EventX), int(ev.EventY)},
				Mbuttons:  Wind.Mbuttons, //uint8
				Time:      time.Duration(ev.Time),
				Modifiers: getModifiers(),
			}

		case xproto.ReparentNotifyEvent:
			//			log.Println("Reparent notify ", ev)

		case xproto.ConfigureNotifyEvent:
			// TODO: убрать для дочерних окон
			w := getWindow(ev.Event)
			// A window's size, position, border, and/or stacking order is reconfigured by calling XConfigureWindow().
			// The window's position in the stacking order is changed by calling XLowerWindow(), XRaiseWindow(), or XRestackWindows().
			// A window is moved by calling XMoveWindow().
			// A window's size is changed by calling XResizeWindow().
			// A window's size and location is changed by calling XMoveResizeWindow().
			// A window is mapped and its position in the stacking order is changed by calling XMapRaised().
			// A window's border width is changed by calling XSetWindowBorderWidth().

			//			log.Println("Configure notify ", ev)
			if ev.Width > uint16(w.Config.MaxSize.X) ||
				ev.Height > uint16(w.Config.MaxSize.Y) {
				SetWindowPos(ev.Event, HWND_TOPMOST,
					int32(w.Config.Position.X),
					int32(w.Config.Position.Y),
					int32(w.Config.MaxSize.X),
					int32(w.Config.MaxSize.Y), 0)

			} else if ev.Width < uint16(w.Config.MinSize.X) ||
				ev.Height < uint16(w.Config.MinSize.Y) {
				SetWindowPos(ev.Event, HWND_TOPMOST,
					int32(w.Config.Position.X),
					int32(w.Config.Position.Y),
					int32(w.Config.MinSize.X),
					int32(w.Config.MinSize.Y), 0)

			} else {
				w.Config.Position.X = int(ev.X)
				w.Config.Position.Y = int(ev.Y)
			}

		case xproto.MapNotifyEvent:
			//			log.Println("Map notify ", ev)

			//	case xproto.ResizeRequestEvent: // WM_SIZE Работает криво, отключил в маске
			//		log.Println("Resize Request ", ev)

		case xproto.ClientMessageEvent:
			log.Println("ClientMessage Event ", ev)

		case xproto.ExposeEvent: // аналог WM_PAINT в Windows
			w := getWindow(ev.Window)
			w.draw()

		case xproto.DestroyNotifyEvent:
			// На закрытие по крестику не приходит
			// Событие приходит для каждого окна (главное и дочерние)
			// Будем отправлять событие только для главного окна
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

func (w *Window) createMouseEvent(evType string, btn xproto.Button, eventX int16, eventY int16, evTime xproto.Timestamp) {
	prevButtons := w.Mbuttons

	// При щелчке в дочернем окне можно оттранслировать координаты относительно дочернего окна
	// в координаты относительно родительского.
	if w != Wind {
		tc := xproto.TranslateCoordinates(X, w.Hwnd, w.Parent, eventX, eventY)
		tcR, err := tc.Reply()
		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Printf("TranslateCoordinates eventX: %d, eventY: %d,  X: %d Y: %d\n", eventX, eventY, tcR.DstX, tcR.DstY)
		log.Printf("TranslateCoordinates Child:%v wn.Parent: %v\n", tcR.Child, w.Parent)
	}

	evnt := Event{
		SWin: w,
		//		Kind:      Press,
		Source:    Mouse,
		Position:  image.Point{int(eventX), int(eventY)},
		Mbuttons:  w.Mbuttons, //uint8
		Time:      time.Duration(evTime),
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
func GetKeyState(key int32) int16 {
	return 0
}

// Заглушка для совместимости с Windows
func SendMessage(hwnd xproto.Window, m uint32, wParam, lParam uint32) {

}

// Меняем положение окна
func SetWindowPos(hwnd xproto.Window,
	HWND_TOPMOST,
	x, y, w, h, move int32,
) {

	wind, exists := WinMap.Load(hwnd)
	wn := &Window{}
	if exists {
		wn = wind.(*Window)
	}
	mask := uint16(xproto.ConfigWindowX | xproto.ConfigWindowY | xproto.ConfigWindowWidth | xproto.ConfigWindowHeight)
	values := []uint32{uint32(x), uint32(y), uint32(w), uint32(h)}

	xproto.ConfigureWindow(X, hwnd, mask, values)
	wn.Config.Position.X = int(x)
	wn.Config.Position.Y = int(y)
	wn.Config.Size.X = int(w)
	wn.Config.Size.Y = int(h)

	wn.draw()
}

func SetIcon() {
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

	ndata, dataP, err := img.LoadIcon()
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