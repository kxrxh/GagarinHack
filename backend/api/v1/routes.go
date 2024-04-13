package v1

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutesV1 sets up the routes for the API version 1.
//
// v1: Pointer to the fiber.Router instance.
// No return value.
func SetupRoutesV1(v1 *fiber.Router) {
	api := (*v1).Group("/v1")

	// auth.SetupAuth(&v)
	completion := api.Group("/completion")
	completion.Post("/yandex", yandexCompletion)
	completion.Post("/yandex/questions", yandexQuestions)
	completion.Post("/yandex/questions/biography", yandexBioQuestions)
	completion.Post("/yandex/epitaph", yandexEpitaph)
	completion.Post("/yandex/biography", yandexBio)
	completion.Post("/yandex/biography/short", yandexBioShort)

	completion.Use(getAccessToken)
	completion.Post("/gigachat", gigachatCompletion)
	completion.Post("/gigachat/questions", gigachatGenerateQuestions)
	completion.Post("/gigachat/questions/biography", gigachatGenerateBioQuestions)
	completion.Post("/gigachat/epitaph", gigachatEpitaph)
	completion.Post("/gigachat/biography", gigachatBio)
	completion.Post("/gigachat/biography/short", gigachatBioShort)

	external := api.Group("/external")
	external.Post("/get-access-token", getApiAccessToken)

	// (*v1).Use(utils.Redirect)
}
