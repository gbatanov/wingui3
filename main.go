//go:generate go-winres make --file-version=v0.3.64.7 --product-version=git-tag
package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"fyne.io/systray"
	"github.com/gbatanov/wingui3/img"
	"github.com/gbatanov/wingui3/winapi"
)

var Version string = "v0.3.64" // Windows - подставится после генерации во время исполнения программы

var serverList []string = []string{"192.168.76.106", "192.168.76.80"}

var quit chan os.Signal
var flag = true

// ---------------------------------------------------------------
func main() {
	quit = make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	defer func() {
		// panic в горутинах здесь не обработается!
		// обработаются только паники из функций основного потока
		if val := recover(); val != nil {
			log.Println("main defer: ", val)
		}
	}()

	getFileVersion()
	config.Title += ("      " + Version)

	win, err := winapi.CreateNativeMainWindow(config)
	if err == nil {
		mainWin := Win{win}
		// Обработчик событий (события от дочерних элементов приходят сюда же)
		go func() {
			// Перехватчик исключений в горутине+
			defer func() {
				if val := recover(); val != nil {
					log.Println("goroutine panic: ", val)
					winapi.CloseWindow()
					flag = false
				}
			}()

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

			winapi.SendMessage(mainWin.Hwnd, winapi.WM_CLOSE, 0, 0)
		}()

		defer winapi.WinMap.Delete(mainWin.Hwnd)
		defer winapi.WinMap.Delete(0)

		var id int = 0

		// Label с текстом
		for _, title := range serverList {
			labelConfig.Title = title
			mainWin.AddLabel(labelConfig, id)
			id++
		}

		// Buttons
		// Ok
		btnConfig1 := btnConfig
		btnConfig1.ID = ID_BUTTON_1
		btnConfig1.Position.Y = 20 + (labelConfig.Size.Y)*(id)
		mainWin.AddButton(btnConfig1, id)
		// Cancel
		id++
		btnConfig2 := btnConfig
		btnConfig2.Title = "Cancel"
		btnConfig2.ID = ID_BUTTON_2
		btnConfig2.Position.Y = btnConfig1.Position.Y
		btnConfig2.Position.X = btnConfig1.Position.X + btnConfig1.Size.X + 10
		btnConfig2.Size.X = 60
		mainWin.AddButton(btnConfig2, id)

		if len(mainWin.Childrens) > 0 {
			for _, w2 := range mainWin.Childrens {
				defer winapi.WinMap.Delete(w2.Hwnd)
			}

			mainWin.Config.Size.Y = 2*labelConfig.Size.Y + btnConfig1.Size.Y + 30
			mainWin.Config.MinSize.Y = win.Config.Size.Y
			mainWin.Config.MaxSize.Y = win.Config.Size.Y
		}

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

func (w *Win) AddLabel(lblConfig winapi.Config, id int) error {

	lblConfig.Position.Y = 10 + (lblConfig.Size.Y)*(id)
	chWin, err := winapi.CreateLabel(w.Window, lblConfig)
	if err == nil {
		w.Childrens[id] = chWin

		return nil
	}
	return err
}

func (w *Win) AddButton(btnConfig winapi.Config, id int) error {

	chWin, err := winapi.CreateButton(w.Window, btnConfig)
	if err == nil {
		w.Childrens[id] = chWin

		return nil
	}
	return err
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
