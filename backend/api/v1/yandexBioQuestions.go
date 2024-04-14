package v1

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// yandexQuestions handles the processing of a Yandex questions request.
//
// c: Context provided by Fiber framework.
// error: An error interface.
func yandexBioQuestions(c *fiber.Ctx) error {
	apiKey := viper.GetString("yandex.api_key")
	folderId := viper.GetString("yandex.folder_id")

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.HumanInfo.Sex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sex query param is required",
		})
	}

	if requestBody.HumanInfo.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name query param is required",
		})
	}

	if requestBody.HumanInfo.DateOfBirth == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "date_of_birth query param is required",
		})
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

	reqBodyBytes, _ := json.Marshal(getYandexRequestBody(folderId, 1024, 0.2, extraContext))
	resp, err := yandexReq(apiKey, reqBodyBytes)

	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	questionRegex := regexp.MustCompile(`^\d*\.?\s*([А-Я].*\?)$`)
	numberDotRegex := regexp.MustCompile(`^\d+\.\s*`)

	var questions []string
	for _, line := range strings.Split(resp.Result.Alternatives[0].Message.Text, "\n") {
		if questionRegex.MatchString(line) {
			cleanLine := strings.ReplaceAll(numberDotRegex.ReplaceAllString(line, ""), "\"", "") // Remove the number and dot at the beginning
			questions = append(questions, cleanLine)
		}
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": questions,
	})

}
