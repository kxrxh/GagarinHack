package v1

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	// Import resty into your code and refer it as `resty`.
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type YandexRequest struct {
	ModelURI          string `json:"modelUri"`
	CompletionOptions struct {
		Stream      bool    `json:"stream"`
		Temperature float32 `json:"temperature"`
		MaxTokens   string  `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []struct {
		Role string `json:"role"`
		Text string `json:"text"`
	} `json:"messages"`
}

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
		ModelURI: "gpt://" + folderId + "/yandexgpt-lite",
		CompletionOptions: struct {
			Stream      bool    "json:\"stream\""
			Temperature float32 "json:\"temperature\""
			MaxTokens   string  "json:\"maxTokens\""
		}{
			Stream:      false,
			Temperature: 0.25,
			MaxTokens:   "1024",
		},
		Messages: []struct {
			Role string "json:\"role\""
			Text string "json:\"text\""
		}{
			{
				Role: "system",
				Text: SYSTEM_PROMT,
			},
			{
				Role: "user",
				Text: requestBody.RequestMessage,
			},
		},
	}

	client := resty.New()
	client.SetTimeout(30 * time.Second)

	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Api-Key %s", apiKey)).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post("https://llm.api.cloud.yandex.net/foundationModels/v1/completion")

	if err != nil {
		return err
	}

	zap.S().Debugln(fmt.Sprintf("response status: %v", resp.Status()))

	return c.SendString(resp.String())
}
