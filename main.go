package main

import (
	"log"

	"gomage/routes"

	"github.com/gofiber/fiber/v2"
)

func main(){
	app := fiber.New()
	routes.V1Routes(app)
	log.Fatal(app.Listen(":3300"))
}


