//go:build linux
// +build linux

package img

import (
	_ "embed"
	"encoding/binary"

	"github.com/gbatanov/wingui3/gd2"
)

//go:embed check.ico
var OkIco []byte

//go:embed stop.ico
var ErrIco []byte

//go:embed  stop.png
var StopPng []byte

// если empty=true будет прозрачная иконка, как бы отсутствующая
func LoadIcon(empty bool) (int, []byte, error) {

	var ndata int = 0
	var data []byte

	icon := gd2.CreateFromPngPtr(StopPng)
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
			if empty {
				alpha = 0
			}
			binary.LittleEndian.PutUint32(data[i:], uint32(pixcolour&0x00ffffff)+uint32(alpha)<<24)
			i = i + 4
		}
	}
	icon.Destroy()
	return ndata, data, nil
}
