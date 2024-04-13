package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func InitApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Authorization,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,DELETE",
	}))

	app.Use(logger.New())

	return app
}

func RegisterRoutes(app *fiber.App) {
	db := app.Group("/db")
	v1 := db.Group("/v1")

	setupV1Router(&v1)
}
