package api

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gomage/filters"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func Startup() {
	vips.Startup(nil)
}

func Shutdown() {
	vips.Shutdown()
}

func checkError(err error, c *fiber.Ctx) error {
	if err != nil {
		//fmt.Println("error:", err)
		c.Set("Content-Type", "application/json")
		c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}

func extractError(err error) []error {
	var errs []error
	for _, err := range err.(validator.ValidationErrors) {
		field := err.StructField()
		errs = mapError(field)
		if len(errs) == 0 {
			errs = append(errs, errors.New(err.Error()))
		}
	}
	return errs
}

func mapError(field string) []error {
	var errs []error
	switch field {
	case "Flip":
		errs = append(errs, errors.New("value for flip filter is invalid. valid param is h or v"))
	case "Pixlate":
		errs = append(errs, errors.New("value for pixlate filter is invalid. valid param is a number between 0 to 100"))
	case "SmartCrop":
		errs = append(errs, errors.New("value for smartcrop filter is invalid. valid param is width,height,(intresting)"))
	case "Rotate":
		errs = append(errs, errors.New("value for rotate filter is invalid. valid param is 0 90 180 270"))
	case "Zoom":
		errs = append(errs, errors.New("value for zoom filter is invalid. valid param is a number between 1 to 10"))

	}
	return errs
}

func Optimize(c *fiber.Ctx) error {
	imgid := c.Params("imgid")
	queryArgs := c.Context().QueryArgs().String()
	if len(queryArgs) >= 1 {
		fmt.Println("query => ", c.Context().QueryArgs())
	}
	img, err := loadImage(imgid)
	if err != nil {
		return checkError(err, c)
	}
	q := new(QueryParams)

	if err = c.QueryParser(q); err != nil {
		return checkError(err, c)
	}

	if err = validate.Struct(q); err != nil {
		fmt.Println("validation error:", err)
		errs := extractError(err)
		if len(errs) > 0 {
			return checkError(errs[0], c)
		}
	}
	qv := structs.Map(q)
	if err = applyFilter(qv, img); err != nil {
		return checkError(err, c)
	}

	exportParams, contentType := getFormat(q.Format)

	imgbytes, _, err := img.Export(exportParams)
	checkError(err, c)
	c.Set("Content-Type", contentType)
	c.Write(imgbytes)
	return nil
}

func getFormat(format string) (*vips.ExportParams, string) {
	switch format {
	case "png":
		return vips.NewDefaultPNGExportParams(), "image/png"
	case "webp":
		return vips.NewDefaultWEBPExportParams(), "image/webp"
	default:
		return vips.NewDefaultJPEGExportParams(), "image/jpeg"
	}
}

func applyFilter(queryMap map[string]interface{}, img *vips.ImageRef) error {
	for filter, val := range queryMap {
		if filter == "format" {
			continue
		}
		if reflect.TypeOf(val).Kind() == reflect.Uint8 && val.(uint8) == 0 {
			//fmt.Println("skipping filter:", filter)
			continue
		}
		if reflect.TypeOf(val).Kind() == reflect.Int && val.(int) == 0 {
			//fmt.Println("skipping filter:", filter)
			continue
		}

		if reflect.TypeOf(val).Kind() == reflect.Float64 && val.(float64) == float64(0) {
			//fmt.Println("skipping filter:", filter)
			continue
		}
		if reflect.TypeOf(val).Kind() == reflect.String && val == "" {
			//fmt.Println("skipping filter:", filter)
			continue
		}
		fmt.Println("query param: ", filter, val, reflect.TypeOf(val))
		_, err := ApplyFilter(img, filter, getValue(val))
		if err != nil {
			return err
		}
	}
	return nil
}

func loadImage(path string) (*vips.ImageRef, error) {
	filepath := "samples/" + path
	return vips.NewImageFromFile(filepath)

}

func getValue(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case float64:
		return fmt.Sprintf("%f", val.(float64))
	case uint8:
		return fmt.Sprintf("%d", val.(uint8))
	case int:
		return fmt.Sprintf("%d", val.(int))
	}
	return ""
}

func ApplyFilter(img *vips.ImageRef, filter string, val string) (*vips.ImageRef, error) {
	//fmt.Println("filter:", filter, "is val:", val, reflect.TypeOf(val))
	switch filter {
	case "Flip":
		fmt.Println("calling flip", val)
		return filters.Flip(img, val)
	case "Pixlate":
		fmt.Println("calling pixlate", val)
		factor, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, errors.New("value for pixlate filter is invalid. valid param is a float")
		}
		return filters.Pixlate(img, factor)
	case "SmartCrop":
		width, height, cropType, err := getSmartCropParams(val)
		if err != nil {
			return nil, err
		}
		fmt.Println("calling smartcrop width", width, "height", height, "croptype", cropType)
		return filters.SmartCrop(img, width, height, cropType)
	case "Sharpen":
		sigma, threshold, slope, err := getSharpenParams(val)
		if err != nil {
			return nil, err
		}
		fmt.Println("call sharpen sigma", sigma, "threshold", threshold, "slope", slope)
		return filters.Sharpen(img, sigma, threshold, slope)
	case "Rotate":
		angle, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New("value for rotate filter is invalid. valid param is an integer")
		}
		fmt.Println("call rotate angle", angle)
		return filters.Rotate(img, angle)
	case "Scale":
		fmt.Println("calling scale", val)
		scale, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, errors.New("value for scale filter is invalid. valid param is a number")
		}
		return filters.Scale(img, scale)
	case "Repeat":
		fmt.Println("calling repeat", val)
		x, y, err := getRepeatParams(val)
		if err != nil {
			return nil, err
		}
		if x > 10 || y > 10 || x < 1 || y < 1 {
			return nil, errors.New("value for repeat filter is invalid. valid param is a number between 1 to 10")
		}
		return filters.Repeat(img, x, y)
	case "Modulate":
		brightness, saturation, hue, err := getModulateParams(val)
		if err != nil {
			return nil, err
		}
		fmt.Println("call modulate brightness", brightness, "saturation", saturation, "hue", hue)
		return filters.Modulate(img, brightness, saturation, hue)
	case "Label":
		fmt.Println("calling label", val)
		label, err := getLabelParams(val)
		if err != nil {
			return nil, err
		}
		return filters.Label(img, label)
	case "Zoom":
		fmt.Println("calling zoom", val)
		zoomBy, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, errors.New("value for zoom filter is invalid. valid param is a number")
		}
		return filters.Zoom(img, zoomBy)
	default:
		return img, nil
	}
}

func getSmartCropParams(val string) (int, int, vips.Interesting, error) {
	params := strings.Split(val, ",")
	errorMsg := errors.New("value for smartcrop filter is invalid. valid args is width,height,(crop interest)")
	if len(params) < 2 {
		return 0, 0, vips.InterestingNone, errorMsg
	}
	width, err := strconv.Atoi(params[0])
	if err != nil {
		return 0, 0, vips.InterestingNone, errorMsg
	}
	height, err := strconv.Atoi(params[1])
	if err != nil {
		return 0, 0, vips.InterestingNone, errorMsg
	}

	cropType := "centre"

	if len(params) == 3 {
		cropType = params[2]
	}

	sc := &SmartCrop{
		Width:    width,
		Height:   height,
		CropType: cropType,
	}
	if err = validate.Struct(sc); err != nil {
		errs := mapError("SmartCrop")
		if len(errs) > 0 {
			return 0, 0, vips.InterestingNone, errs[0]
		}
	}

	return sc.Width, sc.Height, getSmartCropInterst(sc.CropType), nil
}

func getSmartCropInterst(intresting string) vips.Interesting {
	switch intresting {
	case "last":
		return vips.InterestingLast
	case "high":
		return vips.InterestingHigh
	case "entropy":
		return vips.InterestingEntropy
	case "attention":
		return vips.InterestingAttention
	case "all":
		return vips.InterestingAll
	case "low":
		return vips.InterestingLow
	case "none":
		return vips.InterestingNone
	default:
		return vips.InterestingCentre
	}
}

func getSharpenParams(val string) (float64, float64, float64, error) {
	params := strings.Split(val, ",")
	errorMsg := errors.New("value for sharpen filter is invalid. valid param is sigma,threshold,slope")
	if len(params) < 3 {
		return 0, 0, 0, errorMsg
	}
	sigma, err := strconv.ParseFloat(params[0], 64)
	if err != nil {
		return 0, 0, 0, errorMsg
	}
	threshold, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return 0, 0, 0, errorMsg
	}
	slope, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return 0, 0, 0, errorMsg
	}
	return sigma, threshold, slope, nil
}

func getRepeatParams(val string) (int, int, error) {
	params := strings.Split(val, ",")
	errorMsg := errors.New("value for replicate filter is invalid. valid param is x,y")
	if len(params) < 2 {
		return 0, 0, errorMsg
	}
	x, err := strconv.Atoi(params[0])
	if err != nil {
		return 0, 0, errorMsg
	}
	y, err := strconv.Atoi(params[1])
	if err != nil {
		return 0, 0, errorMsg
	}
	return x, y, nil
}

func getModulateParams(val string) (float64, float64, float64, error) {
	params := strings.Split(val, ",")
	errorMsg := errors.New("value for modulate filter is invalid. valid param is sigma,threshold,slope in float")
	if len(params) < 3 {
		return 0, 0, 0, errorMsg
	}
	brightness, err := strconv.ParseFloat(params[0], 64)
	if err != nil {
		return 0, 0, 0, errorMsg
	}
	saturation, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return 0, 0, 0, errorMsg
	}
	hue, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return 0, 0, 0, errorMsg
	}
	return brightness, saturation, hue, nil
}

func getLabelParams(val string) (*vips.LabelParams, error) {
	params := strings.Split(val, ",")
	errorMsg := errors.New("value for label filter is invalid. valid param is text,font,width,height,x,y,opacity,color(r:g:b)")

	labelParams := &vips.LabelParams{}

	if len(params) < 8 {
		return nil, errorMsg
	}
	labelParams.Text = params[0]
	labelParams.Font = params[1]
	width, err := strconv.ParseFloat(params[2], 64)
	if err != nil {
		return nil, errorMsg
	}
	labelParams.Width = vips.Scalar{Value: width, Relative: false}
	height, err := strconv.ParseFloat(params[3], 64)
	if err != nil {
		return nil, errorMsg
	}
	labelParams.Height = vips.Scalar{Value: height, Relative: false}
	x, err := strconv.ParseFloat(params[4], 64)
	if err != nil {
		return nil, errorMsg
	}
	labelParams.OffsetX = vips.Scalar{Value: x, Relative: false}
	y, err := strconv.ParseFloat(params[5], 64)
	if err != nil {
		return nil, errorMsg
	}
	labelParams.OffsetY = vips.Scalar{Value: y, Relative: false}
	opacity, err := strconv.ParseFloat(params[6], 32)
	if err != nil {
		return nil, errorMsg
	}
	labelParams.Opacity = float32(opacity)

	colors := strings.Split(params[7], ":")
	fmt.Println("colors", colors)
	errorMsg = errors.New("value for color in label filter is invalid. valid param is 0-255:0-255:0-255 (R:G:B)")
	if len(colors) < 3 {
		return nil, errorMsg
	}
	r, err := strconv.ParseInt(colors[0], 0, 16)
	if err != nil {
		return nil, errorMsg
	}
	g, err := strconv.ParseInt(colors[1], 0, 16)
	if err != nil {
		return nil, errorMsg
	}
	b, err := strconv.ParseInt(colors[2], 0, 16)
	if err != nil {
		return nil, errorMsg
	}
	color := &vips.Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}
	labelParams.Color = *color

	return labelParams, nil
}
