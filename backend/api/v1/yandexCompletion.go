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

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func yandexCompletion(c *fiber.Ctx) error {
	apiKey := viper.GetString("yandex.api_key")
	folderId := viper.GetString("yandex.folder_id")

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.RequestMessage == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request_message query param is required",
		})
	}

	reqBody := YandexRequest{
		ModelURI: "gpt://" + folderId + "/yandexgpt",
		CompletionOptions: struct {
			Stream      bool    "json:\"stream\""
			Temperature float32 "json:\"temperature\""
			MaxTokens   string  "json:\"maxTokens\""
		}{
			Stream:      false,
			Temperature: 0.1,
			MaxTokens:   "1024",
		},
		Messages: []struct {
			Role string "json:\"role\""
			Text string "json:\"text\""
		}{
			{
				Role: "system",
				Text: SYSTEM_PROMPT,
			},
			{
				Role: "user",
				Text: requestBody.RequestMessage,
			},
		},
	}

	reqBodyBytes, _ := json.Marshal(reqBody)

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

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": data.Result.Alternatives[0].Message.Text,
	})
}
