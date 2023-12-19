//go:build linux
// +build linux

package winapi

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

const ID_BUTTON_1 = 101
const ID_BUTTON_2 = 102

const HWND_TOPMOST = -1
const SWP_NOMOVE = 2

type Window struct {
	Hwnd      *xproto.Window
	Childrens map[int]*Window
	Config    Config
	Buttons   MButtons
}

var X *xgb.Conn
var err error
var win Window

func CreateNativeMainWindow(config Config) (*Window, error) {
	fmt.Println("Create Main Window")
	win = Window{}

	X, err = xgb.NewConn()
	if err != nil {
		return &win, err
	}
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)
	wnd, _ := xproto.NewWindowId(X)
	xproto.CreateWindow(X, screen.RootDepth, wnd, screen.Root,
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
			0xffffffff, //xproto.CwBackPixel
			xproto.EventMaskStructureNotify |
				xproto.EventMaskKeyPress | xproto.EventMaskKeyRelease |
				xproto.EventMaskEnterWindow | xproto.EventMaskLeaveWindow |
				xproto.EventMaskButton1Motion | xproto.EventMaskButton2Motion |
				xproto.EventMaskButton3Motion | xproto.EventMaskButtonMotion |
				xproto.EventMaskButtonPress | xproto.EventMaskButtonRelease |
				xproto.EventMaskPointerMotion,
		})

	err = xproto.MapWindowChecked(X, wnd).Check()
	if err != nil {
		return &win, err
	}

	win.Hwnd = (&wnd)
	win.Childrens = make(map[int]*Window)
	win.Config = config

	return &win, nil
}

func CreateLabel(win *Window, config Config) (*Window, error) {
	return nil, nil
}

// Заглушка
func SendMessage(hwnd *xproto.Window, m uint32, wParam, lParam uint32) int32 {
	return 0
}

func SetWindowPos(hwnd *xproto.Window,
	HWND_TOPMOST,
	x, y, w, h, move int32,
) {

}

func Loop() {
	for {
		ev, xerr := X.WaitForEvent()
		if ev == nil && xerr == nil {
			log.Println("Both event and error are nil. Exiting...")
			return
		}

		//	if ev != nil {
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
				win.Buttons = win.Buttons | ButtonPrimary

			case 2:
				win.Buttons = win.Buttons | ButtonTertiary

			case 3:
				win.Buttons = win.Buttons | ButtonSecondary
			}
			win.Config.EventChan <- Event{
				SWin:      &win,
				Kind:      Press,
				Source:    Mouse,
				Position:  image.Point{int(bpe.EventX), int(bpe.EventY)},
				Mbuttons:  win.Buttons, //uint8
				Time:      time.Duration(bpe.Time),
				Modifiers: getModifiers(),
			}
		case xproto.ButtonReleaseEvent:
			bpe := ev.(xproto.ButtonReleaseEvent)
			btn := bpe.Detail
			switch btn {
			case 1:
				win.Buttons = win.Buttons ^ ButtonPrimary

			case 2:
				win.Buttons = win.Buttons ^ ButtonTertiary

			case 3:
				win.Buttons = win.Buttons ^ ButtonSecondary
			}
			win.Config.EventChan <- Event{
				SWin:      &win,
				Kind:      Release,
				Source:    Mouse,
				Position:  image.Point{int(bpe.EventX), int(bpe.EventY)},
				Mbuttons:  win.Buttons, //uint8
				Time:      time.Duration(bpe.Time),
				Modifiers: getModifiers(),
			}

		case xproto.MotionNotifyEvent:
			mne := ev.(xproto.MotionNotifyEvent)
			//			fmt.Println("Motion notify Event ", mne.Event) // Event  == *win.Hwnd
			//			fmt.Println("Motion notify *win.Hwnd ", *win.Hwnd)
			// mne.State  - номер кнопки
			win.Config.EventChan <- Event{
				SWin:      &win,
				Kind:      Move,
				Source:    Mouse,
				Position:  image.Point{int(mne.EventX), int(mne.EventY)},
				Mbuttons:  win.Buttons, //uint8
				Time:      time.Duration(mne.Time),
				Modifiers: getModifiers(),
			}

		case xproto.ReparentNotifyEvent:
			rne := ev.(xproto.ReparentNotifyEvent)
			fmt.Println("Reparent notify ", rne)

		case xproto.ConfigureNotifyEvent:
			cne := ev.(xproto.ConfigureNotifyEvent)
			fmt.Println("Configure notify ", cne)
			//			fmt.Println("Configure notify ", cne.Event) //  cne.Event == cne.Window == *win.Hwnd
			//			fmt.Println("Configure notify ", cne.Window)
			//			fmt.Println("Configure notify ", *win.Hwnd)

		case xproto.MapNotifyEvent:
			mne := ev.(xproto.MapNotifyEvent)
			fmt.Println("Map notify ", mne)

		case xproto.ResizeRequestEvent:
			mne := ev.(xproto.ResizeRequestEvent)
			fmt.Println("Resize Request ", mne)

		case xproto.ExposeEvent:
			ee := ev.(xproto.ExposeEvent)
			fmt.Println("Expose Event ", ee)

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
