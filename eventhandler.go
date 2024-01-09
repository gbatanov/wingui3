package main

import (
	"log"

	"github.com/gbatanov/wingui3/winapi"
)

// Обработка событий мыши
func MouseEventHandler(ev winapi.Event) {
	mouseX = ev.Position.X
	mouseY = ev.Position.Y

	switch ev.Kind {
	case winapi.Move:
		//		log.Println("Mouse move ", ev.Position)
	case winapi.Press:
		log.Println("Mouse key press ", ev.Position, ev.Mbuttons)
		log.Println("Pressed ", ev.SWin.Mbuttons)
	case winapi.Release:
		log.Println("Mouse key release ", ev.Position, ev.Mbuttons)
		log.Println("Pressed ", ev.SWin.Mbuttons)
	case winapi.Leave:
		log.Println("Mouse lost focus ")
	case winapi.Enter:
		log.Println("Mouse enter focus ")

	}

}

func FrameEventHandler(ev winapi.Event) {
}
