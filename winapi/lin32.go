//go:build linux
// +build linux

package winapi

import (
	"fmt"
	"image"
	"log"
	"time"
	"unicode/utf16"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

const ID_BUTTON_1 = 101
const ID_BUTTON_2 = 102

const HWND_TOPMOST = -1
const SWP_NOMOVE = 2

type Window struct {
	Hwnd      xproto.Window
	Childrens map[int]*Window
	Config    Config
	Mbuttons  MButtons
}

var X *xgb.Conn
var err error
var win *Window

func CreateNativeMainWindow(config Config) (*Window, error) {
	fmt.Println("Create Main Window")
	win = &Window{}

	X, err = xgb.NewConn()
	if err != nil {
		return win, err
	}
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)

	screenX := screen.WidthInPixels
	screenY := screen.HeightInPixels
	if config.Position.X < 0 {
		config.Position.X = int(screenX) + config.Position.X - config.Size.X
	}
	if config.Position.Y < 0 {
		config.Position.Y = int(screenY) + config.Position.Y - config.Size.Y
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
				xproto.EventMaskPointerMotion,
		})

	err = xproto.MapWindowChecked(X, wnd).Check()
	if err != nil {
		return win, err
	}

	win.Hwnd = (wnd)
	win.Childrens = make(map[int]*Window)
	win.Config = config

	WinMap.Store(win.Hwnd, win)

	return win, nil
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

	err = xproto.MapWindowChecked(X, wndL).Check()
	if err != nil {
		return chWin, err
	}

	chWin.Hwnd = (wndL)
	chWin.Childrens = make(map[int]*Window, 0)
	chWin.Config = config

	WinMap.Store(chWin.Hwnd, chWin)

	foreground, err := xproto.NewGcontextId(X)
	if err != nil {
		fmt.Println("error creating foreground context:", err)
		return chWin, nil
	} else {
		draw := xproto.Drawable(wndL)
		mask := uint32(xproto.GcForeground)
		values := []uint32{screen.BlackPixel}
		xproto.CreateGC(X, foreground, draw, mask, values)

		red, err := xproto.NewGcontextId(X)
		if err != nil {
			fmt.Println("error creating red context:", err)
			return chWin, nil
		} else {

			mask = uint32(xproto.GcForeground)
			values = []uint32{0xff0000}
			xproto.CreateGC(X, red, draw, mask, values)

			// We'll create another graphics context that draws thick lines:
			thick, err := xproto.NewGcontextId(X)
			if err != nil {
				fmt.Println("error creating thick context:", err)
				return chWin, nil
			} else {

				mask = uint32(xproto.GcLineWidth)
				values = []uint32{10}
				xproto.CreateGC(X, thick, draw, mask, values)

				blue, err := xproto.NewGcontextId(X)
				if err != nil {
					fmt.Println("error creating blue context:", err)
					return chWin, nil
				} else {

					mask = uint32(xproto.GcForeground | xproto.GcLineWidth)
					values = []uint32{0x0000ff, 4}
					xproto.CreateGC(X, blue, draw, mask, values)

					mask = uint32(xproto.GcLineWidth | xproto.GcCapStyle)

					values = []uint32{3, xproto.CapStyleRound}
					xproto.ChangeGC(X, foreground, mask, values)

				}
			}
		}
	}
	return chWin, nil
}
func convertStringToChar2b(s string) []xproto.Char2b {
	var chars []xproto.Char2b
	var p []uint16

	for _, r := range []rune(s) {
		p = utf16.Encode([]rune{r})
		if len(p) == 1 {
			chars = append(chars, convertUint16ToChar2b(p[0]))
		} else {
			// If the utf16 representation is larger than 2 bytes
			// we can not use it and insert a blank instead:
			chars = append(chars, xproto.Char2b{Byte1: 0, Byte2: 32})
		}
	}

	return chars
}

// convertUint16ToChar2b converts a uint16 (which is basically two bytes)
// into a Char2b by using the higher 8 bits of u as Byte1
// and the lower 8 bits of u as Byte2.
func convertUint16ToChar2b(u uint16) xproto.Char2b {
	return xproto.Char2b{
		Byte1: byte((u & 0xff00) >> 8),
		Byte2: byte((u & 0x00ff)),
	}
}

func Loop() {
	for {
		ev, xerr := X.WaitForEvent()
		if ev == nil && xerr == nil {
			log.Println("Both event and error are nil. Exiting...")
			return
		}

		//		if ev != nil {
		//			log.Printf("Event: %s\n", ev)
		//		}

		if xerr != nil {
			log.Printf("Error: %s\n", xerr)
		}

		switch ev.(type) {
		case xproto.KeyPressEvent:
			kpe := ev.(xproto.KeyPressEvent)
			fmt.Printf("Key pressed: %d\n", kpe.Detail)
			if kpe.Detail == VK_Q { //0x18
				return // exit on q
			}
		case xproto.KeyReleaseEvent:
			kpe := ev.(xproto.KeyReleaseEvent)
			fmt.Printf("Key released: %d\n", kpe.Detail)

		case xproto.ButtonPressEvent:
			bpe := ev.(xproto.ButtonPressEvent)
			btn := bpe.Detail
			switch btn {
			case 1:
				win.Mbuttons = win.Mbuttons | ButtonPrimary

			case 2:
				win.Mbuttons = win.Mbuttons | ButtonTertiary

			case 3:
				win.Mbuttons = win.Mbuttons | ButtonSecondary
			}
			win.Config.EventChan <- Event{
				SWin:      win,
				Kind:      Press,
				Source:    Mouse,
				Position:  image.Point{int(bpe.EventX), int(bpe.EventY)},
				Mbuttons:  win.Mbuttons, //uint8
				Time:      time.Duration(bpe.Time),
				Modifiers: getModifiers(),
			}
		case xproto.ButtonReleaseEvent:
			bpe := ev.(xproto.ButtonReleaseEvent)
			btn := bpe.Detail
			switch btn {
			case 1:
				win.Mbuttons = win.Mbuttons ^ ButtonPrimary

			case 2:
				win.Mbuttons = win.Mbuttons ^ ButtonTertiary

			case 3:
				win.Mbuttons = win.Mbuttons ^ ButtonSecondary
			}
			win.Config.EventChan <- Event{
				SWin:      win,
				Kind:      Release,
				Source:    Mouse,
				Position:  image.Point{int(bpe.EventX), int(bpe.EventY)},
				Mbuttons:  win.Mbuttons, //uint8
				Time:      time.Duration(bpe.Time),
				Modifiers: getModifiers(),
			}

		case xproto.MotionNotifyEvent:
			mne := ev.(xproto.MotionNotifyEvent)
			//			fmt.Println("Motion notify Event ", mne.Event) // Event  == *win.Hwnd
			//			fmt.Println("Motion notify *win.Hwnd ", *win.Hwnd)
			// mne.State  - номер кнопки
			win.Config.EventChan <- Event{
				SWin:      win,
				Kind:      Move,
				Source:    Mouse,
				Position:  image.Point{int(mne.EventX), int(mne.EventY)},
				Mbuttons:  win.Mbuttons, //uint8
				Time:      time.Duration(mne.Time),
				Modifiers: getModifiers(),
			}

		case xproto.ReparentNotifyEvent:
			rne := ev.(xproto.ReparentNotifyEvent)
			log.Println("Reparent notify ", rne)

		case xproto.ConfigureNotifyEvent:
			cne := ev.(xproto.ConfigureNotifyEvent)
			log.Println("Configure notify ", cne)
			//			fmt.Println("Configure notify ", cne.Event) //  cne.Event == cne.Window == *win.Hwnd
			//			fmt.Println("Configure notify ", cne.Window)
			//			fmt.Println("Configure notify ", *win.Hwnd)

		case xproto.MapNotifyEvent:
			mne := ev.(xproto.MapNotifyEvent)
			log.Println("Map notify ", mne)

		case xproto.ResizeRequestEvent:
			mne := ev.(xproto.ResizeRequestEvent)
			fmt.Println("Resize Request ", mne)

		case xproto.ExposeEvent:
			ee := ev.(xproto.ExposeEvent)
			log.Println("Expose Event ", ee)
			// Writing text needs a bit more setup -- we first have
			// to open the required font.

			wind, exists := WinMap.Load(ee.Window)
			w := &Window{}
			if exists {
				w = wind.(*Window)
				log.Println(w.Config.Title)
			}

			draw := xproto.Drawable(ee.Window)
			font, err := xproto.NewFontId(X)
			if err != nil {
				fmt.Println("error creating font id:", err)
				return
			} else {
				fontname := "-*-fixed-*-*-*-*-14-*-*-*-*-*-*-*"
				err = xproto.OpenFontChecked(X, font, uint16(len(fontname)), fontname).Check()

				if err != nil {
					fmt.Println("failed opening the font:", err)
					return
				} else {

					// And create a context from it. We simply pass the font's ID to the GcFont property.
					textCtx, err := xproto.NewGcontextId(X)
					if err != nil {
						fmt.Println("error creating text context:", err)
						return
					}

					mask := uint32(xproto.GcForeground | xproto.GcBackground | xproto.GcFont)
					values := []uint32{w.Config.TextColor, w.Config.BgColor, uint32(font)}
					xproto.CreateGC(X, textCtx, draw, mask, values)
					text := convertStringToChar2b(w.Config.Title) // Unicode capable!
					xproto.ImageText16(X, byte(len(text)), draw, textCtx, 10, 20, text)
					// Close the font handle:
					xproto.CloseFont(X, font)
				}
			}

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
		case xproto.DestroyNotifyEvent:

			return
		}
	}
}
func GetFileVersion() string {
	return ""
}
func GetKeyState(key int32) int16 {
	return 0
}

// Заглушка
func SendMessage(hwnd xproto.Window, m uint32, wParam, lParam uint32) int32 {
	return 0
}

func SetWindowPos(hwnd xproto.Window,
	HWND_TOPMOST,
	x, y, w, h, move int32,
) {

}
