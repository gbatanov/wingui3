package application

import (
	"image"

	"github.com/gbatanov/wingui3/winapi"
)

const ID_BUTTON_1 = 101 // Ok
const ID_BUTTON_2 = 102 // Cancel

const COLOR_GREEN = 0x0011aa11
const COLOR_RED = 0x000000c8
const COLOR_YELLOW = 0x0000c8c8
const COLOR_GRAY_DE = 0x00dedede
const COLOR_GRAY_BC = 0x00bcbcbc
const COLOR_GRAY_AA = 0x00aaaaaa

// Конфиг основного окна приложения
var config = winapi.Config{
	Position:   image.Pt(-20, 20),
	MaxSize:    image.Pt(240, 240),
	MinSize:    image.Pt(240, 100),
	Size:       image.Pt(240, 30),
	Title:      "Wingui3",
	TextColor:  COLOR_GREEN,
	EventChan:  make(chan winapi.Event, 256),
	BorderSize: 0,
	Mode:       winapi.Windowed,
	BgColor:    COLOR_GRAY_DE,
	SysMenu:    2,
	Class:      "GsbWindow",
}
var labelConfig = winapi.Config{
	Class:      "Static",
	Title:      "Static",
	EventChan:  config.EventChan,
	Size:       image.Pt(int(config.Size.X-40), int(30)),
	MinSize:    config.MinSize,
	MaxSize:    config.MaxSize,
	Position:   image.Pt(int(18), int(15)),
	Mode:       winapi.Windowed,
	BorderSize: 0,
	TextColor:  COLOR_GREEN,
	FontSize:   28,
	BgColor:    COLOR_GRAY_DE, //config.BgColor,
}
var btnConfig = winapi.Config{
	Class:      "Button",
	Title:      "Ok",
	EventChan:  config.EventChan,
	Size:       image.Pt(int(40), int(25)),
	MinSize:    image.Pt(0, 0),
	MaxSize:    config.MaxSize,
	Position:   image.Pt(int(18), int(15)),
	Mode:       winapi.Windowed,
	BorderSize: 1,
	TextColor:  COLOR_GREEN,
	FontSize:   16,
	BgColor:    COLOR_GRAY_AA,
}
