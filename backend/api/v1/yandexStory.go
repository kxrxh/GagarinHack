package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	// Import resty into your code and refer it as `resty`.

	"go.uber.org/zap"
)

func yandexStory(c *fiber.Ctx) error {
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
		USER_PROMPT += "\nС именем: " + requestBody.HumanInfo.Name + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Дата смерти: " + requestBody.HumanInfo.DateOfDeath + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex
		for key, value := range requestBody.HumanInfo.Questions {
			USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
		}
	} else if requestBody.TypeOfStory == "biography" {
		USER_PROMPT = USER_PROMPT_BIOGRAPHY
		USER_PROMPT += "\nИмя: " + requestBody.HumanInfo.Name + "\n" + "Пол: " + requestBody.HumanInfo.Sex + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex
		for key, value := range requestBody.HumanInfo.Questions {
			USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
		}
	}

	reqBodyBytes, _ := json.Marshal(yandexRequestBody)

	req, _ := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(reqBodyBytes))

	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")

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

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": data.Result.Alternatives[0].Message.Text,
	})
}
