//go:build linux
// +build linux

package img

import (
	_ "embed"
	"encoding/binary"

	gd "gitee.com/shirdonl/goGd"
)

//go:embed check.ico
var OkIco []byte

//go:embed stop.ico
var ErrIco []byte

//go:embed  stop.png
var StopPng []byte

func LoadIcon() (int, []byte, error) {

	var ndata int = 0
	var data []byte

	icon := gd.CreateFromPngPtr(StopPng)
	width := icon.Sx()
	height := icon.Sy()
	ndata = (width * height) + 2

	data = make([]byte, ndata*4)

	i := 0
	binary.LittleEndian.PutUint32(data[i:], uint32(width))
	i = i + 4
	binary.LittleEndian.PutUint32(data[i:], uint32(height))
	i = i + 4

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixcolour := icon.ColorAt(x, y) // ARGB

			alpha := 127 - uint8((pixcolour&0xff000000)>>24)
			if alpha == 127 {
				alpha = 255
			} else {
				alpha = alpha * 2
			}
			binary.LittleEndian.PutUint32(data[i:], uint32(pixcolour&0x00ffffff)+uint32(alpha)<<24)
			i = i + 4
		}
	}
	icon.Destroy()
	return ndata, data, nil
}
