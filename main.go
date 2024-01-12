//go:generate go-winres make --file-version=v0.3.67.7 --product-version=git-tag
package main

import (
	_ "embed"
	"log"
	"runtime"
	"syscall"

	"fyne.io/systray"
	"github.com/gbatanov/wingui3/application"
	"github.com/gbatanov/wingui3/img"
	"github.com/gbatanov/wingui3/winapi"
)

var Version string = "v0.3.67"

var serverList []string = []string{"192.168.76.106", "192.168.76.80"}
var app *application.Application

// ---------------------------------------------------------------
func main() {

	defer func() {
		if val := recover(); val != nil {
			log.Println("Main thread error: ", val)
			syscall.Exit(1)
		}
	}()

	app = application.AppCreate(Version)
	app.MouseEventHandler = MouseEventHandler
	app.FrameEventHandler = FrameEventHandler

	defer winapi.WinMap.Delete(app.Win.Hwnd)
	defer winapi.WinMap.Delete(0)

	var posY int = 0
	// Label с текстом

	var Labels []*application.Label = make([]*application.Label, len(serverList))
	for id, title := range serverList {
		Labels[id] = app.AddLabel(title)
		posY = 10 + (Labels[id].Config.Size.Y)*(id)
		Labels[id].SetPos(int32(Labels[id].Config.Position.X), int32(posY), int32(Labels[id].Config.Size.X), int32(Labels[id].Config.Size.Y))
	}

	// Buttons
	posY = posY + Labels[0].Config.Size.Y + 10
	log.Println(posY)
	// Ok
	btnOk := app.AddButton(application.ID_BUTTON_1, "Ok")
	btnOk.SetPos(int32(btnOk.Config.Position.X+20), int32(posY), int32(40), int32(btnOk.Config.Size.Y))
	// Cancel
	btnCancel := app.AddButton(application.ID_BUTTON_2, "Cancel")
	btnCancel.SetPos(int32(btnOk.Config.Size.X+btnOk.Config.Position.X+20), int32(posY), int32(60), int32(btnOk.Config.Size.Y))

	for _, w2 := range app.Win.Childrens {
		defer winapi.WinMap.Delete(w2.Hwnd)
	}

	app.Win.Config.Size.Y = len(serverList)*Labels[0].Config.Size.Y + btnOk.Config.Size.Y + 30
	app.Win.Config.MinSize.Y = app.Win.Config.Size.Y
	app.Win.Config.MaxSize.Y = app.Win.Config.Size.Y

	//systray (На Астре-Линукс не работает)
	if runtime.GOOS == "windows" {
		go func() {
			systray.Run(onReady, onExit)
		}()
	}

	app.Start() // Здесь крутимся в цикле, пока не закроем окно

	close(app.Win.Config.EventChan) // Закрываем канал для завершения обработчика событий
	log.Println("Quit")

}
func CorrectSize() {

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
		for app.Flag {
			<-mReconfig.ClickedCh
			log.Println("Reconfig")
		}
	}()

}

// Обработчик завершения трея
func onExit() {
	app.Quit <- syscall.SIGTERM
	app.Flag = false
}
