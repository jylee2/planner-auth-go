package main

import (
	"api-go/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true, // for the cookies
		AllowOrigins:     "http://localhost:3000, https://gofiber.net, https://jylee2.github.io/planner-react/",
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	routes.Setup(app)

	app.Listen(":8101")
}
