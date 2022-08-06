package filters

import (
	"github.com/davidbyttow/govips/v2/vips"
)

func Flip(img *vips.ImageRef, direction string) (*vips.ImageRef, error) {
	if direction == "h" {
		err := img.Flip(vips.DirectionHorizontal)
		if err != nil {
			return nil, err
		}
		return img, nil
	}
	err := img.Flip(vips.DirectionVertical)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func Pixlate(img *vips.ImageRef, factor float64) (*vips.ImageRef, error) {
	err := vips.Pixelate(img, factor)
	return img, err
}

func SmartCrop(img *vips.ImageRef, width int, height int, crop vips.Interesting) (*vips.ImageRef, error) {
	err := img.SmartCrop(width, height, crop)
	return img, err
}
