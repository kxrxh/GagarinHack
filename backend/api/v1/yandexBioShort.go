package v1

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func yandexBioShort(c *fiber.Ctx) error {
	apiKey := viper.GetString("yandex.api_key")
	folderId := viper.GetString("yandex.folder_id")

	var requestBody ShortBody
	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	USER_PROMPT = USER_PROMPT_ENDING + " " + requestBody.HumanInfo.Name + "\n" + "Пол: " + requestBody.HumanInfo.Sex + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth
	USER_PROMPT += "\nЧАСТЬ БИОГРАФИИ №1: " + requestBody.PartYoung
	USER_PROMPT += "\nЧАСТЬ БИОГРАФИИ №2: " + requestBody.PartMiddle
	USER_PROMPT += "\nЧАСТЬ БИОГРАФИИ №3: " + requestBody.PartOld

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
