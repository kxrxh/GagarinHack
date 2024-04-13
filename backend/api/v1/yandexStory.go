package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	// Import resty into your code and refer it as `resty`.

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// yandexStory handles the processing of a Yandex story request.
//
// c: Context provided by Fiber framework.
// error: An error interface.
func yandexStory(c *fiber.Ctx) error {
	apiKey := viper.GetString("yandex.api_key")
	folderId := viper.GetString("yandex.folder_id")

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.TypeOfStory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type_of_story query param is required",
		})
	}

	if requestBody.TypeOfStory == "epitaph" {
		USER_PROMPT = USER_PROMPT_EPITAPH
		USER_PROMPT += "\nС именем: " + requestBody.HumanInfo.Name + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Дата смерти: " + requestBody.HumanInfo.DateOfDeath + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex + "" + "\n"
		for key, value := range requestBody.HumanInfo.Questions {
			USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
		}
	} else if requestBody.TypeOfStory == "biography" {
		USER_PROMPT = USER_PROMPT_BIOGRAPHY
		USER_PROMPT += "\nИмя: " + requestBody.HumanInfo.Name + "\n" + "Пол: " + requestBody.HumanInfo.Sex + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex + "" + "\n" + ""
		for key, value := range requestBody.HumanInfo.Questions {
			USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
		}
	}

	extraContext := []struct {
		Role string
		Text string
	}{}

	reqBodyBytes, _ := json.Marshal(getYandexRequestBody(folderId, 1024, 0.1, extraContext))

	req, _ := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(reqBodyBytes))

	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "PostmanRuntime/7.37.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	var data yandexResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	text := data.Result.Alternatives[0].Message.Text

	// Check if the length of the string is greater than 300
	if len(text) > 300 && requestBody.TypeOfStory == "epitaph" {
		trimmedText := text[:300]

		lastIndex := strings.LastIndexAny(trimmedText, ".!?…")

		if lastIndex != -1 {
			trimmedText = trimmedText[:lastIndex+1]
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"response": trimmedText,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"response": text,
		})
	}

}
