package application

import (
	"github.com/gbatanov/wingui3/winapi"
)

type Control struct {
	*winapi.Window
}

type Button Control
type Label Control

func (b *Control) SetPos(x, y, w, h int32) {
	winapi.SetWindowPos(b.Hwnd, 0, x, y, w, h, 0)
}

func (b *Control) SetTitle(title string) {
	b.Config.Title = title
}

func (b *Button) SetPos(x, y, w, h int32) {
	winapi.SetWindowPos(b.Hwnd, -1, x, y, w, h, 2)
}

func (b *Button) SetTitle(title string) {
	b.Config.Title = title
}

func (b *Label) SetPos(x, y, w, h int32) {
	winapi.SetWindowPos(b.Hwnd, 0, x, y, w, h, 0)
}

func (b *Label) SetTitle(title string) {
	b.Config.Title = title
}
