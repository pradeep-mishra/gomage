package api

import (
	"errors"
	"fmt"

	"gomage/filters"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type QueryParams struct {
	Flip string `query:"flip" validate:"omitempty,oneof=h v"`
}

var validate = validator.New()

func Startup() {
	vips.Startup(nil)
}

func Shutdown() {
	vips.Shutdown()
}

func checkError(err error, c *fiber.Ctx) {
	if err != nil {
		//fmt.Println("error:", err)
		c.Set("Content-Type", "application/json")
		c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}

func validateQuery(err error) []error {
	var errs []error
	for _, err := range err.(validator.ValidationErrors) {
		field := err.StructField()
		switch field {
		case "Flip":
			errs = append(errs, errors.New("value for flip filter is invalid. valid args is h or v"))
		}
	}
	return errs
}

func Optimize(c *fiber.Ctx) {
	imgid := c.Params("imgid")
	queryArgs := c.Context().QueryArgs().String()
	if len(queryArgs) >= 1 {
		fmt.Println(c.Context().QueryArgs())
	}
	img, err := loadImage(imgid)
	if err != nil {
		checkError(err, c)
		return
	}
	q := new(QueryParams)

	if err = c.QueryParser(q); err != nil {
		checkError(err, c)
		return
	}

	if err = validate.Struct(q); err != nil {
		errs := validateQuery(err)
		if len(errs) > 0 {
			checkError(errs[0], c)
			return
		}
	}

	qv := structs.Map(q)

	for filter, val := range qv {
		//fmt.Println("KV Pair: ", filter, val)
		str, _ := val.(string)
		//fmt.Println(len(str))
		if len(str) >= 1 {
			img = ApplyFilter(img, filter, string(str))
		}
	}
	ep := vips.NewDefaultJPEGExportParams()
	imgbytes, _, err := img.Export(ep)
	checkError(err, c)
	c.Set("Content-Type", "image/png")
	c.Write(imgbytes)
}

func loadImage(path string) (*vips.ImageRef, error) {
	filepath := "samples/" + path
	return vips.NewImageFromFile(filepath)

}

func ApplyFilter(img *vips.ImageRef, filter string, val string) *vips.ImageRef {
	switch filter {
	case "Flip":
		println("calling flip")
		return filters.Flip(img, val)
	default:
		return img
	}
}
