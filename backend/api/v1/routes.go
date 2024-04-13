package v1

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutesV1(v1 *fiber.Router) {
	api := (*v1).Group("/v1")

	// auth.SetupAuth(&v)
	completion := api.Group("/completion")
	completion.Post("/yandex", yandexCompletion)
	completion.Post("/yandex/questions", yandexQuestions)
	completion.Post("/yandex/story", yandexStory)

	completion.Use(getAccessToken)
	completion.Post("/gigachat", gigachatCompletion)
	completion.Post("/gigachat/questions", gigachatGenerateQuestions)
	completion.Post("/gigachat/story", gigachatStory)

	external := api.Group("/external")
	external.Post("/get-access-token", getApiAccessToken)

	// (*v1).Use(utils.Redirect)
}
