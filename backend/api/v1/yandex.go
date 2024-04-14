package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type YandexRequest struct {
	ModelURI          string `json:"modelUri"`
	CompletionOptions struct {
		Stream      bool    `json:"stream"`
		Temperature float32 `json:"temperature"`
		MaxTokens   uint    `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []struct {
		Role string `json:"role"`
		Text string `json:"text"`
	} `json:"messages"`
}

// getYandexRequestBody generates a Yandex request body based on the provided folder ID, max tokens, and temperature.
//
// folderId: The ID of the folder.
// maxTokens: The maximum number of tokens.
// temperature: The temperature for the completion.
// Returns a pointer to a YandexRequest.
func getYandexRequestBody(folderId string, maxTokens uint, temperature float32, extraContext []struct {
	Role string
	Text string
}) *YandexRequest {
	if folderId == "" || maxTokens == 0 || temperature < 0 {
		return nil
	}

	messages := []struct {
		Role string `json:"role"`
		Text string `json:"text"`
	}{
		{
			Role: "system",
			Text: SYSTEM_PROMPT,
		},
		{
			Role: "user",
			Text: USER_PROMPT,
		},
	}

	for _, ctx := range extraContext {
		messages = append(messages, struct {
			Role string `json:"role"`
			Text string `json:"text"`
		}{
			Role: ctx.Role,
			Text: ctx.Text,
		})
	}

	return &YandexRequest{
		ModelURI: fmt.Sprintf("gpt://%s/yandexgpt-lite", folderId),
		CompletionOptions: struct {
			Stream      bool    `json:"stream"`
			Temperature float32 `json:"temperature"`
			MaxTokens   uint    `json:"maxTokens"`
		}{
			Stream:      false,
			Temperature: temperature,
			MaxTokens:   maxTokens,
		},
		Messages: messages,
	}
}

// yandexReq sends a POST request to the Yandex API and returns the response or an error.
func yandexReq(apiKey string, reqBody []byte) (*yandexResponseData, error) {
	req, err := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
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
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data yandexResponseData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

type yandexResult struct {
	Alternatives []yandexAlternative `json:"alternatives"`
}

type yandexAlternative struct {
	Message yandexMessage `json:"message"`
	Status  string        `json:"status"`
}

type yandexMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type yandexResponseData struct {
	Result yandexResult `json:"result"`
}
