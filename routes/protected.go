package routes

import (
	"github.com/akhil-is-watching/authservice/controllers"
	"github.com/akhil-is-watching/authservice/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupProtectedRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/securitycheck", middleware.DeserializeUser, controllers.SecurityCheck)
}
