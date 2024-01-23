package application

import (
	"github.com/gbatanov/wingui3/winapi"
)

type Control struct {
	*winapi.Window
}

type Button struct{ Control }
type Label struct{ Control }

func (c *Control) SetPos(x, y, w, h int32) {
	winapi.SetWindowPos(c.Hwnd, 0, x, y, w, h, 0)
}

func (c *Control) SetTitle(title string) {
	c.Config.Title = title
}
