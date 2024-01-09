//go:generate go-winres make --file-version=v0.1.45.6 --product-version=git-tag
package main

import (
	_ "embed"
	"fmt"
	"image"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"fyne.io/systray"
	"github.com/gbatanov/wingui3/img"
	"github.com/gbatanov/wingui3/winapi"
)

var Version string = "v0.2.53" // Windows - подставится после генерации во время исполнения программы

const COLOR_GREEN = 0x0011aa11
const COLOR_RED = 0x000000c8
const COLOR_YELLOW = 0x0000c8c8
const COLOR_GRAY_DE = 0x00dedede
const COLOR_GRAY_BC = 0x00bcbcbc
const COLOR_GRAY_AA = 0x00aaaaaa

var mouseX, mouseY int = 0, 0
var serverList []string = []string{"192.168.76.106", "192.168.76.80"}

// Конфиг основного окна
var config = winapi.Config{
	Position:   image.Pt(1345, 20),
	MaxSize:    image.Pt(480, 240),
	MinSize:    image.Pt(200, 100),
	Size:       image.Pt(240, 100),
	Title:      "wingui3",
	TextColor:  COLOR_GREEN,
	EventChan:  make(chan winapi.Event, 256),
	BorderSize: image.Pt(1, 1),
	Mode:       winapi.Windowed,
	BgColor:    COLOR_GRAY_DE,
	SysMenu:    2,
	Class:      "GsbWindow",
}
var labelConfig = winapi.Config{
	Class:      "Static",
	Title:      "Static",
	EventChan:  config.EventChan,
	Size:       image.Pt(int(config.Size.X-40), int(30)),
	MinSize:    config.MinSize,
	MaxSize:    config.MaxSize,
	Position:   image.Pt(int(18), int(15)),
	Mode:       winapi.Windowed,
	BorderSize: image.Pt(0, 0),
	TextColor:  COLOR_GREEN,
	FontSize:   28,
	BgColor:    COLOR_GRAY_DE, //config.BgColor,
}
var btnConfig = winapi.Config{
	Class:      "Button",
	Title:      "Ok",
	EventChan:  config.EventChan,
	Size:       image.Pt(int(40), int(25)),
	Position:   image.Pt(int(18), int(15)),
	Mode:       winapi.Windowed,
	BorderSize: image.Pt(1, 1),
	TextColor:  COLOR_GREEN,
	FontSize:   16,
	BgColor:    COLOR_GRAY_AA,
}

var quit chan os.Signal
var flag = true

// ---------------------------------------------------------------
func main() {
	quit = make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	getFileVersion()
	config.Title += (" " + Version)
	win, err := winapi.CreateNativeMainWindow(config)
	if err == nil {

		// Обработчик событий
		go func() {
			for flag {
				select {
				case ev, ok := <-config.EventChan:
					if !ok {
						// канал закрыт
						flag = false
						break
					}
					switch ev.Source {
					case winapi.Mouse:
						MouseEventHandler(ev)
					case winapi.Frame:
						FrameEventHandler(ev)
					}

				case <-quit: // сообщение при закрытии трея
					flag = false
				} //select
			} //for

			winapi.SendMessage(win.Hwnd, winapi.WM_CLOSE, 0, 0)
		}()

		defer winapi.WinMap.Delete(win.Hwnd)

		var id int = 0

		// Label с текстом
		for _, title := range serverList {
			labelConfig.Title = title
			AddLabel(win, labelConfig, id)
			//			labelConfig.BgColor = COLOR_GRAY_AA
			id++
		}

		// Buttons
		// Ok
		btnConfig1 := btnConfig
		btnConfig1.ID = winapi.ID_BUTTON_1
		btnConfig1.Position.Y = 20 + (labelConfig.Size.Y)*(id)
		AddButton(win, btnConfig1, id)
		// Cancel
		id++
		btnConfig2 := btnConfig
		btnConfig2.Title = "Cancel"
		btnConfig2.ID = winapi.ID_BUTTON_2
		btnConfig2.Position.Y = btnConfig1.Position.Y
		btnConfig2.Position.X = btnConfig1.Position.X + btnConfig1.Size.X + 10
		btnConfig2.Size.X = 60
		AddButton(win, btnConfig2, id)

		if len(win.Childrens) > 0 {
			for _, w2 := range win.Childrens {
				defer winapi.WinMap.Delete(w2.Hwnd)
			}

			win.Config.Size.Y = 2*labelConfig.Size.Y + +btnConfig1.Size.Y + 30
			win.Config.MinSize.Y = win.Config.Size.Y
			win.Config.MaxSize.Y = win.Config.Size.Y
			///	win.Config.Position.Y // будет либо 27 (до 35), либо Y+27(35 и больше)
		}
		//		log.Printf("Before SetWindowPos PositionX %d PositionY %d", win.Config.Position.X, win.Config.Position.Y)
		winapi.SetWindowPos(win.Hwnd,
			winapi.HWND_TOPMOST,
			int32(win.Config.Position.X),
			int32(win.Config.Position.Y),
			int32(win.Config.Size.X),
			int32(win.Config.Size.Y),
			winapi.SWP_NOMOVE)

		//systray (На Астре-Линукс не работает)
		if runtime.GOOS == "windows" {
			go func() {
				systray.Run(onReady, onExit)
			}()
		}
		winapi.SetIcon()
		winapi.Loop()

		close(config.EventChan)
		fmt.Println("Quit")
	} else {
		panic(err.Error())
	}

}

func AddLabel(win *winapi.Window, lblConfig winapi.Config, id int) error {

	lblConfig.Position.Y = 10 + (lblConfig.Size.Y)*(id)
	chWin, err := winapi.CreateLabel(win, lblConfig)
	if err == nil {
		win.Childrens[id] = chWin

		return nil
	}
	return err
}

func AddButton(win *winapi.Window, btnConfig winapi.Config, id int) error {

	chWin, err := winapi.CreateButton(win, btnConfig)
	if err == nil {
		win.Childrens[id] = chWin

		return nil
	}
	return err
}

// Обработка событий мыши
func MouseEventHandler(ev winapi.Event) {
	mouseX = ev.Position.X
	mouseY = ev.Position.Y
	buttons := uint8(ev.SWin.Mbuttons)

	switch ev.Kind {
	case winapi.Move:
		//		log.Println("Mouse move ", ev.Position)
	case winapi.Press:
		log.Println("Mouse key press ", ev.Position, buttons)
	case winapi.Release:
		log.Println("Mouse key release ", ev.Position, buttons)
	case winapi.Leave:
		log.Println("Mouse lost focus ")
	case winapi.Enter:
		log.Println("Mouse enter focus ")

	}
}

func FrameEventHandler(ev winapi.Event) {
}

func getFileVersion() {
	Version = winapi.GetFileVersion()
}

// трей готов к работе
func onReady() {

	if len(img.ErrIco) > 0 {
		systray.SetIcon(img.ErrIco)
		systray.SetTooltip("WinGUI3 example")
	}
	systray.SetTitle("WinGUI3 systray")
	mQuit := systray.AddMenuItem("Quit", "Выход")
	mQuit.Enable()
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	systray.AddSeparator()
	mReconfig := systray.AddMenuItem("Reconfig", "Перечитать конфиг")
	mReconfig.Enable()
	go func() {
		for flag {
			<-mReconfig.ClickedCh
			log.Println("Reconfig")
		}
	}()

}

// Обработчик завершения трея
func onExit() {
	quit <- syscall.SIGTERM
	flag = false
}
