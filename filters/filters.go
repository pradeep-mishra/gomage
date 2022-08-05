package filters

import (
	"fmt"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func Flip(img *vips.ImageRef, direction string) *vips.ImageRef {
	//fmt.Println("direction", direction)
	if direction == "h" {
		err := img.Flip(vips.DirectionHorizontal)
		checkError(err)
		return img
	}
	err := img.Flip(vips.DirectionVertical)
	checkError(err)
	return img
}
