package routes

import (
	"github.com/akhil-is-watching/authservice/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/healthcheck", controllers.HealthCheck)
}
