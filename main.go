package main

import (
	"fmt"

	"api-go/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
  app := fiber.New()
  routes.Setup(app)

  fmt.Println("--------Hello, World!")
  app.Listen(":8000")
}