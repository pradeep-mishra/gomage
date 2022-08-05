package optimize

import (
	"fmt"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func Startup() {
	vips.Startup(nil)
	//vips.LoggingSettings(logfunc ,vips.LogLevelError)
}

func Shutdown() {
	vips.Shutdown()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func Flip(path string) []byte {
	filepath := "samples/" + path
	img, err := vips.NewImageFromFile(filepath)
	checkError(err)

	err = img.Flip(vips.DirectionVertical)
	checkError(err)

	ep := vips.NewDefaultJPEGExportParams()
	imgbytes, _, err := img.Export(ep)
	checkError(err)
	return imgbytes

}
