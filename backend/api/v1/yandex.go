package v1

import "github.com/spf13/viper"

var folderId = viper.GetString("yandex.folder_id")
var apiKey = viper.GetString("yandex.api_key")

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

var yandexRequestBody = YandexRequest{
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
			Text: USER_PROMPT,
		},
	},
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
