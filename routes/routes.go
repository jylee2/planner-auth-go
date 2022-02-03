package routes

import (
	"api-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
  app.Get("/api/user", controllers.GetUserFromCookie)

  app.Post("/api/register", controllers.Register)
  app.Post("/api/login", controllers.Login)
  app.Post("/api/logout", controllers.Logout)
}