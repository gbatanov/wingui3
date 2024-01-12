package application

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gbatanov/wingui3/winapi"
)

type EventHandler func(winapi.Event)

type Application struct {
	Win               *winapi.Window
	Version           string
	Quit              chan os.Signal
	Flag              bool
	MouseEventHandler func(winapi.Event)
	FrameEventHandler func(winapi.Event)
}

func AppCreate(Version string) *Application {
	var err error
	app := Application{}
	app.Quit = make(chan os.Signal)
	signal.Notify(app.Quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	app.Version = Version
	app.GetFileVersion()
	config.Title += ("      " + app.Version)
	_, err = winapi.CreateNativeMainWindow(config)
	app.Win = winapi.Wind

	if err != nil {
		panic(err)
	}
	app.Flag = true
	winapi.SetIcon()
	return &app
}

func (app *Application) GetFileVersion() {
	vers := winapi.GetFileVersion()
	if vers != "" {
		app.Version = vers
	}
}

func (app *Application) Start() {
	app.eventHandler()
	/*
		winapi.SetWindowPos(app.Win.Hwnd,
			winapi.HWND_TOPMOST,
			int32(app.Win.Config.Position.X),
			int32(app.Win.Config.Position.Y),
			int32(app.Win.Config.Size.X),
			int32(app.Win.Config.Size.Y),
			winapi.SWP_NOMOVE)
	*/
	winapi.Loop()
}

func (app *Application) eventHandler() {
	go func() {
		// Перехватчик исключений в горутине
		// Поскольку горутина закроется, сообщения от окна обрабатываться не будут
		// В частности, сообщение "Destroy" не придет в основной поток в этом случае
		defer func() {
			if val := recover(); val != nil {
				log.Println("goroutine panic: ", val)
				winapi.CloseWindow()
			}
		}()

		for app.Flag {

			select {
			case ev, ok := <-app.Win.Config.EventChan:
				if !ok {
					// канал закрыт
					app.Flag = false
					break
				}

				switch ev.Source {
				case winapi.Mouse:
					app.MouseEventHandler(ev)
				case winapi.Frame:
					app.FrameEventHandler(ev)
				}

			case <-app.Quit: // сообщение при закрытии трея
				app.Flag = false
			} //select
		} //for

		winapi.SendMessage(app.Win.Hwnd, winapi.WM_CLOSE, 0, 0)
	}()
}

func (app *Application) AddLabel(title string) *Label {
	var lbl Label = Label{}
	lblConfig := labelConfig
	lblConfig.Title = title
	chWin, err := winapi.CreateLabel(app.Win, lblConfig)
	if err == nil {
		id := len(app.Win.Childrens)
		app.Win.Childrens[id] = chWin
		lbl = Label{chWin}
		return &lbl
	}
	panic(err)

}

func (app *Application) AddButton(ID int, title string) *Button {

	var btn Button = Button{}
	config := btnConfig
	config.ID = uintptr(ID)
	config.Title = title
	chWin, err := winapi.CreateButton(app.Win, config)
	if err == nil {
		id := len(app.Win.Childrens)
		app.Win.Childrens[id] = chWin
		btn = Button{chWin}
		return &btn
	}
	panic(err)
}
