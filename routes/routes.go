package routes

import (
	"gomage/api"

	"github.com/gofiber/fiber/v2"
)

func V1Routes(app *fiber.App){
	v1 := app.Group("/v1")
	v1.Get("/optimize/:imgid", func(c *fiber.Ctx) error {
		return api.Optimize(c)
	})
}