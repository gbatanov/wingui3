//go:build linux
// +build linux

package winapi

import (
	"io"
	"log"
	"unicode/utf16"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

/*
// Атомы прописанные в X11
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
	Hwnd              xproto.Window
	Childrens         map[int]*Window
	Config            Config   // настройки окна
	Mbuttons          MButtons // здесь состав нажатых кнопок
	Parent            xproto.Window
	IsMain            bool
	Keymap            []xproto.Keysym
	KeysymsPerKeycode byte
	FirstCode         byte
	ModKeyState       Modifiers // Состояние клавиш-модификаторов
}

var X *xgb.Conn
var err error
var Wind *Window

// Create Main Window
func CreateNativeMainWindow(config Config) (*Window, error) {
	xgb.Logger = log.New(io.Discard, "", 0) // Давим внутренние сообщения от XGB

	Wind = &Window{}

	X, err = xgb.NewConn()
	if err != nil {
		return Wind, err
	}

	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)

	screenX := screen.WidthInPixels
	screenY := screen.HeightInPixels
	if config.Position.X < 0 {
		config.Position.X = int(screenX) + config.Position.X - config.Size.X
	}
	if config.Position.Y < 0 {
		config.Position.Y = int(screenY) + config.Position.Y - config.Size.Y - 48
	}

	Wind.FirstCode = byte(setup.MinKeycode)
	rep, err := xproto.GetKeyboardMapping(X, xproto.Keycode(Wind.FirstCode), byte(setup.MaxKeycode-setup.MinKeycode)).Reply()
	if err == nil {
		Wind.Keymap = rep.Keysyms
		Wind.KeysymsPerKeycode = rep.KeysymsPerKeycode
	} else {
		Wind.Keymap = make([]xproto.Keysym, 0)
		Wind.KeysymsPerKeycode = 0
		return Wind, err
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
		uint16(config.BorderSize),
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
		return Wind, err
	}

	// Установка заголовка окна
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

	Wind.Hwnd = (wnd)
	Wind.Childrens = make(map[int]*Window)
	Wind.Config = config
	Wind.Parent = screen.Root
	Wind.IsMain = true
	WinMap.Store(Wind.Hwnd, Wind)
	WinMap.Store(0, Wind) // Основное окно дублируем с нулевым ключчом, чтобы иметь доступ всегда

	return Wind, nil
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
		uint16(config.BorderSize),
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
