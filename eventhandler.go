package main

import (
	"log"
	"strings"

	"github.com/gbatanov/wingui3/winapi"
)

// Обработка событий мыши
func MouseEventHandler(ev winapi.Event) {
	w := ev.SWin
	if w == nil {
		return
	}
	if strings.ToUpper(w.Config.Class) == "BUTTON" {
		if ev.Kind == winapi.Release {
			w.HandleButton(w, w.Config.ID)
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
