package api

import (
	"gomage/optimize"

	"github.com/gofiber/fiber/v2"
)

func Optimize(c *fiber.Ctx) {
	imgid := c.Params("imgid")
	img := optimize.Flip(imgid)
	c.Set("Content-Type", "image/png")
	c.Write(img)
}
