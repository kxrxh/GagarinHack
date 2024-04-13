package v1

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func yandexBio(c *fiber.Ctx) error {
	apiKey := viper.GetString("yandex.api_key")
	folderId := viper.GetString("yandex.folder_id")

	var requestBody RequestBody
	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.TypeOfStory == "youth" {
		USER_PROMPT = USER_PROMPT_BIOGRAPHY + " биографии человека в детстве и юношестве "
	} else if requestBody.TypeOfStory == "middle_age" {
		USER_PROMPT = USER_PROMPT_BIOGRAPHY + " биографии человека в средение года жизни "
	} else if requestBody.TypeOfStory == "old_age" {
		USER_PROMPT = USER_PROMPT_BIOGRAPHY + " биографии человека в последние года жизни "
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type_of_story query param is required",
		})
	}
	USER_PROMPT += "\nИмя: " + requestBody.HumanInfo.Name + "\n" + "Пол: " + requestBody.HumanInfo.Sex + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex
	for key, value := range requestBody.HumanInfo.Questions {
		USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
	}

	extraContext := []struct {
		Role string
		Text string
	}{}

	reqBodyBytes, _ := json.Marshal(getYandexRequestBody(folderId, 1024, 0.1, extraContext))
	response, err := yandexReq(apiKey, reqBodyBytes)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": strings.ReplaceAll(response.Result.Alternatives[0].Message.Text, "\"", ""),
	})
}
