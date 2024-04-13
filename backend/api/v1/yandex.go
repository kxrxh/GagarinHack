package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
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

	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.SendString(string(body))
}
