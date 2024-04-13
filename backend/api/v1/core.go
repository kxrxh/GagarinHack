package v1

import (
	"github.com/gofiber/fiber/v2"

	// "github.com/gagarin/backend/api/auth"
	"github.com/gagarin/backend/utils"
)

func SetupRoutesV1(v1 *fiber.Router) {
	api := (*v1).Group("/v1")

	// auth.SetupAuth(&v)
	completion := api.Group("/completion")
	completion.Post("/yandex", yandexCompletion)

	completion.Use("/gigachat", getAccessToken)
	completion.Post("/gigachat", gigachatCompletion)

	(*v1).Use(utils.Redirect)
}
