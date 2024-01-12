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

func AppCreate(config winapi.Config) (*Application, error) {
	var err error
	app := Application{}
	app.Quit = make(chan os.Signal)
	signal.Notify(app.Quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	app.GetFileVersion()
	config.Title += ("      " + app.Version)
	_, err = winapi.CreateNativeMainWindow(config)
	app.Win = winapi.Wind

	if err != nil {
		return nil, err
	}
	app.Flag = true
	return &app, nil
}

func (app *Application) GetFileVersion() {
	app.Version = winapi.GetFileVersion()
}

func (app *Application) Start() {
	go app.eventHandler()
	winapi.SetWindowPos(app.Win.Hwnd,
		winapi.HWND_TOPMOST,
		int32(app.Win.Config.Position.X),
		int32(app.Win.Config.Position.Y),
		int32(app.Win.Config.Size.X),
		int32(app.Win.Config.Size.Y),
		winapi.SWP_NOMOVE)

	winapi.Loop()
}

func (app *Application) eventHandler() {

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

}

func (app *Application) AddLabel(lblConfig winapi.Config, id int) error {

	lblConfig.Position.Y = 10 + (lblConfig.Size.Y)*(id)
	chWin, err := winapi.CreateLabel(app.Win, lblConfig)
	if err == nil {
		app.Win.Childrens[id] = chWin

		return nil
	}
	return err

}

func (app *Application) AddButton(btnConfig winapi.Config, id int) error {

	chWin, err := winapi.CreateButton(app.Win, btnConfig)
	if err == nil {
		app.Win.Childrens[id] = chWin

		return nil
	}
	return err
}
