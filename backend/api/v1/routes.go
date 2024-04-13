package v1

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutesV1(v1 *fiber.Router) {
	api := (*v1).Group("/v1")

	// auth.SetupAuth(&v)
	completion := api.Group("/completion")
	completion.Post("/yandex", yandexCompletion)

	completion.Use(getAccessToken)
	completion.Post("/gigachat", gigachatCompletion)

	external := api.Group("/external")
	external.Post("/get-access-token", getApiAccessToken)
	// external.Use(getApiAccessToken)
	// external.Post("/get-access-token", login)
	// external.Post("/relative", connectPageToRelative)
	// external.Get("/individual-pages", connectPageToRelative)

	// (*v1).Use(utils.Redirect)
}
