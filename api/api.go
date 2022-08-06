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
	}
	return errs
}

func mapError(field string) []error {
	var errs []error
	switch field {
	case "Flip":
		errs = append(errs, errors.New("value for flip filter is invalid. valid args is h or v"))
	case "Pixlate":
		errs = append(errs, errors.New("value for pixlate filter is invalid. valid args is a float between 0 to 100"))
	case "SmartCrop":
		errs = append(errs, errors.New("value for smartcrop filter is invalid. valid args is width,height,(intresting)"))
	}
	return errs
}

func Optimize(c *fiber.Ctx) error {
	imgid := c.Params("imgid")
	queryArgs := c.Context().QueryArgs().String()
	if len(queryArgs) >= 1 {
		fmt.Println(c.Context().QueryArgs())
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
		errs := extractError(err)
		if len(errs) > 0 {
			return checkError(errs[0], c)
		}
	}
	qv := structs.Map(q)
	if err = applyFilter(qv, img); err != nil {
		return checkError(err, c)
	}

	ep := vips.NewDefaultJPEGExportParams()
	imgbytes, _, err := img.Export(ep)
	checkError(err, c)
	c.Set("Content-Type", "image/png")
	c.Write(imgbytes)
	return nil
}

func applyFilter(queryMap map[string]interface{}, img *vips.ImageRef) error {
	for filter, val := range queryMap {
		if reflect.TypeOf(val).Kind() == reflect.Uint8 && fmt.Sprintf("%d", val.(uint8)) == "0" {
			fmt.Println("skipping filter:", filter)
			continue
		}
		if reflect.TypeOf(val).Kind() == reflect.String && val == "" {
			fmt.Println("skipping filter:", filter)
			continue
		}
		fmt.Println("KV Pair: ", filter, val, reflect.TypeOf(val))
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
	}
	return ""
}

func ApplyFilter(img *vips.ImageRef, filter string, val string) (*vips.ImageRef, error) {
	fmt.Println("filter:", filter, "is val:", val, reflect.TypeOf(val))
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
			return nil, errors.New("value for pixlate filter is invalid. valid args is a float")
		}
		fmt.Println("factor", factor)
		return filters.Pixlate(img, factor)
	case "SmartCrop":

		width, height, cropType, err := getSmartCropParams(val)
		if err != nil {
			return nil, err
		}
		fmt.Println("calling smartcrop", width, height, cropType)
		return filters.SmartCrop(img, width, height, cropType)
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
