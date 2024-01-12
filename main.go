//go:generate go-winres make --file-version=v0.3.66.7 --product-version=git-tag
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

var Version string = "v0.3.67" // Windows - подставится после генерации во время исполнения программы

var serverList []string = []string{"192.168.76.106", "192.168.76.80"}
var app *application.Application

// ---------------------------------------------------------------
func main() {
	//quit = make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	defer func() {
		// panic в горутинах здесь не обработается!
		// обработаются только паники из функций основного потока
		if val := recover(); val != nil {
			log.Println("main defer: ", val)
		}
	}()

	app, err := application.AppCreate(config)
	if err != nil || app == nil {
		return
	}
	app.MouseEventHandler = MouseEventHandler
	app.FrameEventHandler = FrameEventHandler

	defer winapi.WinMap.Delete(app.Win.Hwnd)
	defer winapi.WinMap.Delete(0)

	var id int = 0

	// Label с текстом
	for _, title := range serverList {
		labelConfig.Title = title
		app.AddLabel(labelConfig, id)
		id++
	}

	// Buttons
	// Ok
	btnConfig1 := btnConfig
	btnConfig1.ID = ID_BUTTON_1
	btnConfig1.Position.Y = 20 + (labelConfig.Size.Y)*(id)
	app.AddButton(btnConfig1, id)
	// Cancel
	id++
	btnConfig2 := btnConfig
	btnConfig2.Title = "Cancel"
	btnConfig2.ID = ID_BUTTON_2
	btnConfig2.Position.Y = btnConfig1.Position.Y
	btnConfig2.Position.X = btnConfig1.Position.X + btnConfig1.Size.X + 10
	btnConfig2.Size.X = 60
	app.AddButton(btnConfig2, id)

	if len(app.Win.Childrens) > 0 {
		for _, w2 := range app.Win.Childrens {
			defer winapi.WinMap.Delete(w2.Hwnd)
		}

		app.Win.Config.Size.Y = 2*labelConfig.Size.Y + btnConfig1.Size.Y + 30
		app.Win.Config.MinSize.Y = app.Win.Config.Size.Y
		app.Win.Config.MaxSize.Y = app.Win.Config.Size.Y
	}

	//systray (На Астре-Линукс не работает)
	if runtime.GOOS == "windows" {
		go func() {
			systray.Run(onReady, onExit)
		}()
	}

	winapi.SetIcon()

	app.Start() // Здесь крутимся в цикле, пока не закроем окно

	close(app.Win.Config.EventChan)
	log.Println("Quit")

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
