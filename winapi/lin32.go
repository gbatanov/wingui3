//go:build linux
// +build linux

package winapi

import (
	"fmt"
	"log"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

const ID_BUTTON_1 = 101
const ID_BUTTON_2 = 102

const HWND_TOPMOST = -1
const SWP_NOMOVE = 2

var X *xgb.Conn
var err error

type Window struct {
	Hwnd      *xproto.Window
	Childrens map[int]*Window
	Config    Config
}

func GetKeyState(key int32) int16 {
	return 0
}

func CreateNativeMainWindow(config Config) (*Window, error) {
	fmt.Println("Create Main Window")
	win := Window{}

	X, err = xgb.NewConn()
	if err != nil {
		fmt.Println(err)
		return &win, err
	}
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)
	wnd, _ := xproto.NewWindowId(X)
	xproto.CreateWindow(X, screen.RootDepth, wnd, screen.Root,
		0, 0, 500, 500, 0,
		xproto.WindowClassInputOutput, screen.RootVisual, 0, []uint32{})
	xproto.ChangeWindowAttributes(X, wnd,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{ // values must be in the order defined by the protocol
			0xffffffff,
			xproto.EventMaskStructureNotify |
				xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease,
		})

	err = xproto.MapWindowChecked(X, wnd).Check()
	if err != nil {
		return &win, err
	} else {
		log.Printf("Map window %d successful!\n", wnd)
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
			fmt.Println("Both event and error are nil. Exiting...")
			return
		}

		if ev != nil {
			fmt.Printf("Event: %s\n", ev)
		}
		if xerr != nil {
			fmt.Printf("Error: %s\n", xerr)
		}

		switch ev.(type) {
		case xproto.KeyPressEvent:
			// See https://pkg.go.dev/github.com/jezek/xgb/xproto#KeyPressEvent
			// for documentation about a key press event.
			kpe := ev.(xproto.KeyPressEvent)
			fmt.Printf("Key pressed: %d\n", kpe.Detail)
			// The Detail value depends on the keyboard layout,
			// for QWERTY, q is #24.
			if kpe.Detail == 24 {
				return // exit on q
			}
		case xproto.DestroyNotifyEvent:

			return
		}
	}
}
func GetFileVersion() string {
	return ""
}
