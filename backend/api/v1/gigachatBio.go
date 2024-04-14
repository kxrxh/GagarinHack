package v1

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// gigachatStory handles the story generation based on the type of story requested.
//
// It takes a Fiber context object as a parameter.
// Returns an error.
// gigachatReq sends a POST request to the Gigachat API and returns the response or an error.

func gigachatBio(c *fiber.Ctx) error {
	accessToken := c.Locals("sber_access_token").(string)
	baseUrl := viper.GetString("gigachat.baseUrl")
	var requestBody RequestBody
	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	USER_PROMPT = USER_PROMPT_BIOGRAPHY_START + requestBody.Previous + "\n"
	if requestBody.TypeOfStory == "youth" {
		USER_PROMPT += USER_PROMPT_BIOGRAPHY + " биографии человека в детстве и юношестве "
	} else if requestBody.TypeOfStory == "middle_age" {
		USER_PROMPT += USER_PROMPT_BIOGRAPHY + " биографии человека в средение года жизни "
	} else if requestBody.TypeOfStory == "old_age" {
		USER_PROMPT += USER_PROMPT_BIOGRAPHY + " биографии человека в последние года жизни "
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type_of_story query param is required",
		})
	}
	USER_PROMPT += "\nИмя: " + requestBody.HumanInfo.Name + "\n" + "Пол: " + requestBody.HumanInfo.Sex + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth
	for key, value := range requestBody.HumanInfo.Questions {
		USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
	}

	extraContext := []struct {
		Role    string
		Content string
	}{}
	reqBodyBytes, _ := json.Marshal(getGigachatRequestBody(1024, 0.6, extraContext))
	response, err := gigachatReq(accessToken, baseUrl, reqBodyBytes)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	USER_PROMPT = "\nСделай краткий заголовок к части биографии: " + response.Choices[0].Message.Content
	reqBodyBytesHeader, _ := json.Marshal(getGigachatRequestBody(1024, 0.6, extraContext))
	responseHeader, err := gigachatReq(accessToken, baseUrl, reqBodyBytesHeader)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": strings.ReplaceAll(response.Choices[0].Message.Content, "\"", ""),
		"header":   strings.ReplaceAll(responseHeader.Choices[0].Message.Content, "\"", ""),
	})

}
