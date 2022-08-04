package api

import (
	"github.com/gofiber/fiber/v2"
)

func Optimize(c *fiber.Ctx) error {
	imgid := c.Params("imgid")
	data := fiber.Map{
		"status": 200,
		"imageId": imgid,
		"notice": "Api called successfully",
	}
	return c.Status(fiber.StatusOK).JSON(data)
}