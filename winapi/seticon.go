//go:build linux
// +build linux

package winapi

import (
	"log"
	"os"

	gd "gitee.com/shirdonl/goGd"
)

func LoadIcon(filename string) (int, []uint32, error) {

	var ndata int = 0
	var data []uint32 = make([]uint32, 0)

	iconfile, err := os.OpenFile("/home/user/work/src/wingui3/img/stop.png", os.O_RDONLY, 0)
	if err != nil {
		log.Println(err.Error())
		return ndata, data, err
	}
	iconfile.Close()

	icon := gd.CreateFromPng(filename)
	width := icon.Sx()
	height := icon.Sy()
	ndata = (width * height) + 2
	data = make([]uint32, ndata)

	i := 0
	data[i] = uint32(width)
	i++
	data[i] = uint32(height)
	i++

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixcolour := icon.ColorAt(x, y) // ARGB

			alpha := 127 - uint8((pixcolour&0xff000000)>>24)
			if alpha == 127 {
				alpha = 255
			} else {
				alpha = alpha * 2
			}
			data[i] = uint32(pixcolour&0x00ffffff) + uint32(alpha)<<24
			i++
		}
	}
	icon.Destroy()
	return ndata, data, nil
}
