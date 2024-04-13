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
func gigachatBioShort(c *fiber.Ctx) error {
	accessToken := c.Locals("sber_access_token").(string)
	baseUrl := viper.GetString("gigachat.baseUrl")
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
		Role    string
		Content string
	}{}
	reqBodyBytes, _ := json.Marshal(getGigachatRequestBody(1024, 0.6, extraContext))
	response, err := gigachatReq(accessToken, baseUrl, reqBodyBytes)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": strings.ReplaceAll(response.Choices[0].Message.Content, "\"", ""),
	})
}
