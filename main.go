package main

import (
	"fmt"

	"api-go/database"
	"api-go/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
  database.Connect()

  app := fiber.New()
  routes.Setup(app)

  fmt.Println("--------Hello, World!")
  app.Listen(":8000")
}