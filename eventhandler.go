package main

import (
	"log"
	"strings"

	"github.com/gbatanov/wingui3/winapi"
)

var mouseX, mouseY int = 0, 0

type Win struct {
	*winapi.Window
}

// Обработка событий мыши
func MouseEventHandler(ev winapi.Event) {
	w := ev.SWin
	w2 := Win{w}
	if w == nil {
		return
	}
	if strings.ToUpper(w.Config.Class) == "BUTTON" {
		if ev.Kind == winapi.Release {
			w2.HandleButton()
		}
		return
	}

	mouseX = ev.Position.X
	mouseY = ev.Position.Y

	switch ev.Kind {
	case winapi.Move:
		//		log.Println("Mouse move ", ev.Position)
	case winapi.Press:
		log.Println("Mouse key press ", ev.Position, ev.Mbuttons)
		log.Println("Mbuttons ", ev.SWin.Mbuttons)
		log.Println(ev.SWin.Config.Title)
	case winapi.Release:
		log.Println("Mouse key release ", ev.Position, ev.Mbuttons)
		log.Println("Mbuttons ", ev.SWin.Mbuttons)
		log.Println(ev.SWin.Config.Title)

	case winapi.Leave:
		log.Println("Mouse lost focus ")
	case winapi.Enter:
		log.Println("Mouse enter focus ")

	}

}

func FrameEventHandler(ev winapi.Event) {
}

// Обработчик нажатий кнопок
func (w Win) HandleButton() {

	switch w.Config.ID {
	case ID_BUTTON_1:
		log.Println("Click ", w.Config.Title)
		// И какие-то действия
		panic("Что-то пошло не тудысь!")

	case ID_BUTTON_2:
		log.Println(w.Config.Title)
		winapi.CloseWindow()
	}

}
