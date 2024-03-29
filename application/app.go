package application

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"fyne.io/systray"
	"github.com/gbatanov/wingui3/winapi"
)

type Application struct {
	Win               *winapi.Window
	Version           string
	Quit              chan os.Signal
	Flag              bool
	MouseEventHandler func(winapi.Event)
	FrameEventHandler func(winapi.Event)
	KbEventHandler    func(winapi.Event)
	SystrayOnReady    func()
	ChildMutex        sync.Mutex
}

func AppCreate(Version string) *Application {
	var err error
	app := Application{}
	app.Quit = make(chan os.Signal)
	signal.Notify(app.Quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	app.Version = Version
	app.GetFileVersion()
	if Config.SysMenu == 0 {
		Config.Title = ""
	} else {
		Config.Title += ("      " + app.Version)
	}
	_, err = winapi.CreateNativeMainWindow(Config)
	app.Win = winapi.Wind

	if err != nil {
		panic(err)
	}
	app.Flag = true

	winapi.SetIcon(Config.SysMenu)

	return &app
}

func (app *Application) GetFileVersion() {
	vers := winapi.GetFileVersion()
	if vers != "" {
		app.Version = vers
	}
}

func (app *Application) Start() {

	app.Win.Config.MinSize.X = app.Win.Config.Size.X
	app.Win.Config.MaxSize.X = app.Win.Config.Size.X
	app.Win.Config.MinSize.Y = app.Win.Config.Size.Y
	app.Win.Config.MaxSize.Y = app.Win.Config.Size.Y

	app.eventHandler()

	winapi.SetWindowPos(app.Win.Hwnd,
		winapi.HWND_TOPMOST,
		int32(app.Win.Config.Position.X),
		int32(app.Win.Config.Position.Y),
		int32(app.Win.Config.Size.X),
		int32(app.Win.Config.Size.Y),
		winapi.SWP_NOMOVE)

	//systray (На Астре-Линукс не работает)

	var startSystray, endSystray func()
	if app.Win.Config.WithSystray {
		startSystray, endSystray = systray.RunWithExternalLoop(app.SystrayOnReady, app.onExit)
		startSystray()
	}

	winapi.Loop()

	if app.Win.Config.WithSystray {
		endSystray()
	}
	winapi.WinMap.Delete(app.Win.Hwnd)
	winapi.WinMap.Delete(0)

	ch := app.Win.GetChildren()
	for _, chWin := range ch {
		winapi.WinMap.Delete(chWin.Hwnd)
	}
	close(app.Win.Config.EventChan)

}

func (app *Application) MoveWindow(dx, dy int) {
	winapi.SetWindowPos(app.Win.Hwnd,
		winapi.HWND_TOPMOST,
		int32(app.Win.Config.Position.X+dx),
		int32(app.Win.Config.Position.Y+dy),
		int32(app.Win.Config.Size.X),
		int32(app.Win.Config.Size.Y),
		0)
	if runtime.GOOS == "windows" {
		app.Win.Config.Position.X += dx
		app.Win.Config.Position.Y += dy
	}
}

func (app *Application) eventHandler() {
	go func() {
		// Перехватчик исключений в горутине
		// Поскольку горутина закроется, сообщения от окна обрабатываться не будут
		// В частности, сообщение "Destroy" не придет в основной поток в этом случае
		defer func() {
			if val := recover(); val != nil {
				log.Println("eventHandler panic: ", val)
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
				case winapi.Keyboard:
					app.KbEventHandler(ev)
				}

			case <-app.Quit: // сообщение при закрытии трея
				app.Flag = false
				app.CloseWindow()
				return
			} //select
		} //for

		app.CloseWindow()
	}()
}

func (app Application) CloseWindow() {
	winapi.CloseWindow()
}

func (app *Application) AddLabel(title string) *Label {
	var lbl Label = Label{}
	lblConfig := labelConfig
	lblConfig.Title = title
	chWin, err := winapi.CreateLabel(app.Win, lblConfig)
	if err == nil {
		app.Win.Childrens = append(app.Win.Childrens, chWin)
		lbl = Label{Control{chWin}}
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
		app.Win.Childrens = append(app.Win.Childrens, chWin)
		//		id := len(*app.Win.Childrens)
		//		(*app.Win.Childrens)[id] = chWin
		btn = Button{Control{chWin}}
		return &btn
	}
	panic(err)
}

// Обработчик завершения трея
func (app *Application) onExit() {
	log.Println("Systray On exit")
	app.Quit <- syscall.SIGTERM
	app.Flag = false
}

func (app *Application) SysLog(level int, msg string) {
	if level != 0 {
		level = 1
	}
	winapi.SysLog(level, msg)
}

func (app *Application) GetChildren() []*winapi.Window {

	return app.Win.GetChildren()
}
