package v1

import "fmt"

// var folderId = viper.GetString("yandex.folder_id")

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

func getYandexRequestBody(folderId string, maxTokens string, temperature float32) *YandexRequest {
	if folderId == "" || maxTokens == "" || temperature < 0 {
		return nil
	}

	return &YandexRequest{
		ModelURI: fmt.Sprintf("gpt://%s/yandexgpt-lite", folderId),
		CompletionOptions: struct {
			Stream      bool    "json:\"stream\""
			Temperature float32 "json:\"temperature\""
			MaxTokens   string  "json:\"maxTokens\""
		}{
			Stream:      false,
			Temperature: temperature,
			MaxTokens:   maxTokens,
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
				Text: USER_PROMPT,
			},
		},
	}
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
