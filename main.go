//go:generate go-winres make --file-version=v0.3.88.10 --product-version=git-tag
package main

import (
	_ "embed"
	"log"
	"runtime"
	"strings"
	"syscall"

	"fyne.io/systray"
	"github.com/gbatanov/wingui3/application"
	"github.com/gbatanov/wingui3/img"
	"github.com/gbatanov/wingui3/winapi"
)

var Version string = "v0.3.88"

var serverList []string = []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"}
var app *application.Application

// ---------------------------------------------------------------
func main() {

	defer func() {
		if val := recover(); val != nil {
			log.Println("Main thread error: ", val)
			syscall.Exit(1)
		}
	}()

	application.Config.SysMenu = 1
	app = application.AppCreate(Version)
	app.MouseEventHandler = MouseEventHandler
	app.FrameEventHandler = FrameEventHandler
	app.KbEventHandler = KbEventHandler
	app.SystrayOnReady = onReady

	defer winapi.WinMap.Delete(app.Win.Hwnd)
	defer winapi.WinMap.Delete(0)

	var posY int = 10

	// Label с текстом
	var Labels []*application.Label = make([]*application.Label, len(serverList))
	for id, title := range serverList {
		Labels[id] = app.AddLabel(title)
		Labels[id].SetPos(int32(Labels[id].Config.Position.X), int32(posY), int32(Labels[id].Config.Size.X), int32(Labels[id].Config.Size.Y))
		posY += Labels[0].Config.Size.Y

	}
	app.Win.Config.Size.Y += posY

	// Buttons
	posY += 10
	// Ok
	btnOk := app.AddButton(application.ID_BUTTON_1, "Ok")
	btnOk.SetPos(int32(btnOk.Config.Position.X+20), int32(posY), int32(40), int32(btnOk.Config.Size.Y))

	// Cancel
	btnCancel := app.AddButton(application.ID_BUTTON_2, "Cancel")
	btnCancel.SetPos(int32(btnOk.Config.Size.X+btnOk.Config.Position.X+20), int32(posY), int32(60), int32(btnOk.Config.Size.Y))
	app.Win.Config.Size.Y += btnCancel.Config.Size.Y

	if runtime.GOOS == "windows" {
		if application.Config.SysMenu == 0 {
			app.Win.Config.Size.Y -= 10
		} else {
			app.Win.Config.Size.Y += 20
		}
	} else {
		app.Win.Config.Size.Y -= 20
	}

	ch := app.GetChildren()
	for _, w2 := range ch {
		defer winapi.WinMap.Delete(w2.Hwnd)
	}

	app.Start() // Здесь крутимся в цикле, пока не закроем окно

	close(app.Win.Config.EventChan) // Закрываем канал для завершения обработчика событий
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
		for app.Flag {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func KbEventHandler(ev winapi.Event) {
	w := ev.SWin
	if w == nil {
		return
	}

	switch ev.Kind {
	case winapi.Press:
	case winapi.Release:
		if strings.ToUpper(ev.Name) == "Q" {
			winapi.CloseWindow()
		}
	}
}

var x, y int = 0, 0

// Обработка событий мыши
func MouseEventHandler(ev winapi.Event) {
	w := ev.SWin
	if w == nil {
		return
	}

	if strings.ToUpper(w.Config.Class) == "BUTTON" {
		HandleButton(ev)
		return
	}

	switch ev.Kind {
	case winapi.Move:
		//		log.Println("Mouse move ", ev.Position)
		if (ev.SWin.Mbuttons & winapi.ButtonPrimary) != 0 {
			// В бесшапочном режиме двигаем окно
			if application.Config.SysMenu == 0 {
				x1, y1, _ := ev.SWin.WinTranslateCoordinates(ev.Position.X, ev.Position.Y)
				dx := x1 - x
				dy := y1 - y
				app.MoveWindow(dx, dy)
				x = x1
				y = y1
			}
		}
	case winapi.Press:
		log.Println("Mouse key press ", ev.Position, ev.Mbuttons)
		if (ev.SWin.Mbuttons & winapi.ButtonPrimary) != 0 {
			x, y, _ = ev.SWin.WinTranslateCoordinates(ev.Position.X, ev.Position.Y)
			log.Printf("Mouse key press posX; %d posY: %d x: %d y: %d", ev.Position.X, ev.Position.Y, x, y)
		}
	//	log.Println("Mbuttons ", ev.SWin.Mbuttons)
	//	log.Println(ev.SWin.Config.Title)
	case winapi.Release:
		log.Println("Mouse key release ", ev.Position, ev.Mbuttons)
		if ev.Mbuttons == winapi.ButtonSecondary {
			if application.Config.SysMenu == 0 {
				winapi.CloseWindow()
			}
		}
		if (ev.SWin.Mbuttons & winapi.ButtonPrimary) == 0 {
			x = 0
			y = 0
		}
	//	log.Println("Mbuttons ", ev.SWin.Mbuttons)
	//	log.Println(ev.SWin.Config.Title)

	case winapi.Leave:
	//	log.Println("Mouse lost focus ")
	case winapi.Enter:
		//	log.Println("Mouse enter focus ")

	}
}

// Обработчик событий от окна, отличных от кнопок и мыши
func FrameEventHandler(ev winapi.Event) {
	switch ev.Kind {
	case winapi.Destroy:
		// Большого смысла в обработке этого события в основном потоке нет, чисто информативно.
		// Не придет, если работа завершается при панике в горутине приема событий от окна.
		log.Println("Window destroy")
	}
}

// Обработчик нажатий кнопок-объектов
func HandleButton(ev winapi.Event) {
	//	log.Println("Click ", w.Config.Title)
	w := ev.SWin
	if w == nil {
		return
	}

	if ev.Kind == winapi.Release {
		switch w.Config.ID {
		case application.ID_BUTTON_1:
			// какие-то действия
		//	panic("Что-то пошло не так!") // имитация сбоя в работе

		case application.ID_BUTTON_2:
			// какие-то действия
			winapi.CloseWindow() // выход из программы
		}
	}
}
