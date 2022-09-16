package main

import (
	"github.com/akpor-kofi/auth/handling/rest"
	"github.com/akpor-kofi/auth/store/mongo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

const webPort = ":80"

func main() {
	log.Println("Starting authentication service...")

	userStore := mongo.NewUserStore()
	userHandler := rest.NewUserHandler(userStore)
	authHandler := rest.NewAuthHandler(userStore)

	app := fiber.New(fiber.Config{
		ErrorHandler: rest.GlobalHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://*, http://*",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.Post("/", userHandler.CreateUser)
	app.Post("/authenticate", authHandler.Authenticate)

	log.Fatal(app.Listen(webPort))
}
