package main

import (
	"log"

	"github.com/akhil-is-watching/authservice/config"
	"github.com/akhil-is-watching/authservice/routes"
	"github.com/akhil-is-watching/authservice/storage"
	"github.com/gofiber/fiber/v2"
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	storage.ConnectDB(&config)
}
func main() {

	app := fiber.New()
	routes.SetupRoutes(app)
	routes.SetupPublicRoutes(app)
	routes.SetupProtectedRoutes(app)

	app.Listen(":3000")
}
