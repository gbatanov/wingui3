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
	"unicode/utf16"

	"github.com/gbatanov/wingui3/img"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

/*
atomPRIMARY            = 1
atomSECONDARY          = 2
atomARC                = 3
atomATOM               = 4
atomBITMAP             = 5
atomCARDINAL           = 6
atomCOLORMAP           = 7
atomCURSOR             = 8
atomCUT_BUFFER0        = 9
atomCUT_BUFFER1        = 10
atomCUT_BUFFER2        = 11
atomCUT_BUFFER3        = 12
atomCUT_BUFFER4        = 13
atomCUT_BUFFER5        = 14
atomCUT_BUFFER6        = 15
atomCUT_BUFFER7        = 16
atomDRAWABLE           = 17
atomFONT               = 18
atomINTEGER            = 19
atomPIXMAP             = 20
atomPOINT              = 21
atomRECTANGLE          = 22
atomRESOURCE_MANAGER   = 23
atomRGB_COLOR_MAP      = 24
atomRGB_BEST_MAP       = 25
atomRGB_BLUE_MAP       = 26
atomRGB_DEFAULT_MAP    = 27
atomRGB_GRAY_MAP       = 28
atomRGB_GREEN_MAP      = 29
atomRGB_RED_MAP        = 30
atomSTRING             = 31
atomVISUALID           = 32
atomWINDOW             = 33
atomWM_COMMAND         = 34
atomWM_HINTS           = 35
atomWM_CLIENT_MACHINE  = 36
atomWM_ICON_NAME       = 37
atomWM_ICON_SIZE       = 38
atomWM_NAME            = 39
atomWM_NORMAL_HINTS    = 40
atomWM_SIZE_HINTS      = 41
atomWM_ZOOM_HINTS      = 42
atomMIN_SPACE          = 43
atomNORM_SPACE         = 44
atomMAX_SPACE          = 45
atomEND_SPACE          = 46
atomSUPERSCRIPT_X      = 47
atomSUPERSCRIPT_Y      = 48
atomSUBSCRIPT_X        = 49
atomSUBSCRIPT_Y        = 50
atomUNDERLINE_POSITION = 51
atomUNDERLINE_THICKNESS= 52
atomSTRIKEOUT_ASCENT   = 53
atomSTRIKEOUT_DESCENT  = 54
atomITALIC_ANGLE       = 55
atomX_HEIGHT           = 56
atomQUAD_WIDTH         = 57
atomWEIGHT             = 58
atomPOINT_SIZE         = 59
atomRESOLUTION         = 60
atomCOPYRIGHT          = 61
atomNOTICE             = 62
atomFONT_NAME          = 63
atomFAMILY_NAME        = 64
atomFULL_NAME          = 65
atomCAP_HEIGHT         = 66
atomWM_CLASS           = 67
atomWM_TRANSIENT_FOR   = 68
*/

const HWND_TOPMOST = -1
const SWP_NOMOVE = 2

type WND_KIND int

type Window struct {
	Hwnd      xproto.Window
	Childrens map[int]*Window
	Config    Config
	Mbuttons  MButtons // здесь состав нажатых кнопок
	Parent    xproto.Window
	IsMain    bool
}

var X *xgb.Conn
var err error
var win *Window

func CreateNativeMainWindow(config Config) (*Window, error) {
	//log.Println("Create Main Window")
	win = &Window{}

	X, err = xgb.NewConn()
	if err != nil {
		return win, err
	}
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)
	/*
		err := randr.Init(X)
		if err != nil {
			log.Fatal(err)
		}


		// Gets the current screen resources. Screen resources contains a list
		// of names, crtcs, outputs and modes, among other things.
		resources, err := randr.GetScreenResources(X, screen.Root).Reply()
		if err != nil {
			log.Fatal(err)
		}

		// Iterate through all of the outputs and show some of their info.
		for _, output := range resources.Outputs {
			info, err := randr.GetOutputInfo(X, output, 0).Reply()
			if err != nil {
				log.Fatal(err)
			}

			if info.Connection == randr.ConnectionConnected {
				bestMode := info.Modes[0]
				for _, mode := range resources.Modes {
					if mode.Id == uint32(bestMode) {
						log.Printf("Best mode: Width: %d, Height: %d\n", mode.Width, mode.Height)
					}
				}
			}
		}

		// Iterate through all of the crtcs and show some of their info.
		for _, crtc := range resources.Crtcs {
			info, err := randr.GetCrtcInfo(X, crtc, 0).Reply()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("CRT: X: %d, Y: %d, Width: %d, Height: %d\n", info.X, info.Y, info.Width, info.Height)
		}

		// Tell RandR to send us events. (I think these are all of them, as of 1.3.)
		err = randr.SelectInputChecked(X, screen.Root,
			randr.NotifyMaskScreenChange|
				randr.NotifyMaskCrtcChange|
				randr.NotifyMaskOutputChange|
				randr.NotifyMaskOutputProperty).Check()
		if err != nil {
			log.Fatal(err)
		}
	*/
	screenX := screen.WidthInPixels
	screenY := screen.HeightInPixels
	if config.Position.X < 0 {
		config.Position.X = int(screenX) + config.Position.X - config.Size.X
	}
	if config.Position.Y < 0 {
		config.Position.Y = int(screenY) + config.Position.Y - config.Size.Y - 48
	}

	wnd, _ := xproto.NewWindowId(X)

	xproto.CreateWindow(X,
		screen.RootDepth,
		wnd,
		screen.Root,
		int16(config.Position.X),
		int16(config.Position.Y),
		uint16(config.Size.X),
		uint16(config.Size.Y),
		uint16(config.BorderSize.X),
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		0,
		[]uint32{})

	xproto.ChangeWindowAttributes(X, wnd,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{ // values must be in the order defined by the protocol
			config.BgColor,
			xproto.EventMaskStructureNotify | xproto.EventMaskExposure |
				xproto.EventMaskKeyPress | xproto.EventMaskKeyRelease |
				xproto.EventMaskEnterWindow | xproto.EventMaskLeaveWindow |
				xproto.EventMaskButton1Motion | xproto.EventMaskButton2Motion |
				xproto.EventMaskButton3Motion | xproto.EventMaskButtonMotion |
				xproto.EventMaskButtonPress | xproto.EventMaskButtonRelease |
				xproto.EventMaskPointerMotion /*| xproto.EventMaskResizeRedirect*/})

	//log.Println("Before MapWindow Main")
	err = xproto.MapWindowChecked(X, wnd).Check()
	if err != nil {
		return win, err
	}

	// Установка заголовка окна (работает)
	var mode byte = xproto.PropModeReplace
	var property xproto.Atom = xproto.AtomWmName
	var ptype xproto.Atom = xproto.AtomString
	var pformat byte = 8
	var data []byte = []byte(config.Title)
	datalen := len(data)

	err = xproto.ChangePropertyChecked(X, mode, wnd, property, ptype, pformat, uint32(datalen), data).Check()
	if err != nil {
		log.Println(err.Error())
		err = nil
	}

	win.Hwnd = (wnd)
	win.Childrens = make(map[int]*Window)
	win.Config = config
	win.Parent = screen.Root
	win.IsMain = true

	WinMap.Store(win.Hwnd, win)
	WinMap.Store(0, win) // Основное окно дублируем с нулевым ключчом, чтобы иметь доступ всегда

	return win, nil
}

func CreateButton(win *Window, config Config) (*Window, error) {
	w, err := CreateLabel(win, config)
	return w, err
}

func CreateLabel(win *Window, config Config) (*Window, error) {
	chWin := &Window{}

	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)
	wndL, _ := xproto.NewWindowId(X)

	xproto.CreateWindow(X,
		screen.RootDepth,
		wndL,
		win.Hwnd,
		int16(config.Position.X),
		int16(config.Position.Y),
		uint16(config.Size.X),
		uint16(config.Size.Y),
		uint16(config.BorderSize.X),
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		0,
		[]uint32{})

	xproto.ChangeWindowAttributes(X, wndL,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{
			config.BgColor,
			xproto.EventMaskStructureNotify | xproto.EventMaskExposure |
				xproto.EventMaskKeyPress | xproto.EventMaskKeyRelease |
				xproto.EventMaskEnterWindow | xproto.EventMaskLeaveWindow |
				xproto.EventMaskButton1Motion | xproto.EventMaskButton2Motion |
				xproto.EventMaskButton3Motion | xproto.EventMaskButtonMotion |
				xproto.EventMaskButtonPress | xproto.EventMaskButtonRelease |
				xproto.EventMaskPointerMotion,
		})

	//log.Println("Before MapWindow Child")
	err = xproto.MapWindowChecked(X, wndL).Check()
	if err != nil {
		return chWin, err
	}

	chWin.Hwnd = (wndL)
	chWin.Childrens = make(map[int]*Window, 0)
	chWin.Config = config
	chWin.Parent = win.Hwnd
	chWin.IsMain = false

	WinMap.Store(chWin.Hwnd, chWin)

	return chWin, nil
}

func convertStringToChar2b(s string) []xproto.Char2b {
	var chars []xproto.Char2b
	var p []uint16

	for _, r := range []rune(s) {
		p = utf16.Encode([]rune{r})
		if len(p) == 1 { // uint16, 2 байта
			chars = append(chars, convertUint16ToChar2b(p[0]))
		} else {
			// Если вернулось больше двух байт, вставляем пробел(TODO: символ-заменитель)
			chars = append(chars, xproto.Char2b{Byte1: 0, Byte2: 0x20})
		}
	}

	return chars
}

// Переводит uint16 в двух байтовую структуру
func convertUint16ToChar2b(u uint16) xproto.Char2b {
	return xproto.Char2b{
		Byte1: byte((u & 0xff00) >> 8),
		Byte2: byte((u & 0x00ff)),
	}
}

// Основной цикл обработки событий
func Loop() {
	for {
		ev, xerr := X.WaitForEvent()
		if ev == nil && xerr == nil {
			log.Println("Both event and error are nil. Exiting...")
			return
		}

		if xerr != nil {
			log.Printf("Error: %s\n", xerr)
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

			win.Config.EventChan <- Event{
				SWin:      win,
				Kind:      Move,
				Source:    Mouse,
				Position:  image.Point{int(ev.EventX), int(ev.EventY)},
				Mbuttons:  win.Mbuttons, //uint8
				Time:      time.Duration(ev.Time),
				Modifiers: getModifiers(),
			}

		case xproto.ReparentNotifyEvent:
			//			log.Println("Reparent notify ", ev)

		case xproto.ConfigureNotifyEvent: // Идет только для главного окна
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

			return
		} // switch
	} //for
} //Loop

func getWindow(wev xproto.Window) *Window {
	w := win
	if wev != win.Hwnd {
		wind, exists := WinMap.Load(wev)
		if exists {
			w = wind.(*Window)
		}
	}
	return w
}

func (w *Window) createMouseEvent(evType string, btn xproto.Button, eventX int16, eventY int16, evTime xproto.Timestamp) {
	prevButtons := w.Mbuttons
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
	win.Config.EventChan <- evnt
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
	//	log.Printf("PosX: %d,  PosY: %d , x: %d, Y: %d \n", wn.Config.Position.X, wn.Config.Position.Y, x, y)
	/*
		tc := xproto.TranslateCoordinates(X, hwnd, wn.Parent, int16(0), int16(0))
		tcR, err := tc.Reply()
		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Printf("TranslateCoordinates X:%d Y:%d\n", tcR.DstX, tcR.DstY)

			log.Printf("TranslateCoordinates Child:%v wn.Parent: %v\n", tcR.Child, wn.Parent)

				xwa := xproto.GetWindowAttributes(X, hwnd)
				xwaR, err := xwa.Reply()
				if err != nil {
					log.Println(err.Error())
					return
				}
					// На Ubuntu это дает верную координату по Y , но некорректную по X
			// На Астра-Линукс обе кординаты некорректны
			//values := []uint32{uint32(tcR.DstX), uint32(tcR.DstY), uint32(w), uint32(h)}

	*/
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
		err = xproto.ChangePropertyChecked(X, mode, win.Hwnd, property, ptype, pformat, uint32(ndata), dataP).Check()
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}

}

func CloseWindow() {
	xproto.DestroyWindow(X, win.Hwnd)
}

// Отрисовка окна
func (w *Window) draw() {
	if !w.IsMain { // дочернее окно
		draw := xproto.Drawable(w.Hwnd)
		font, err := xproto.NewFontId(X)
		if err != nil {
			fmt.Println("error creating font id:", err)
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