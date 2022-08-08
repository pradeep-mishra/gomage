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

func Sharpen(img *vips.ImageRef, sigma float64, threshold float64, slope float64) (*vips.ImageRef, error) {
	err := img.Sharpen(sigma, threshold, slope)
	return img, err
}

func Rotate(img *vips.ImageRef, angle int) (*vips.ImageRef, error) {
	var err error = nil
	switch angle {
	case 0:
		err = img.Rotate(vips.Angle0)
	case 90:
		err = img.Rotate(vips.Angle90)
	case 180:
		err = img.Rotate(vips.Angle180)
	case 270:
		err = img.Rotate(vips.Angle270)
	}
	return img, err
}

func Scale(img *vips.ImageRef, scale float64) (*vips.ImageRef, error) {
	if scale == 0 {
		return img, nil
	}
	err := img.Resize(scale, vips.KernelAuto)
	return img, err
}

func Repeat(img *vips.ImageRef, x int, y int) (*vips.ImageRef, error) {
	err := img.Replicate(x, y)
	return img, err
}

func Modulate(img *vips.ImageRef, brightness float64, saturation float64, hue float64) (*vips.ImageRef, error) {
	err := img.Modulate(brightness, saturation, hue)
	return img, err
}

func Label(img *vips.ImageRef, label *vips.LabelParams) (*vips.ImageRef, error) {
	err := img.Label(label)
	return img, err
}

func Zoom(img *vips.ImageRef, zoomBy int) (*vips.ImageRef, error) {
	err := img.Zoom(zoomBy, zoomBy)
	return img, err
}
