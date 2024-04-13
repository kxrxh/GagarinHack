package api

import (
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	v1 "github.com/gagarin/backend/api/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"go.uber.org/zap"
	// "github.com/gagarin/backend/api/auth"
)

type Api struct {
	appAddress string
	appPort    string
	app        *fiber.App
}

// CreateApi creates a new API instance with the given address and port
func CreateApi(address, port string) *Api {
	if port == "" {
		zap.S().Panic("app port is not provided")
	}

	app := fiber.New()
	// prom := prometheusMiddleware(app)
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Authorization,Content-Type,Origin,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "https://github.com",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	// app.Use(prom.Middleware)

	return &Api{appAddress: address, appPort: port, app: app}
}

// prometheusMiddleware generates a FiberPrometheus middleware for the given app.
func prometheusMiddleware(app *fiber.App) *fiberprometheus.FiberPrometheus {
	// TODO: Make sure middleware works
	prometheus := fiberprometheus.New("prometheus-service")
	prometheus.RegisterAt(app, "/metrics")
	return prometheus
}

// ConfigureApp sets up the routes for the API
func (api *Api) ConfigureApp() *Api {
	api.app.Get("/", func(c *fiber.Ctx) error {
		zap.S().Debugln("GET /api")
		return c.JSON(fiber.Map{
			"message": "Gagarin Hack API",
		})
	})

	apiGroup := api.app.Group("/api")

	apiGroup.Use(healthcheck.New())
	// auth.SetupAuth(&apiGroup)
	v1.SetupRoutesV1(&apiGroup)
	return api
}

// Run starts the API server on the given address and port
func (api *Api) Run() {
	zap.S().Debugln(fmt.Sprintf("Starting server on %s:%s", api.appAddress, api.appPort))
	api.app.Listen(fmt.Sprintf("%s:%s", api.appAddress, api.appPort))
}
